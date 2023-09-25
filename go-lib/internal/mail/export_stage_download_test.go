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
	"testing"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestDownloadMessageAndAttachments_NoAttachments(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)

	const msgID = "msgID"

	metaData := proton.MessageMetadata{
		ID: msgID,
	}

	msgData := proton.Message{
		MessageMetadata: metaData,
		Header:          "Foo",
		Body:            "MsgBody",
		MIMEType:        "",
		Attachments:     nil,
	}

	expected := proton.FullMessage{
		Message: msgData,
	}

	client.EXPECT().GetMessage(gomock.Any(), gomock.Eq(msgID)).Return(msgData, nil)

	fullMsg, err := downloadMessageAndAttachments(context.Background(), client, metaData)
	require.NoError(t, err)
	require.Equal(t, expected, fullMsg)
}

func TestDownloadMessageAndAttachments_WithAttachments(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)

	const msgID = "msgID"
	const attID1 = "att1"
	const attID2 = "att1"
	attData1 := []byte("hello")
	attData2 := []byte("world!")

	metaData := proton.MessageMetadata{
		ID: msgID,
	}

	msgData := proton.Message{
		MessageMetadata: metaData,
		Header:          "Foo",
		Body:            "MsgBody",
		MIMEType:        "",
		Attachments: []proton.Attachment{
			{
				ID:   attID1,
				Size: int64(len(attData1)),
			},
			{
				ID:   attID2,
				Size: int64(len(attData2)),
			},
		},
	}

	expected := proton.FullMessage{
		Message: msgData,
		AttData: [][]byte{attData1, attData2},
	}

	client.EXPECT().GetMessage(gomock.Any(), gomock.Eq(msgID)).Return(msgData, nil)
	client.EXPECT().GetAttachmentInto(gomock.Any(), gomock.Eq(attID1), gomock.Any()).DoAndReturn(func(_ context.Context, _ string, b *bytes.Buffer) error {
		_, err := b.Write(attData1)
		require.NoError(t, err)
		return nil
	})
	client.EXPECT().GetAttachmentInto(gomock.Any(), gomock.Eq(attID2), gomock.Any()).DoAndReturn(func(_ context.Context, _ string, b *bytes.Buffer) error {
		_, err := b.Write(attData2)
		require.NoError(t, err)
		return nil
	})

	fullMsg, err := downloadMessageAndAttachments(context.Background(), client, metaData)
	require.NoError(t, err)
	require.Equal(t, expected, fullMsg)
}

func TestDownloadStage_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)
	stage := NewDownloadStage(client, 2, logrus.WithField("test", "test"), &async.NoopPanicHandler{})

	input := make(chan []proton.MessageMetadata)

	const msgID1 = "msgID1"
	const msgID2 = "msgID2"
	const attID1 = "att1"
	const attID2 = "att1"
	attData1 := []byte("hello")
	attData2 := []byte("world!")

	metaData := proton.MessageMetadata{
		ID: msgID1,
	}

	msgData := proton.Message{
		MessageMetadata: metaData,
		Header:          "Foo",
		Body:            "MsgBody",
		MIMEType:        "",
		Attachments: []proton.Attachment{
			{
				ID:   attID1,
				Size: int64(len(attData1)),
			},
			{
				ID:   attID2,
				Size: int64(len(attData2)),
			},
		},
	}

	msgError := &proton.APIError{Status: 422}

	inputMetadata := []proton.MessageMetadata{
		{
			ID: msgID1,
		},
		{
			ID: msgID2,
		},
	}

	expected := DownloadStageOutput{
		messages: []proton.FullMessage{
			{
				Message: msgData,
				AttData: [][]byte{attData1, attData2},
			},
		},
	}

	client.EXPECT().GetMessage(gomock.Any(), gomock.Eq(msgID1)).Return(proton.Message{}, msgError)
	client.EXPECT().GetMessage(gomock.Any(), gomock.Eq(msgID2)).Return(msgData, nil)
	client.EXPECT().GetAttachmentInto(gomock.Any(), gomock.Eq(attID1), gomock.Any()).DoAndReturn(func(_ context.Context, _ string, b *bytes.Buffer) error {
		_, err := b.Write(attData1)
		require.NoError(t, err)
		return nil
	})
	client.EXPECT().GetAttachmentInto(gomock.Any(), gomock.Eq(attID2), gomock.Any()).DoAndReturn(func(_ context.Context, _ string, b *bytes.Buffer) error {
		_, err := b.Write(attData2)
		require.NoError(t, err)
		return nil
	})

	go func() {
		stage.Run(context.Background(), input, errReporter)
	}()

	input <- inputMetadata
	close(input)

	result := <-stage.outputCh

	require.Equal(t, expected, result)
}

func TestDownloadStage_RunOtherErrorsReported(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)
	stage := NewDownloadStage(client, 2, logrus.WithField("test", "test"), &async.NoopPanicHandler{})

	input := make(chan []proton.MessageMetadata)

	const msgID1 = "msgID1"

	msgError := errors.New("unexpected error")

	inputMetadata := []proton.MessageMetadata{
		{
			ID: msgID1,
		},
	}

	client.EXPECT().GetMessage(gomock.Any(), gomock.Eq(msgID1)).Return(proton.Message{}, msgError)
	errReporter.EXPECT().ReportStageError(gomock.Eq(msgError))

	go func() {
		stage.Run(context.Background(), input, errReporter)
	}()

	input <- inputMetadata
	close(input)

	<-stage.outputCh
}
