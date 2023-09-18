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
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/gluon/rfc822"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message"
	"github.com/bradenaw/juniper/parallel"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type WriteStage struct {
	tempPath         string
	dirPath          string
	panicHandler     async.PanicHandler
	log              *logrus.Entry
	progressReporter StageProgressReporter
	parallelWriters  int
}

func NewWriteStage(
	tempPath string,
	dirPath string,
	parallelWriters int,
	log *logrus.Entry,
	progressReporter StageProgressReporter,
	panicHandler async.PanicHandler,
) *WriteStage {
	return &WriteStage{
		tempPath:         tempPath,
		dirPath:          dirPath,
		panicHandler:     panicHandler,
		parallelWriters:  parallelWriters,
		progressReporter: progressReporter,
		log:              log.WithField("stage", "write"),
	}
}

func (w *WriteStage) Run(ctx context.Context, inputs <-chan BuildStageOutput, errReporter StageErrorReporter) {
	w.log.Debug("Starting")
	defer w.log.Debug("Exiting")

	for input := range inputs {
		if ctx.Err() != nil {
			return
		}

		if err := parallel.DoContext(ctx, w.parallelWriters, len(input.messages), func(ctx context.Context, i int) error {
			metadata := input.messages[i].GetMetadata()
			metadataPath := filepath.Join(w.dirPath, metadata.ID+".metadata.json")

			integrityChecker := &utils.Sha256IntegrityChecker{}

			metadataBytes, err := metadata.toBytes()
			if err != nil {
				w.log.WithField("msg-id", metadata.ID).WithError(err).Error("Failed to generate metadata")
				return fmt.Errorf("failed to generate message metadata: %w", err)
			}

			if err := utils.WriteFileSafe(w.tempPath, metadataPath, metadataBytes, integrityChecker); err != nil {
				w.log.WithField("msg-id", metadata.ID).WithError(err).Errorf("Failed to write %v", metadataPath)
				return fmt.Errorf("failed to write '%v': %w", metadata, err)
			}

			return input.messages[i].WriteMessage(w.dirPath, w.tempPath, w.log, integrityChecker)
		}); err != nil {
			errReporter.ReportStageError(err)
			return
		}

		w.progressReporter.OnProgress(len(input.messages))
	}
}

type MessageMetadata struct {
	proton.MessageMetadata
	Attachments []proton.Attachment
	MIMEType    rfc822.MIMEType
	Headers     string
	WriterType  MessageWriterType
}

func NewMessageMetadata(writerType MessageWriterType, msg *proton.Message) MessageMetadata {
	return MessageMetadata{
		MessageMetadata: msg.MessageMetadata,
		Headers:         msg.Header,
		Attachments:     msg.Attachments,
		MIMEType:        msg.MIMEType,
		WriterType:      writerType,
	}
}

func (m *MessageMetadata) toBytes() ([]byte, error) {
	// GODT-2925 version metadata.
	return json.MarshalIndent(m, "", "  ")
}

type MessageWriterType int

const (
	MessageWriterTypeDecryptedAndBuilt MessageWriterType = iota
	MessageWriterTypeFailedToAssemble
	MessageWriterTypeNoAddrKey
)

type MessageWriter interface {
	WriteMessage(dir string, tempDir string, log *logrus.Entry, checker utils.IntegrityChecker) error
	GetMetadata() MessageMetadata
}

type DecryptedAndBuiltMessageWriter struct {
	msg proton.FullMessage
	eml bytes.Buffer
}

func (d *DecryptedAndBuiltMessageWriter) WriteMessage(dir string, tempDir string, log *logrus.Entry, integrityChecker utils.IntegrityChecker) error {
	filePath := filepath.Join(dir, d.msg.ID)
	filePath += ".eml"

	if err := utils.WriteFileSafe(tempDir, filePath, d.eml.Bytes(), integrityChecker); err != nil {
		log.WithField("msg-id", d.msg.ID).WithError(err).Errorf("Failed to write file %v", filePath)
		return fmt.Errorf("failed to write metadata '%v': %w", filePath, err)
	}

	return nil
}

func (d *DecryptedAndBuiltMessageWriter) GetMetadata() MessageMetadata {
	return NewMessageMetadata(MessageWriterTypeDecryptedAndBuilt, &d.msg.Message)
}

type AssembleFailedMessageWriter struct {
	decrypted message.DecryptedMessage
}

func (a *AssembleFailedMessageWriter) WriteMessage(dir string, tempDir string, log *logrus.Entry, integrityChecker utils.IntegrityChecker) error {
	// Failed to assemble message, write body and attachments in a folder with the message id.
	exportDir := filepath.Join(dir, a.decrypted.Msg.ID)
	var bodyPath string

	if err := os.MkdirAll(exportDir, 0o700); err != nil {
		return fmt.Errorf("failed to create '%v': %w", exportDir, err)
	}

	// write body.
	var bodyBytes []byte
	if a.decrypted.BodyErr == nil {
		bodyBytes = a.decrypted.Body.Bytes()
		bodyPath = filepath.Join(exportDir, bodyFileName())
	} else {
		bodyBytes = []byte(a.decrypted.Msg.Body)
		bodyPath = filepath.Join(exportDir, bodyFileNameEncrypted())
	}

	if err := utils.WriteFileSafe(tempDir, bodyPath, bodyBytes, integrityChecker); err != nil {
		log.WithField("msg-id", a.decrypted.Msg.ID).WithError(err).Errorf("Failed to write %v", bodyPath)
		return fmt.Errorf("failed to write '%v': %w", bodyPath, err)
	}

	// Write attachments.
	for idx, attachment := range a.decrypted.Attachments {
		attachmentInfo := a.decrypted.Msg.Attachments[idx]
		var attachmentPath string

		var attBytes []byte
		if attachment.Err == nil {
			attBytes = attachment.Data.Bytes()
			attachmentPath = filepath.Join(exportDir, attachmentFileName(attachmentInfo.ID, attachmentInfo.Name))
		} else {
			attBytes = attachment.Encrypted
			attachmentPath = filepath.Join(exportDir, attachmentFileNameEncrypted(attachmentInfo.ID, attachmentInfo.Name))
		}

		if err := utils.WriteFileSafe(tempDir, attachmentPath, attBytes, integrityChecker); err != nil {
			log.WithField("msg-id", a.decrypted.Msg.ID).WithField("attID", attachmentInfo.ID).WithError(err).Errorf("Failed to write %v", attachmentPath)
			return fmt.Errorf("failed to write '%v': %w", attachmentPath, err)
		}

	}

	return nil
}

func (a *AssembleFailedMessageWriter) GetMetadata() MessageMetadata {
	return NewMessageMetadata(MessageWriterTypeFailedToAssemble, &a.decrypted.Msg)
}

type AddrKeyRingMissingMessageWriter struct {
	msg proton.FullMessage
}

func (a *AddrKeyRingMissingMessageWriter) GetMetadata() MessageMetadata {
	return NewMessageMetadata(MessageWriterTypeNoAddrKey, &a.msg.Message)
}

func (a *AddrKeyRingMissingMessageWriter) WriteMessage(dir string, tempDir string, log *logrus.Entry, integrityChecker utils.IntegrityChecker) error {
	// Failed decrypt due to lack of addr keyring. Write everything as pgp files to disk.
	exportDir := filepath.Join(dir, a.msg.ID)

	if err := os.MkdirAll(exportDir, 0o700); err != nil {
		return fmt.Errorf("failed to create '%v': %w", exportDir, err)
	}

	// write body.
	bodyPath := filepath.Join(exportDir, bodyFileNameEncrypted())

	if err := utils.WriteFileSafe(tempDir, bodyPath, []byte(a.msg.Body), integrityChecker); err != nil {
		log.WithField("msg-id", a.msg.ID).WithError(err).Errorf("Failed to write %v", bodyPath)
		return fmt.Errorf("failed to write '%v': %w", bodyPath, err)
	}

	// Write attachments.
	for idx, attachment := range a.msg.Attachments {
		attachmentPath := filepath.Join(exportDir, attachmentFileNameEncrypted(attachment.ID, attachment.Name))

		if err := utils.WriteFileSafe(tempDir, attachmentPath, a.msg.AttData[idx], integrityChecker); err != nil {
			log.WithField("msg-id", a.msg.ID).WithField("attID", attachment.ID).WithError(err).Errorf("Failed to write %v", attachmentPath)
			return fmt.Errorf("failed to write '%v': %w", attachmentPath, err)
		}

	}

	return nil
}

func attachmentFileName(id, name string) string {
	return fmt.Sprintf("%v_%v", id, name)
}

func attachmentFileNameEncrypted(id, name string) string {
	return fmt.Sprintf("%v_%v.pgp", id, name)
}

func bodyFileName() string {
	return "body.txt"
}

func bodyFileNameEncrypted() string {
	return "body.pgp"
}
