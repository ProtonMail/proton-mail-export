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
	"strings"
)

func (r *RestoreTask) validateBackupDir() error {
	r.log.Info("Verifying backup folder")

	dirEntry, err := os.ReadDir(r.backupDir)
	if err != nil {
		return err
	}

	var labelsFileFound bool
	var importableCount int
	var dirs []string
	for _, entry := range dirEntry {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		name := entry.Name()
		if entry.IsDir() {
			if mailFolderRegExp.MatchString(name) {
				r.log.WithField("name", name).Info("Found a potential backup sub-folder")
				dirs = append(dirs, name)
			}
			continue
		}

		if strings.EqualFold(name, getLabelFileName()) {
			labelsFileFound = true
		}

		if !strings.HasSuffix(name, ".eml") {
			if !strings.HasSuffix(name, ".metadata.json") {
				r.log.WithField("fileName", name).Warn("Ignoring unknown file")
			}
			continue
		}

		jsonFile := strings.TrimSuffix(name, ".eml") + ".metadata.json"
		stats, err := os.Stat(filepath.Join(r.backupDir, jsonFile))
		if err != nil {
			r.log.WithError(err).WithField("jsonFile", jsonFile).Warn("EML file has no associated JSON file")
			continue
		}
		if stats.IsDir() {
			r.log.WithField("jsonFile", jsonFile).Warn("JSON file is directory")
			continue
		}
		importableCount++
	}

	if importableCount > 0 {
		if !labelsFileFound {
			return fmt.Errorf("the labels file '%v' could not be found", getLabelFileName())
		}

		r.log.WithField("mailCount", importableCount).Info("Importable emails found")
		return nil
	}

	if len(dirs) == 0 {
		return errors.New("no importable mail found")
	}

	if len(dirs) > 1 {
		return errors.New("the specified folder contains more than one backup sub-folder")
	}

	r.log.WithField("folderName", dirs[0]).Info("A potential backup sub-folder has been found and will be inspected")
	r.backupDir = filepath.Join(r.backupDir, dirs[0])
	return r.validateBackupDir()
}
