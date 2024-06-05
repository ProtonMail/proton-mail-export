// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
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
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message/parser"
	"github.com/bradenaw/juniper/stream"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Message struct {
	literal  []byte
	metadata proton.MessageMetadata
}

const messageBatchSize = 10

func (r *RestoreTask) importMails(reporter Reporter) error {
	return r.withAddrKR(func(addrID string, addrKR *crypto.KeyRing) error {
		messages := make([]Message, 0, messageBatchSize)
		err := r.walkBackupDir(func(emlPath string) {
			literal, err := os.ReadFile(emlPath) //nolint:gosec
			if err != nil {
				logrus.WithField("path", emlPath).Error("Could not read EML file. Skipping.")
				return
			}

			metadataBytes, err := os.ReadFile(emlToMetadataFilename(emlPath)) //nolint:gosec
			if err != nil {
				logrus.WithField("path", emlPath).Error("Could not read JSON metadata file. Skipping.")
				return
			}

			metadata, err := utils.NewVersionedJSON[proton.MessageMetadata](MessageMetadataVersion, metadataBytes)
			if err != nil {
				logrus.WithError(err).WithField("path", emlPath).Error("Message metadata is invalid.")
				return
			}

			messages = append(messages, Message{literal: literal, metadata: metadata.Payload})
			if len(messages) >= messageBatchSize {
				r.importMailBatch(addrID, addrKR, messages, reporter)
				messages = messages[:0]
			}
		})

		if len(messages) > 0 {
			r.importMailBatch(addrID, addrKR, messages, reporter)
		}

		return err
	})
}

func (r *RestoreTask) importMailBatch(addrID string, addrKR *crypto.KeyRing, messages []Message, reporter Reporter) {
	reqs := make([]proton.ImportReq, 0, len(messages))
	for _, message := range messages {
		log := r.log.WithField("messageID", message.metadata.AddressID)
		labelIDs, err := r.getLabelList(message.metadata.LabelIDs)
		if err != nil {
			log.WithField("messageID", message.metadata.ID).WithError(err).Error("Could not map label to remote labels.")
			continue
		}

		parser, err := parser.New(bytes.NewReader(message.literal))
		if err != nil {
			log.WithField(message.metadata.ID, message.metadata).WithError(err).Error("Failed to parse literal for message.")
			continue
		}

		// multipart body requires at least one text part to be properly encrypted.
		if parser.AttachEmptyTextPartIfNoneExists() {
			buf := new(bytes.Buffer)
			if err := parser.NewWriter().Write(buf); err != nil {
				log.WithError(err).Error("failed to add an empty text body.")
				continue
			}
			message.literal = buf.Bytes()
		}

		reqs = append(reqs, proton.ImportReq{
			Metadata: proton.ImportMetadata{
				AddressID: addrID,
				LabelIDs:  labelIDs,
				Unread:    message.metadata.Unread,
				Flags:     message.metadata.Flags,
			},
			Message: message.literal,
		})
	}

	if len(reqs) == 0 {
		return
	}
	str, err := r.session.GetClient().ImportMessages(r.ctx, addrKR, -1, -1, reqs...)
	if err != nil {
		r.log.WithError(err).Error("failed to prepare message batch for import")
		return
	}

	results, err := stream.Collect(r.ctx, stream.Stream[proton.ImportRes](str))
	reporter.OnProgress(len(results))
	for _, result := range results {
		if err != nil {
			r.log.WithField("messageID", result.MessageID).WithError(err).Error("failed to import message")
		}

		r.log.WithField("messageID", result.MessageID).WithField("newMessageID", result.MessageID).Info("Message was imported.)")
	}
}

func (r *RestoreTask) getLabelList(labels []string) ([]string, error) {
	var result = make([]string, 0, len(labels)+1)
	result = append(result, r.importLabelID)

	for _, label := range labels {
		if !IsAcceptableLabel(label) {
			continue
		}
		remoteLabel, ok := r.labelMapping[label]
		if !ok {
			return nil, fmt.Errorf("could not find a remote label matching backup label %v", label)
		}

		result = append(result, remoteLabel)
	}

	return result, nil
}

func IsAcceptableLabel(label string) bool {
	var acceptableLabel = []string{
		proton.InboxLabel,
		proton.AllDraftsLabel,
		proton.AllSentLabel,
		proton.TrashLabel,
		proton.SpamLabel,
		proton.ArchiveLabel,
		proton.SentLabel,
		proton.DraftsLabel,
		proton.OutboxLabel,
		proton.StarredLabel,
		proton.AllScheduledLabel,
	} // proton.AllMailLabel is discarded on purpose as it would cause import to fail.
	return (len(label) > 4) || slices.Contains(acceptableLabel, label)
}
