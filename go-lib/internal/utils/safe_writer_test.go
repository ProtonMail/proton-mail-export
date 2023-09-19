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

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteFileSafe(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testFile.txt")
	data := []byte("Proton Mail Bridge is free software: you can redistribute it and/or modify")
	require.NoError(t, WriteFileSafe(tmpDir, filePath, data, &Sha256IntegrityChecker{}))
}

func TestSha256IntegrityChecker_Check(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testFile.txt")
	data := []byte("Proton Mail Bridge is free software: you can redistribute it and/or modify")
	dataCorrupt := []byte("Proton Mail Bridge is free software: you can redistribute it and/or")

	checker := &Sha256IntegrityChecker{}
	checker.Initialize(data)

	require.NoError(t, os.WriteFile(filePath, dataCorrupt, 0o700))
	require.ErrorIs(t, ErrIntegrityCheckFailed, checker.Check(filePath))
}
