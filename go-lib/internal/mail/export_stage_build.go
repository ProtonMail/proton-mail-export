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
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message"
	"github.com/bradenaw/juniper/parallel"
	"github.com/sirupsen/logrus"
)

type BuildStageOutput struct {
	messages []proton.FullMessage
	result   []BuildResult
}

type BuildStage struct {
	panicHandler     async.PanicHandler
	log              *logrus.Entry
	outputCh         chan BuildStageOutput
	parallelBuilders int
}

type BuildResult struct {
	eml bytes.Buffer
	err error
}

var ErrBuildNoAddrKey = errors.New("no key found for address")

func NewBuildStage(parallelBuilders int, log *logrus.Entry, panicHandler async.PanicHandler) *BuildStage {
	return &BuildStage{
		panicHandler:     panicHandler,
		log:              log.WithField("stage", "build"),
		outputCh:         make(chan BuildStageOutput),
		parallelBuilders: parallelBuilders,
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
		if ctx.Err() != nil {
			return
		}

		results := make([]BuildResult, len(input.metadata))

		if err := parallel.DoContext(ctx, b.parallelBuilders, len(results), func(ctx context.Context, i int) error {

			addrID := input.messages[i].AddressID

			kr, ok := keys.GetAddrKeyRing(addrID)
			if !ok {
				results[i].err = ErrBuildNoAddrKey
				return nil
			}

			results[i].eml.Grow(input.messages[i].Size)

			// Detected decryption errors GODT-2915.
			if err := message.BuildRFC822Into(kr, input.messages[i].Message, input.messages[i].AttData, defaultMessageJobOpts(), &results[i].eml); err != nil {
				results[i].err = err
				return nil
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
			messages: input.messages,
			result:   results,
		}:
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
