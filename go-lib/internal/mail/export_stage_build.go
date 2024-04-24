// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Mail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message"
	"github.com/bradenaw/juniper/parallel"
	"github.com/sirupsen/logrus"
)

type BuildStageOutput struct {
	lastMessageID string
	messages      []MessageWriter
}

type BuildStage struct {
	panicHandler     async.PanicHandler
	log              *logrus.Entry
	outputCh         chan BuildStageOutput
	parallelBuilders int
	maxBuildMemMB    uint64
	reporter         reporter.Reporter
	userID           string
}

var ErrBuildNoAddrKey = errors.New("no key found for address")

func NewBuildStage(
	parallelBuilders int,
	log *logrus.Entry,
	maxBuildMemMB uint64,
	panicHandler async.PanicHandler,
	reporter reporter.Reporter,
	userID string,
) *BuildStage {
	return &BuildStage{
		panicHandler:     panicHandler,
		log:              log.WithField("stage", "build"),
		outputCh:         make(chan BuildStageOutput),
		parallelBuilders: parallelBuilders,
		maxBuildMemMB:    maxBuildMemMB,
		reporter:         reporter,
		userID:           userID,
	}
}

func (b *BuildStage) Run(
	ctx context.Context,
	inputs <-chan DownloadStageOutput,
	keys *apiclient.UnlockedKeyRing,
	errReporter StageErrorReporter,
) {
	b.log.Debug("Starting")
	defer b.log.Debug("Exiting")
	defer close(b.outputCh)

	for input := range inputs {
		for _, chunk := range chunkMemLimitFullMessage(input.messages, b.maxBuildMemMB) {
			if ctx.Err() != nil {
				return
			}

			results := make([]MessageWriter, len(chunk))

			if err := parallel.DoContext(ctx, b.parallelBuilders, len(results), func(_ context.Context, i int) error {
				addrID := chunk[i].AddressID

				kr, ok := keys.GetAddrKeyRing(addrID)
				if !ok {
					b.log.WithField("addrID", addrID).Warn("Address has no key ring")
					results[i] = &AddrKeyRingMissingMessageWriter{msg: chunk[i]}
					return nil
				}

				var buffer bytes.Buffer
				buffer.Grow(chunk[i].Size)

				decrypted := message.DecryptMessage(kr, chunk[i].Message, chunk[i].AttData)

				if err := message.BuildRFC822Into(kr, &decrypted, defaultMessageJobOpts(), &buffer); err != nil {
					b.log.WithError(err).WithField("addrID", addrID).Warn("Failed to build message")
					b.reporter.ReportError(fmt.Errorf("failed to build message: %w", err), reporter.Context{
						"msgID":  chunk[i].Message.ID,
						"userID": b.userID,
					})
					results[i] = &AssembleFailedMessageWriter{decrypted: decrypted}
					return nil
				}

				results[i] = &DecryptedAndBuiltMessageWriter{
					msg: chunk[i],
					eml: buffer,
				}

				return nil
			}); err != nil {
				errReporter.ReportStageError(err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case b.outputCh <- BuildStageOutput{
				lastMessageID: chunk[len(chunk)-1].ID,
				messages:      results,
			}:
			}
		}
	}
}

func defaultMessageJobOpts() message.JobOptions {
	return message.JobOptions{
		IgnoreDecryptionErrors: true, // Whether to ignore decryption errors and create a "custom message" instead.
		SanitizeDate:           true, // Whether to replace all dates before 1970 with RFC822's birthdate.
		AddInternalID:          true, // Whether to include MessageID as X-Pm-Internal-Id.
		AddExternalID:          true, // Whether to include ExternalID as X-Pm-External-Id.
		AddMessageDate:         true, // Whether to include message time as X-Pm-Date.
		AddMessageIDReference:  true, // Whether to include the MessageID in References.
	}
}

func chunkMemLimitFullMessage(batch []proton.FullMessage, maxMemory uint64) [][]proton.FullMessage {
	// Message are alive for 2 stages.
	const stageMultiplier = 2

	return chunkMemLimit(batch, maxMemory, stageMultiplier, func(message proton.FullMessage) uint64 {
		var dataSize uint64
		for _, a := range message.Attachments {
			dataSize += uint64(a.Size)
		}
		dataSize += uint64(len(message.Body))

		return dataSize
	})
}
