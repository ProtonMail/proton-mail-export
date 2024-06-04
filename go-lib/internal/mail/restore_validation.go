// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
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
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (r *RestoreTask) validateBackupDir(reporter Reporter) error {
	r.log.Info("Verifying backup folder")

	var importableCount uint64
	err := r.walkBackupDir(func(_ string) {
		importableCount++
	})

	if err != nil {
		return err
	}

	if importableCount > 0 {
		labelsFilename := getLabelFileName()
		if _, err := os.Stat(filepath.Join(r.backupDir, labelsFilename)); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the labels file '%v' could not be found", labelsFilename)
		}

		reporter.SetMessageTotal(importableCount)
		reporter.SetMessageProcessed(0)
		r.log.WithField("messageCount", importableCount).Info("Found importable messages")
		return nil
	}

	subDirs, err := r.getTimestampedBackupDirs()
	if err != nil {
		return err
	}

	if len(subDirs) == 0 {
		return errors.New("no importable mail found")
	}

	if len(subDirs) > 1 {
		return errors.New("the specified folder contains more than one backup sub-folder")
	}

	r.log.WithField("folderName", subDirs[0]).Info("A potential backup sub-folder has been found and will be inspected")
	r.backupDir = subDirs[0]

	return r.validateBackupDir(reporter)
}
