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
	"github.com/bradenaw/juniper/parallel"
	"github.com/bradenaw/juniper/xslices"
	"github.com/sirupsen/logrus"
)

type DownloadStageOutput struct {
	messages []proton.FullMessage
}

type DownloadStage struct {
	client           apiclient.Client
	log              *logrus.Entry
	outputCh         chan DownloadStageOutput
	parallelWorkers  int
	maxDownloadMemMB uint64
	panicHandler     async.PanicHandler
}

func NewDownloadStage(
	client apiclient.Client,
	parallelWorkers int,
	log *logrus.Entry,
	maxDownloadMemMB uint64,
	panicHandler async.PanicHandler,
) *DownloadStage {
	return &DownloadStage{
		client:           client,
		log:              log.WithField("stage", "download"),
		outputCh:         make(chan DownloadStageOutput),
		parallelWorkers:  parallelWorkers,
		panicHandler:     panicHandler,
		maxDownloadMemMB: maxDownloadMemMB,
	}
}

func (d *DownloadStage) Run(ctx context.Context, input <-chan []proton.MessageMetadata, errReporter StageErrorReporter) {
	d.log.Debug("Starting")
	defer d.log.Debug("Exiting")

	const Failed422ID = "MsgFailed422"

	defer close(d.outputCh)
	for metadata := range input {
		memChucked := chunkMemLimitMetadata(metadata, d.maxDownloadMemMB)
		for _, chunk := range memChucked {
			if ctx.Err() != nil {
				return
			}

			result := DownloadStageOutput{
				messages: make([]proton.FullMessage, len(chunk)),
			}

			if err := parallel.DoContext(ctx, d.parallelWorkers, len(chunk), func(ctx context.Context, i int) error {
				defer async.HandlePanic(d.panicHandler)

				msg, err := downloadMessageAndAttachments(ctx, d.client, chunk[i])
				if err != nil {
					var apiErr *proton.APIError
					if errors.As(err, &apiErr) && apiErr.Status == 422 {
						d.log.WithField("msgID", chunk[i].ID).Warn("Failed to download message due to 422")
						result.messages[i].ID = Failed422ID
						return nil
					}

					d.log.WithError(err).WithField("msgID", chunk[i].ID).Error("Failed to download message or attachment")
					return err
				}

				result.messages[i] = msg

				return nil
			}); err != nil {
				errReporter.ReportStageError(err)
				return
			}

			// Remove any failed 422 downloads.
			result.messages = xslices.Filter(result.messages, func(t proton.FullMessage) bool {
				return t.ID != Failed422ID
			})

			select {
			case <-ctx.Done():
				return
			case d.outputCh <- result:
			}
		}
	}
}

func downloadMessageAndAttachments(ctx context.Context, client apiclient.Client, metadata proton.MessageMetadata) (proton.FullMessage, error) {
	msg, err := client.GetMessage(ctx, metadata.ID)
	if err != nil {
		return proton.FullMessage{}, err
	}

	full := proton.FullMessage{
		Message: msg,
		AttData: nil,
	}

	if len(msg.Attachments) != 0 {
		attData := make([][]byte, len(msg.Attachments))

		for i, a := range msg.Attachments {
			buffer := bytes.Buffer{}

			buffer.Grow(int(a.Size))

			if err := client.GetAttachmentInto(ctx, a.ID, &buffer); err != nil {
				return proton.FullMessage{}, err
			}

			attData[i] = buffer.Bytes()
		}

		full.AttData = attData
	}

	return full, nil
}

func chunkMemLimitMetadata(batch []proton.MessageMetadata, maxMemory uint64) [][]proton.MessageMetadata {
	// Message are alive for 4 stages. Even though there are technically 2 stages after this one
	// Due to pipelining up to 4 batches can be in circulation at any given time.
	const stageMultiplier = 4

	return chunkMemLimit(batch, maxMemory, stageMultiplier, func(message proton.MessageMetadata) uint64 {
		return uint64(message.Size)
	})
}
