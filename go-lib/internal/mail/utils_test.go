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
	"testing"

	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/xslices"
	"github.com/stretchr/testify/require"
)

func TestSyncChunkBuilderBatch(t *testing.T) {
	const totalMessageCount = 100

	msg := proton.FullMessage{
		Message: proton.Message{
			Attachments: []proton.Attachment{
				{
					Size: int64(8 * 1024 * 1024),
				},
			},
		},
		AttData: nil,
	}

	messages := xslices.Repeat(msg, totalMessageCount)

	chunks := chunkMemLimitFullMessage(messages, 16*1024*1024)

	var totalMessagesInChunks int

	for _, v := range chunks {
		totalMessagesInChunks += len(v)
	}

	require.Equal(t, totalMessagesInChunks, totalMessageCount)
}
