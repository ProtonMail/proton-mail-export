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

package apiclient

import (
	"context"
	"encoding/json"
	"github.com/ProtonMail/go-proton-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestAutoRetryClientRepetition(t *testing.T) {
	type test struct {
		name        string
		err         error
		expectRetry bool
	}

	var tests = []test{
		{
			name:        "429",
			err:         &proton.APIError{Status: 429},
			expectRetry: true,
		},
		{
			name:        "500",
			err:         &proton.APIError{Status: 500},
			expectRetry: true,
		},
		{
			name:        "505",
			err:         &proton.APIError{Status: 505},
			expectRetry: true,
		},
		{
			name:        "ProtonNetError",
			err:         &proton.NetError{},
			expectRetry: true,
		},
		{
			name:        "OSNetError",
			err:         &net.OpError{},
			expectRetry: true,
		},
		{
			name:        "UnexpectedEOF",
			err:         io.ErrUnexpectedEOF,
			expectRetry: true,
		},
		{
			name:        "CtxCancel",
			err:         context.Canceled,
			expectRetry: false,
		},
		{
			name:        "CtxCancel",
			err:         context.Canceled,
			expectRetry: false,
		},
		{
			name:        "422",
			err:         &proton.APIError{Status: 422},
			expectRetry: false,
		},
		{
			name:        "jsonUnmarshall",
			err:         &json.UnmarshalTypeError{},
			expectRetry: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			strategy := NewMockRetryStrategy(mockCtrl)
			mockClient := NewMockClient(mockCtrl)

			client := NewAutoRetryClient(mockClient, &mockRetryStrategyBuilder{s: strategy})

			call1 := mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any()).Times(1).Return(proton.Message{}, test.err)
			if test.expectRetry {
				strategy.EXPECT().HandleRetry(gomock.Any()).Times(1)

				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any()).Times(1).After(call1).Return(proton.Message{}, nil)
			}

			_, err := client.GetMessage(context.Background(), "msgid")
			if test.expectRetry {
				require.NoError(t, err)
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}

type mockRetryStrategyBuilder struct {
	s *MockRetryStrategy
}

func (m mockRetryStrategyBuilder) NewRetryStrategy() RetryStrategy {
	return m.s
}
