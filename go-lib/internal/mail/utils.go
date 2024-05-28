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
	"github.com/ProtonMail/go-proton-api"
)

const emlExtension = ".eml"
const jsonMetadataExtension = ".metadata.json"

func nonSystemLabel(label proton.Label) bool {
	return label.Type != proton.LabelTypeSystem
}

func chunkMemLimit[T any](batch []T, maxMemory uint64, stageMultiplier uint64, getSize func(T) uint64) [][]T {
	var expectedMemUsage uint64
	var chunks [][]T
	var lastIndex int
	var index int

	for _, v := range batch {
		dataSize := getSize(v)

		// 2x increase for attachment due to extra memory needed for decrypting and writing
		// in memory buffer.
		dataSize *= stageMultiplier

		nextMemSize := expectedMemUsage + dataSize
		if nextMemSize >= maxMemory {
			chunks = append(chunks, batch[lastIndex:index])
			lastIndex = index
			expectedMemUsage = dataSize
		} else {
			expectedMemUsage = nextMemSize
		}

		index++
	}

	if lastIndex < len(batch) {
		chunks = append(chunks, batch[lastIndex:])
	}

	return chunks
}
