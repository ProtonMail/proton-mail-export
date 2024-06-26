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
	"time"

	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var errCircularLabelReference = errors.New("unable to sort labels because of a circular reference1")

func (r *RestoreTask) restoreLabels() error {
	backupLabels, err := r.readLabelFile()
	if err != nil {
		return err
	}

	backupLabels, err = sortLabels(backupLabels)
	if err != nil {
		return err
	}

	remoteLabels, err := r.session.GetClient().GetLabels(r.ctx, proton.LabelTypeFolder, proton.LabelTypeLabel, proton.LabelTypeSystem)
	if err != nil {
		return err
	}

	for _, label := range backupLabels {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		matchIndex := slices.IndexFunc(remoteLabels, func(remoteLabel proton.Label) bool {
			return (label.ID == remoteLabel.ID) || strings.EqualFold(label.Name, remoteLabel.Name)
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
	var remoteParentID string
	if len(label.ParentID) > 0 {
		var ok bool
		remoteParentID, ok = r.labelMapping[label.ParentID]
		if !ok {
			// this should not happen has we have sorted the label beforehand.
			return proton.Label{}, fmt.Errorf("could not find parent label for %s", label.ID)
		}
	}

	return r.session.GetClient().CreateLabel(
		r.ctx,
		proton.CreateLabelReq{
			Name:     label.Name,
			Color:    label.Color,
			Type:     label.Type,
			ParentID: remoteParentID,
		},
	)
}

func (r *RestoreTask) createImportLabel() error {
	label, err := r.session.GetClient().CreateLabel(
		r.ctx,
		proton.CreateLabelReq{
			Name:     "Import " + time.Now().Format("2006-01-02 15:04:05"),
			Color:    "#f66",
			Type:     proton.LabelTypeLabel,
			ParentID: "",
		},
	)

	if err != nil {
		return err
	}

	r.importLabelID = label.ID
	return nil
}

// sortLabels Sorts the labels ensuring that parent labels are listed before their children.
func sortLabels(labels []proton.Label) ([]proton.Label, error) {
	result := make([]proton.Label, 0, len(labels))
	remaining := make([]proton.Label, 0)

	for _, label := range labels {
		// if a label has no parent, no problem
		if len(label.ParentID) == 0 {
			result = append(result, label)
			continue
		}

		// if the parent is already in the result list, no problem, otherwise we queue it for later processing.
		if slices.ContainsFunc(result, func(sortedLabel proton.Label) bool {
			return sortedLabel.ID == label.ParentID
		}) {
			result = append(result, label)
			continue
		}
		remaining = append(remaining, label)
	}

	// we retry the remaining labels until everything is processed.
	// if we cannot processed at least one remaining item per loop, we're stuck. It should not happen.
	for len(remaining) > 0 {
		var newRemaining []proton.Label
		for _, label := range remaining {
			if slices.ContainsFunc(result, func(sortedLabel proton.Label) bool {
				return sortedLabel.ID == label.ParentID
			}) {
				result = append(result, label)
				continue
			}

			newRemaining = append(newRemaining, label)
		}

		if len(remaining) == len(newRemaining) {
			return nil, errCircularLabelReference
		}

		remaining = newRemaining
	}

	return result, nil
}
