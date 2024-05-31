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
	"github.com/ProtonMail/proton-bridge/v3/pkg/message/parser"
	"github.com/bradenaw/juniper/stream"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func (r *RestoreTask) importMails() error {
	var count int
	return r.walkBackupDir(func(emlPath, metadataPath string) {
		count++

		literal, err := os.ReadFile(emlPath) //nolint:gosec
		if err != nil {
			logrus.WithField("path", emlPath).Error("Could not read EML file. Skipping.")
			return
		}

		metadataBytes, err := os.ReadFile(metadataPath) //nolint:gosec
		if err != nil {
			logrus.WithField("path", emlPath).Error("Could not read JSON metadata file. Skipping.")
			return
		}

		metadata, err := utils.NewVersionedJSON[proton.MessageMetadata](MessageMetadataVersion, metadataBytes)
		if err != nil {
			logrus.WithError(err).WithField("path", emlPath).Error("Message metadata is invalid.")
			return
		}
		r.log.WithField("count", count).Info("Counting")
		r.importMail(literal, metadata.Payload)
	})
}

func (r *RestoreTask) importMail(literal []byte, metadata proton.MessageMetadata) {
	log := r.log.WithField("messageID", metadata.AddressID)
	labelIDs, err := r.getLabelList(metadata.LabelIDs)
	if err != nil {
		log.WithField("messageID", metadata.ID).WithError(err).Error("Could not map label to remote labels.")
		return
	}

	importMetadata := proton.ImportMetadata{
		AddressID: r.addrID,
		LabelIDs:  labelIDs,
		Unread:    metadata.Unread,
		Flags:     metadata.Flags,
	}

	parser, err := parser.New(bytes.NewReader(literal))
	if err != nil {
		log.WithField(metadata.ID, metadata).WithError(err).Error("Failed to parse literal for message.")
		return
	}

	addrKR, ok := r.addrKR.GetAddrKeyRing(r.addrID)
	if !ok {
		log.WithError(err).WithField("addressID", r.addrID).Error("unable to get keyring for address.")
		return
	}

	primaryKey, err := addrKR.FirstKey()
	if err != nil {
		log.WithError(err).WithField("addressID", r.addrID).Error("unable to get primary key for address.")
		return
	}

	// multipart body requires at least one text part to be properly encrypted.
	if parser.AttachEmptyTextPartIfNoneExists() {
		buf := new(bytes.Buffer)
		if err := parser.NewWriter().Write(buf); err != nil {
			log.WithError(err).Error("failed to add an empty text body.")
			return
		}
		literal = buf.Bytes()
	}

	str, err := r.session.GetClient().ImportMessages(r.ctx, primaryKey, 1, 1, []proton.ImportReq{{
		Metadata: importMetadata, Message: literal,
	}}...)
	if err != nil {
		log.WithError(err).Error("failed to prepare message for import")
		return
	}

	res, err := stream.Collect(r.ctx, stream.Stream[proton.ImportRes](str))
	if err != nil {
		log.WithError(err).Error("failed to import message")
	}

	log.WithField("newMessageID", res[0].MessageID).Info("Message was imported.)")
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
