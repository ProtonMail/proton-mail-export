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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestAddrKeyRingMissingMessageWriter(t *testing.T) {
	attData := []byte("hello attachment")
	attID := "attachment_id"
	msgID := "msg_id"
	msgBody := "hello body"

	msg := proton.FullMessage{
		Message: proton.Message{
			MessageMetadata: proton.MessageMetadata{
				ID: msgID,
			},
			Header:   "",
			Body:     msgBody,
			MIMEType: "",
			Attachments: []proton.Attachment{
				{
					ID:          attID,
					Name:        "foo",
					Size:        int64(len(attData)),
					MIMEType:    "",
					Disposition: "",
					KeyPackets:  "",
					Signature:   "",
				},
			},
		},
		AttData: [][]byte{attData},
	}

	writer := AddrKeyRingMissingMessageWriter{msg: msg}
	writeDir := t.TempDir()
	tmpDir := t.TempDir()

	checker := &utils.Sha256IntegrityChecker{}
	require.NoError(t, writer.WriteMessage(writeDir, tmpDir, logrus.WithField("t", "t"), checker))

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, attachmentFileNameEncrypted(attID, "foo")))
		require.NoError(t, err)
		require.Equal(t, attData, data)
	}

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, "body.pgp"))
		require.NoError(t, err)
		require.Equal(t, []byte(msgBody), data)
	}
}

func TestAssembleFailedMessageWriter_NotDecrypted(t *testing.T) {
	attData := []byte("hello attachment")
	attID := "attachment_id"
	msgID := "msg_id"
	msgBody := "hello body"

	msg := proton.FullMessage{
		Message: proton.Message{
			MessageMetadata: proton.MessageMetadata{
				ID: msgID,
			},
			Header:   "",
			Body:     msgBody,
			MIMEType: "",
			Attachments: []proton.Attachment{
				{
					ID:          attID,
					Name:        "foo",
					Size:        int64(len(attData)),
					MIMEType:    "",
					Disposition: "",
					KeyPackets:  "",
					Signature:   "",
				},
			},
		},
		AttData: [][]byte{attData},
	}

	writer := AssembleFailedMessageWriter{
		decrypted: message.DecryptedMessage{
			Msg:     msg.Message,
			Body:    bytes.Buffer{},
			BodyErr: fmt.Errorf("failed to decrypt body"),
			Attachments: []message.DecryptedAttachment{
				{
					Packet:    nil,
					Encrypted: attData,
					Data:      bytes.Buffer{},
					Err:       fmt.Errorf("failed to decrypt attachment"),
				},
			},
		},
	}

	writeDir := t.TempDir()
	tmpDir := t.TempDir()

	checker := &utils.Sha256IntegrityChecker{}
	require.NoError(t, writer.WriteMessage(writeDir, tmpDir, logrus.WithField("t", "t"), checker))

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, attachmentFileNameEncrypted(attID, "foo")))
		require.NoError(t, err)
		require.Equal(t, attData, data)
	}

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, "body.pgp"))
		require.NoError(t, err)
		require.Equal(t, []byte(msgBody), data)
	}
}

func TestAssembleFailedMessageWriter_AllDecrypted(t *testing.T) {
	attID := "attachment_id"
	msgID := "msg_id"

	msgBodyDecrypted := []byte("decrypted body")
	attachmentDecrypted := []byte("decrypted attachment")

	msg := proton.FullMessage{
		Message: proton.Message{
			MessageMetadata: proton.MessageMetadata{
				ID: msgID,
			},
			Header:   "",
			MIMEType: "",
			Attachments: []proton.Attachment{
				{
					ID:   attID,
					Name: "foo",
				},
			},
		},
	}

	writer := AssembleFailedMessageWriter{
		decrypted: message.DecryptedMessage{
			Msg:  msg.Message,
			Body: bytes.Buffer{},
			Attachments: []message.DecryptedAttachment{
				{
					Packet: nil,
					Data:   bytes.Buffer{},
				},
			},
		},
	}

	{
		_, err := writer.decrypted.Body.Write(msgBodyDecrypted)
		require.NoError(t, err)
	}
	{
		_, err := writer.decrypted.Attachments[0].Data.Write(attachmentDecrypted)
		require.NoError(t, err)
	}

	writeDir := t.TempDir()
	tmpDir := t.TempDir()

	checker := &utils.Sha256IntegrityChecker{}
	require.NoError(t, writer.WriteMessage(writeDir, tmpDir, logrus.WithField("t", "t"), checker))

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, attachmentFileName(attID, "foo")))
		require.NoError(t, err)
		require.Equal(t, attachmentDecrypted, data)
	}

	{
		data, err := os.ReadFile(filepath.Join(writeDir, msg.ID, "body.txt"))
		require.NoError(t, err)
		require.Equal(t, msgBodyDecrypted, data)
	}
}

func TestFileMetadataFileChecker_HasMessage_MetadataMissing(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.False(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithoutEMLOrDir(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, MessageMetadata{}, metadataPath)

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.False(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithEML(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, MessageMetadata{}, metadataPath)

	emlFile := filepath.Join(dir, getEMLFileName(messageID))
	require.NoError(t, os.WriteFile(emlFile, []byte{0, 1}, 0o600))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.True(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithDirButNoFiles(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, MessageMetadata{}, metadataPath)

	require.NoError(t, os.MkdirAll(filepath.Join(dir, metadataPath), 0o700))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.False(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithDirButMissingFiles(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadata := getTestMessageMetadata()

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, metadata, metadataPath)

	msgDir := filepath.Join(dir, metadataPath)
	require.NoError(t, os.MkdirAll(msgDir, 0o700))

	// write only body
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, bodyFileName()), []byte{0, 1}, 0o600))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.False(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithDirAndAllFiles(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadata := getTestMessageMetadata()

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, metadata, metadataPath)

	msgDir := filepath.Join(dir, messageID)
	require.NoError(t, os.MkdirAll(msgDir, 0o700))

	// write message parts.
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, bodyFileName()), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileName(metadata.Attachments[0].ID, metadata.Attachments[0].Name)), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileName(metadata.Attachments[1].ID, metadata.Attachments[1].Name)), []byte{0, 1}, 0o600))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.True(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithDirAndAllFilesNotDecrypted(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadata := getTestMessageMetadata()

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, metadata, metadataPath)

	msgDir := filepath.Join(dir, messageID)
	require.NoError(t, os.MkdirAll(msgDir, 0o700))

	// write message parts.
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, bodyFileNameEncrypted()), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileNameEncrypted(metadata.Attachments[0].ID, metadata.Attachments[0].Name)), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileNameEncrypted(metadata.Attachments[1].ID, metadata.Attachments[1].Name)), []byte{0, 1}, 0o600))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.True(t, hasMessage)
}

func TestFileMetadataFileChecker_HasMessage_MetadataWithDirAndAllFilesEncryptionMix(t *testing.T) {
	const messageID = "msg-1"
	dir := t.TempDir()
	checker := NewFileMetadataFileChecker(dir)

	metadata := getTestMessageMetadata()

	metadataPath := filepath.Join(dir, getMetadataFileName(messageID))
	writeTestMetadata(t, metadata, metadataPath)

	msgDir := filepath.Join(dir, messageID)
	require.NoError(t, os.MkdirAll(msgDir, 0o700))

	// write message parts.
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, bodyFileName()), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileNameEncrypted(metadata.Attachments[0].ID, metadata.Attachments[0].Name)), []byte{0, 1}, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, attachmentFileName(metadata.Attachments[1].ID, metadata.Attachments[1].Name)), []byte{0, 1}, 0o600))

	hasMessage, err := checker.HasMessage(messageID)
	require.NoError(t, err)
	require.True(t, hasMessage)
}

func writeTestMetadata(t *testing.T, metadata MessageMetadata, path string) {
	b, err := utils.GenerateVersionedJSON(MessageMetadataVersion, metadata)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, b, 0o600))
}

func getTestMessageMetadata() MessageMetadata {
	return MessageMetadata{
		MessageMetadata: proton.MessageMetadata{},
		Attachments: []proton.Attachment{
			{
				ID:   "attachment1",
				Name: "foo.png",
			},
			{
				ID:   "attachment2",
				Name: "bar.pdf",
			},
		},
		MIMEType:   "",
		Headers:    "",
		WriterType: 0,
	}
}
