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
	"path/filepath"

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

const messageBatchSize = 10 // max batch size supported by go-proton-api (larger batches will be split).

func (r *RestoreTask) importMails(messageInfoList []messageInfo, reporter Reporter) error {
	return r.withAddrKR(func(addrID string, addrKR *crypto.KeyRing) error {
		messages := make([]Message, 0, messageBatchSize)
		for _, info := range messageInfoList {
			emlPath := filepath.Join(r.backupDir, info.messageID+emlExtension)
			literal, err := os.ReadFile(emlPath) //nolint:gosec
			if err != nil {
				logrus.WithField("path", emlPath).Error("Could not read EML file. Skipping.")
				continue
			}

			metadataPath := emlToMetadataFilename(emlPath)
			metadata, err := loadMetadataFile(metadataPath)
			if err != nil {
				logrus.WithField("path", metadataPath).Error("Could not load metadata file. Skipping.")
			}

			messages = append(messages, Message{literal: literal, metadata: metadata.MessageMetadata})
			if len(messages) >= messageBatchSize {
				if err := r.importMailBatch(addrID, addrKR, messages, reporter); err != nil {
					return err
				}
				messages = messages[:0]
			}
		}

		if len(messages) > 0 {
			if err := r.importMailBatch(addrID, addrKR, messages, reporter); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RestoreTask) importMailBatch(addrID string, addrKR *crypto.KeyRing, messages []Message, reporter Reporter) error {
	reqs := make([]proton.ImportReq, 0, len(messages))
	for _, message := range messages {
		log := r.log.WithField("messageID", message.metadata.AddressID)
		labelIDs, err := r.getLabelList(message.metadata.LabelIDs)
		if err != nil {
			log.WithField("messageID", message.metadata.ID).WithError(err).Error("Could not map label to remote labels.")
			continue
		}

		msgParser, err := parser.New(bytes.NewReader(message.literal))
		if err != nil {
			log.WithField(message.metadata.ID, message.metadata).WithError(err).Error("Failed to parse literal for message.")
			continue
		}

		// multipart body requires at least one text part to be properly encrypted.
		if msgParser.AttachEmptyTextPartIfNoneExists() {
			buf := new(bytes.Buffer)
			if err := msgParser.NewWriter().Write(buf); err != nil {
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
		return nil
	}

	str, err := r.session.GetClient().ImportMessages(r.ctx, addrKR, -1, -1, reqs...)
	if err != nil {
		r.log.WithError(err).Error("failed to prepare message batch for import")
		return err
	}

	results, err := stream.Collect(r.ctx, stream.Stream[proton.ImportRes](str))
	if err != nil {
		r.log.WithError(err).Error("An error occurred while importing a batch of messages. Retrying one by one.")
		r.importOneByOne(reqs, messages, addrKR, reporter)
		return nil
	}

	for i, result := range results {
		if result.Code != 1000 {
			r.log.WithField("messageID", messages[i].metadata.ID).WithError(result.APIError).Error("Failed to import message")
			r.failedCount++
		} else {
			r.importedCount++
		}
	}
	reporter.OnProgress(len(results))

	return nil
}

func (r *RestoreTask) importOneByOne(requests []proton.ImportReq, messages []Message, addrKR *crypto.KeyRing, reporter Reporter) {
	for i, request := range requests {
		resultStream, err := r.session.GetClient().ImportMessages(r.ctx, addrKR, -1, -1, request)
		if err != nil {
			r.log.WithError(err).WithField("messageID", messages[i].metadata.ID).Error("Failed to import message")
			r.failedCount++
			continue
		}

		results, err := stream.Collect(r.ctx, stream.Stream[proton.ImportRes](resultStream))
		if err != nil {
			r.log.WithError(err).WithField("messageID", messages[i].metadata.ID).Error("Failed to import message")
			r.failedCount++
			continue
		}

		if results[0].Code != 1000 {
			r.log.WithField("messageID", messages[i].metadata.ID).WithError(results[0].APIError).Error("Failed to import message")
			r.failedCount++
		} else {
			r.importedCount++
		}
	}
	reporter.OnProgress(len(requests))
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

func IsAcceptableLabel(labelID string) bool {
	var acceptableLabel = []string{
		proton.InboxLabel,
		proton.TrashLabel,
		proton.SpamLabel,
		proton.ArchiveLabel,
		proton.SentLabel,
		proton.DraftsLabel,
		proton.OutboxLabel,
		proton.StarredLabel,
	} // proton.AllMailLabel is discarded on purpose as it would cause import to fail.
	return (!isSystemLabel(labelID)) || slices.Contains(acceptableLabel, labelID)
}
