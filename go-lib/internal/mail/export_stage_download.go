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
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/parallel"
	"github.com/bradenaw/juniper/xslices"
	"github.com/sirupsen/logrus"
)

type DownloadStageOutput struct {
	metadata []proton.MessageMetadata
	messages []proton.FullMessage
	errors   []error
}

type DownloadStage struct {
	client          apiclient.Client
	log             *logrus.Entry
	outputCh        chan DownloadStageOutput
	parallelWorkers int
	panicHandler    async.PanicHandler
}

func NewDownloadStage(
	client apiclient.Client,
	parallelWorkers int,
	log *logrus.Entry,
	panicHandler async.PanicHandler,
) *DownloadStage {
	return &DownloadStage{
		client:          client,
		log:             log.WithField("stage", "download"),
		outputCh:        make(chan DownloadStageOutput),
		parallelWorkers: parallelWorkers,
		panicHandler:    panicHandler,
	}
}

func (d *DownloadStage) Run(ctx context.Context, input <-chan []proton.MessageMetadata, errReporter StageErrorReporter) {
	d.log.Debug("Starting")
	defer d.log.Debug("Exiting")

	defer close(d.outputCh)
	for metadata := range input {
		for _, chunk := range xslices.Chunk(metadata, d.parallelWorkers) {
			if ctx.Err() != nil {
				return
			}

			result := DownloadStageOutput{
				metadata: chunk,
				messages: make([]proton.FullMessage, len(chunk)),
				errors:   make([]error, len(chunk)),
			}

			if err := parallel.DoContext(ctx, d.parallelWorkers, len(chunk), func(ctx context.Context, i int) error {
				defer async.HandlePanic(d.panicHandler)

				msg, err := downloadMessageAndAttachments(ctx, d.client, chunk[i])

				result.messages[i] = msg
				result.errors[i] = err

				return nil
			}); err != nil {
				errReporter.ReportStageError(err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case d.outputCh <- result:
			}
		}
	}
}

func downloadMessageAndAttachments(ctx context.Context, client apiclient.Client, metadata proton.MessageMetadata) (proton.FullMessage, error) {
	// GODT-2900: Handle network errors/loss.
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
			// GODT-2900: Handle network errors/loss.
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
