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
	"strconv"
	"strings"

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

func emlToMetadataFilename(emlPath string) string {
	result, _ := strings.CutSuffix(emlPath, emlExtension)
	return result + jsonMetadataExtension
}

// isSystemLabel returns true if the label is a built-in label (Inbox, All Mail, etc...).
func isSystemLabel(labelID string) bool {
	// At the moment system folder are reported as regular folders by backend unless client is Bridge or Web.
	// A new version of the route will correct the issue (IMEX-36). For the time being, we consider a label to be
	// system if its ID is an integer. Others labels have a base64 encoded ID.
	_, err := strconv.Atoi(labelID)
	return err == nil
}
