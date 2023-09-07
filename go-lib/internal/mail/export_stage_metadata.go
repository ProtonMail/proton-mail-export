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
	"context"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type MetadataStage struct {
	client   apiclient.Client
	log      *logrus.Entry
	outputCh chan []proton.MessageMetadata
	pageSize int
}

func NewMetadataStage(client apiclient.Client, entry *logrus.Entry, pageSize int) *MetadataStage {
	return &MetadataStage{
		client:   client,
		log:      entry.WithField("stage", "metadata"),
		outputCh: make(chan []proton.MessageMetadata),
		pageSize: pageSize,
	}
}

func (m *MetadataStage) Run(ctx context.Context, errReporter StageErrorReporter) {
	m.log.Debug("Starting")
	defer m.log.Debug("Exiting")
	defer close(m.outputCh)

	var lastMessageID string

	// GODT-2900: Handle network errors/loss.
	client := m.client

	for {
		if ctx.Err() != nil {
			return
		}

		var metadata []proton.MessageMetadata

		if lastMessageID != "" {
			meta, err := client.GetMessageMetadataPage(ctx, 0, m.pageSize, proton.MessageFilter{
				EndID: lastMessageID,
				Desc:  true,
			})

			if err != nil {
				errReporter.ReportStageError(err)
				return
			}

			// * There is only one message returned and it matches the EndID query.
			if len(meta) != 0 && meta[0].ID == lastMessageID {
				meta = meta[1:]
			}

			metadata = meta
		} else {
			meta, err := client.GetMessageMetadataPage(ctx, 0, m.pageSize, proton.MessageFilter{
				Desc: true,
			})
			if err != nil {
				errReporter.ReportStageError(err)
				return
			}
			metadata = meta
		}

		// Nothing left to do
		if len(metadata) == 0 {
			return
		}

		lastMessageID = metadata[len(metadata)-1].ID

		select {
		case <-ctx.Done():
			return
		case m.outputCh <- metadata:
		}
	}
}
