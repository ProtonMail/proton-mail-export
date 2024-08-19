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

var errCircularLabelReference = errors.New("unable to sort labels because of a circular reference")

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

		labelID, name := matchLocalLabelWithRemote(label, remoteLabels)
		if len(labelID) > 0 {
			r.labelMapping[label.ID] = labelID
			continue
		}

		label.Name = name
		if err = r.createAndMapLabel(label); err != nil {
			return err
		}
	}

	return nil
}

// matchLocalLabelWithRemote match a label from a backup with remote labels.
// if a label of the same name and type is found, its labelID and returned, otherwise and empty labelID and the name of the matching
// label to create on the server is returned as newName. The name of the label to create will be the name of the label in the backup, unless
// a label of the same name but different type already exists, in which case a number in appended at the end of the name.
func matchLocalLabelWithRemote(label proton.Label, remoteLabels []proton.Label) (labelID, newName string) {
	if isSystemLabel(label.ID) {
		return label.ID, ""
	}

	index := slices.IndexFunc(remoteLabels, func(remoteLabel proton.Label) bool {
		return (label.ID == remoteLabel.ID) || strings.EqualFold(label.Name, remoteLabel.Name)
	})

	// label does not exist.
	if index == -1 {
		return "", label.Name
	}

	// label exists remotely and is of the correct type. We map it.
	if remoteLabels[index].Type == label.Type {
		return remoteLabels[index].ID, ""
	}

	// label exists remotely but not of the right type, we need a new name
	return "", findFirstAvailableLabelIncrementalName(label.Name, remoteLabels)
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

// createAndMapLabel create the given label and adds the mapping of its remote ID to r.labelMappings.
func (r *RestoreTask) createAndMapLabel(label proton.Label) error {
	var remoteParentID string
	if len(label.ParentID) > 0 {
		var ok bool
		remoteParentID, ok = r.labelMapping[label.ParentID]
		if !ok {
			// this should not happen has we have sorted the label beforehand.
			return fmt.Errorf("could not find parent label for %s", label.ID)
		}
	}

	newLabel, err := r.session.GetClient().CreateLabel(
		r.ctx,
		proton.CreateLabelReq{
			Name:     label.Name,
			Color:    label.Color,
			Type:     label.Type,
			ParentID: remoteParentID,
		},
	)
	if err != nil {
		return err
	}

	r.labelMapping[label.ID] = newLabel.ID
	r.log.WithFields(logrus.Fields{"backupLabelID": label.ID, "remoteLabelID": newLabel.ID}).Info("Recreated remote label")
	return nil
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

// findFirstAvailableLabelIncrementalName return a variant of name that is not in the remoteLabels list, ensuring uniqueness by append
// a number between parenthesis to the name. For instance if the list contains 'Folder' and 'Folder (1)', the function will return 'Folder (2)'.
func findFirstAvailableLabelIncrementalName(name string, remoteLabels []proton.Label) string {
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s (%d)", name, i)
		index := slices.IndexFunc(remoteLabels, func(label proton.Label) bool { return strings.EqualFold(label.Name, candidate) })
		if index == -1 {
			return candidate
		}
	}
}
