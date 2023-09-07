package mail

import (
	"context"
	"fmt"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadataStage_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)

	const pageSize = 2

	expected := testMetadata(20)
	encodeMetadataExpectations(client, expected, pageSize)

	metadata := NewMetadataStage(client, logrus.WithField("test", "test"), pageSize)

	go func() {
		metadata.Run(context.Background(), errReporter)
	}()

	result := make([]proton.MessageMetadata, 0, 20)
	for out := range metadata.outputCh {
		result = append(result, out...)
	}

	require.Equal(t, expected, result)
}

func testMetadata(count int) []proton.MessageMetadata {
	result := make([]proton.MessageMetadata, count)

	for i := 0; i < count; i++ {
		result[i].ID = fmt.Sprintf("msg-%v", i)
	}

	return result
}

func encodeMetadataExpectations(client *apiclient.MockClient, metadata []proton.MessageMetadata, pageSize int) {
	filter := proton.MessageFilter{
		Desc: true,
	}

	for i := 0; i < len(metadata); i += pageSize - 1 {

		if i != 0 {
			filter.EndID = metadata[i].ID
		}

		if i+pageSize > len(metadata) {
			client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(pageSize), gomock.Eq(filter)).Return(metadata[i:], nil)
		} else {
			client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(pageSize), gomock.Eq(filter)).Return(metadata[i:i+pageSize], nil)
		}
	}

	if pageSize > 2 {
		client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(pageSize), gomock.Eq(proton.MessageFilter{
			EndID: metadata[len(metadata)-1].ID,
			Desc:  true,
		})).Return(metadata[len(metadata)-1:], nil)
	}
}
