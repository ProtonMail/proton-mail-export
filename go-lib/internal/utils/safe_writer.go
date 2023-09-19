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
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var ErrIntegrityCheckFailed = errors.New("integrity check failed")

type IntegrityChecker interface {
	Initialize([]byte)
	Check(path string) error
}

// WriteFileSafe writes the contents to a temporary location first, before moving to the designated location.
func WriteFileSafe(tempPath, dstPath string, data []byte, integrityChecker IntegrityChecker) error {
	if integrityChecker != nil {
		integrityChecker.Initialize(data)
	}

	file, err := os.CreateTemp(tempPath, "export-tool-*")
	if err != nil {
		return fmt.Errorf("failed to create tmp file: %w", err)
	}

	filePath := file.Name()

	written, err := file.Write(data)
	if err != nil {
		if err := file.Close(); err != nil {
			logrus.WithField("dstPath", filePath).WithError(err).Error("Failed to close tmp file after io error")
		}
		return fmt.Errorf("failed to write contents: %w", err)
	}

	if written != len(data) {
		if err := file.Close(); err != nil {
			logrus.WithField("dstPath", filePath).WithError(err).Error("Failed to close tmp file")
		}
		return fmt.Errorf("not all contents written to file")
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close tmp file: %w", err)
	}

	if integrityChecker != nil {
		if err := integrityChecker.Check(filePath); err != nil {
			return err
		}
	}

	if err := os.Rename(filePath, dstPath); err != nil {
		return fmt.Errorf("failed to move file to location: %w", err)
	}

	return nil
}

type Sha256IntegrityChecker struct {
	hash []byte
}

func (s *Sha256IntegrityChecker) Initialize(i []byte) {
	hash := sha256.Sum256(i)
	s.hash = hash[:]
}

func (s *Sha256IntegrityChecker) Check(path string) error {
	input, err := os.Open(path) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open tmp file for checksum validation:%w", err)
	}

	hasher := sha256.New()

	if _, err := io.Copy(hasher, input); err != nil {
		if err := input.Close(); err != nil {
			logrus.WithField("path", path).WithError(err).Error("Failed to close file during checksum validation")
		}
		return fmt.Errorf("failed to hash written tmp file: %w", err)
	}

	if err := input.Close(); err != nil {
		return fmt.Errorf("failed to close file after checksum validation")
	}

	onDiskHash := hasher.Sum(nil)
	if !bytes.Equal(onDiskHash, s.hash) {
		return ErrIntegrityCheckFailed
	}

	return nil
}
