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
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/proton-bridge/v3/pkg/message"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
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
			Header:        "",
			ParsedHeaders: nil,
			Body:          msgBody,
			MIMEType:      "",
			Attachments: []proton.Attachment{
				{
					ID:          attID,
					Name:        "foo",
					Size:        int64(len(attData)),
					MIMEType:    "",
					Disposition: "",
					Headers:     nil,
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
			Header:        "",
			ParsedHeaders: nil,
			Body:          msgBody,
			MIMEType:      "",
			Attachments: []proton.Attachment{
				{
					ID:          attID,
					Name:        "foo",
					Size:        int64(len(attData)),
					MIMEType:    "",
					Disposition: "",
					Headers:     nil,
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
			Header:        "",
			ParsedHeaders: nil,
			MIMEType:      "",
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
