package mail

import (
	"context"
	"fmt"
	"testing"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/xmaps"
	"github.com/bradenaw/juniper/xslices"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMetadataStage_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)
	fileChecker := NewMockMetadataFileChecker(mockCtrl)
	reporter := NewMockReporter(mockCtrl)

	const pageSize = 2

	expected := testMetadata(20)
	encodeMetadataExpectations(client, expected, pageSize)
	fileChecker.EXPECT().HasMessage(gomock.Any()).AnyTimes().Return(false, nil)

	metadata := NewMetadataStage(client, logrus.WithField("test", "test"), pageSize, 1)

	go func() {
		metadata.Run(context.Background(), errReporter, fileChecker, reporter)
	}()

	result := make([]proton.MessageMetadata, 0, 20)
	for out := range metadata.outputCh {
		result = append(result, out...)
	}

	require.Equal(t, expected, result)
}

func TestMetadataStage_RunWithCached(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)
	fileChecker := NewMockMetadataFileChecker(mockCtrl)
	reporter := NewMockReporter(mockCtrl)

	const pageSize = 2

	expected := testMetadata(20)
	encodeMetadataExpectations(client, expected, pageSize)

	reporter.EXPECT().OnProgress(1).MinTimes(20 / 3)

	filteredIDs := make(xmaps.Set[string])

	for idx, m := range expected {
		if idx%3 == 0 {
			fileChecker.EXPECT().HasMessage(gomock.Eq(m.ID)).Return(true, nil)
			filteredIDs.Add(m.ID)
		} else {
			fileChecker.EXPECT().HasMessage(gomock.Eq(m.ID)).Return(false, nil)
		}
	}

	metadata := NewMetadataStage(client, logrus.WithField("test", "test"), pageSize, 1)

	go func() {
		metadata.Run(context.Background(), errReporter, fileChecker, reporter)
	}()

	result := make([]proton.MessageMetadata, 0, 20)
	for out := range metadata.outputCh {
		result = append(result, out...)
	}

	expectedFiltered := xslices.Filter(expected, func(t proton.MessageMetadata) bool {
		return !filteredIDs.Contains(t.ID)
	})

	require.Equal(t, len(expected)-len(filteredIDs), len(expectedFiltered))
	require.Equal(t, expectedFiltered, result)
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
