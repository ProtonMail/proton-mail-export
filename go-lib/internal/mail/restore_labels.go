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
	"os"
	"path/filepath"

	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func (r *RestoreTask) restoreLabels() error {
	backupLabels, err := r.readLabelFile()
	if err != nil {
		return err
	}

	remoteLabels, err := r.session.GetClient().GetLabels(r.ctx, proton.LabelTypeFolder, proton.LabelTypeLabel, proton.LabelTypeSystem)
	if err != nil {
		return err
	}

	for _, label := range backupLabels {
		matchIndex := slices.IndexFunc(remoteLabels, func(remoteLabel proton.Label) bool {
			return (label.ID == remoteLabel.ID) || (label.Name == remoteLabel.Name)
		})
		if matchIndex != -1 {
			r.labelMapping[label.ID] = remoteLabels[matchIndex].ID
		} else {
			newLabel, err := r.recreateLabel(label)
			if err != nil {
				return err
			}
			r.labelMapping[label.ID] = newLabel.ID
			r.log.WithFields(logrus.Fields{"backupLabelID": label.ID, "remoteLabelID": newLabel.ID}).Info("Recreated remote folder")
		}
	}

	return nil
}

func (r *RestoreTask) readLabelFile() ([]proton.Label, error) {
	data, err := os.ReadFile(filepath.Join(r.backupDir, getLabelFileName()))
	if err != nil {
		return nil, err
	}

	versionedLabels, err := utils.NewVersionedJSON[[]proton.Label](LabelMetadataVersion, data)
	if err != nil {
		return nil, err
	}

	return versionedLabels.Payload, nil
}

func (r *RestoreTask) recreateLabel(label proton.Label) (proton.Label, error) {
	req := proton.CreateLabelReq{
		Name:     label.Name,
		Color:    label.Color,
		Type:     label.Type,
		ParentID: "",
	}

	return r.session.GetClient().CreateLabel(r.ctx, req)
}
