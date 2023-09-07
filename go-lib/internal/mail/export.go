// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is Free software: you can redistribute it and/or modify
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
	"context"
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/xslices"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const NumParallelDownloads = 10
const NumParallelBuilders = 4
const NumParallelWriters = 4
const MetadataPageSize = 64

// Mail Exports will be created in the given directory and will be structured:
// <email>
//  |- mail
//	    |- export.log
//      |- labels.json
//      |- msg-id.eml
//      |- msg-id.meta.json

type ExportTask struct {
	group     *async.Group
	tmpDir    string
	exportDir string
	session   *session.Session
	log       *logrus.Entry
}

func NewExportTask(
	ctx context.Context,
	tmpDir, exportPath string,
	session *session.Session,
) *ExportTask {
	return &ExportTask{
		group:     async.NewGroup(ctx, session.GetPanicHandler()),
		tmpDir:    tmpDir,
		exportDir: exportPath,
		session:   session,
		log:       logrus.WithField("export", "mail"),
	}
}

type Reporter interface {
	StageProgressReporter
	StageErrorReporter
}

func (e *ExportTask) Run(ctx context.Context, reporter Reporter) error {
	// GODT-2900: Handle network errors/loss.
	client := e.session.GetClient()
	// Get user info
	user, err := client.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to load user info: %w", err)
	}

	// Get User addresses
	addresses, err := client.GetAddresses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user addresses: %w", err)
	}

	keyRing, err := apiclient.NewUnlockedKeyRing(&user, addresses, e.session.GetMailboxPassword())
	if err != nil {
		return fmt.Errorf("failed to unlock user keyring:%w", err)
	}
	defer keyRing.Close()

	exportDir := filepath.Join(e.exportDir, user.Email, "mail")

	if err := os.MkdirAll(exportDir, 0o700); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	// Create required folders
	if err := e.WriteLabelMetadata(ctx, e.tmpDir, exportDir); err != nil {
		return err
	}

	// get message count
	msgCountPerLabel, err := client.GetGroupedMessageCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get message count: %w", err)
	}

	var totalMessageCount uint64
	var foundAllMailLabel bool

	for _, c := range msgCountPerLabel {
		if c.LabelID == proton.AllMailLabel {
			totalMessageCount = uint64(c.Total)
			foundAllMailLabel = true
			break
		}
	}

	reporter.SetMessageTotal(totalMessageCount)

	e.log.Infof("Found %v Messages for download", totalMessageCount)

	if !foundAllMailLabel {
		return fmt.Errorf("failed to determine total message count")
	}

	// Build stages
	metaStage := NewMetadataStage(client, e.log, MetadataPageSize)
	downloadStage := NewDownloadStage(client, NumParallelDownloads, e.log, e.session.GetPanicHandler())
	buildStage := NewBuildStage(NumParallelBuilders, e.log, e.session.GetPanicHandler())
	writeStage := NewWriteStage(e.tmpDir, e.exportDir, NumParallelWriters, e.log, reporter, e.session.GetPanicHandler())

	// GODT-2900: Handle network errors/loss.
	errReporter := &NullErrorReporter{}

	// Start pipeline.
	e.group.Once(func(ctx context.Context) {
		metaStage.Run(ctx, errReporter)
	})
	e.group.Once(func(ctx context.Context) {
		downloadStage.Run(ctx, metaStage.outputCh, errReporter)
	})
	e.group.Once(func(ctx context.Context) {
		buildStage.Run(ctx, downloadStage.outputCh, keyRing, errReporter)
	})
	e.group.Once(func(ctx context.Context) {
		writeStage.Run(ctx, buildStage.outputCh, errReporter)
	})

	// wait for downloads to finish.
	e.group.WaitToFinish()
	// collect errors.

	return nil
}

func (e *ExportTask) WriteLabelMetadata(ctx context.Context, tmpDir, exportPath string) error {
	// GODT-2925 version metadata.
	apiLabels, err := e.session.GetClient().GetLabels(ctx, proton.LabelTypeSystem, proton.LabelTypeFolder, proton.LabelTypeLabel)
	if err != nil {
		return fmt.Errorf("failed to retrieve labels: %w", err)
	}

	apiLabels = xslices.Filter(apiLabels, func(t proton.Label) bool {
		return wantLabel(t)
	})

	labelData, err := json.Marshal(apiLabels)
	if err != nil {
		return fmt.Errorf("failed to json encode labels: %w", err)
	}

	labelFile := filepath.Join(exportPath, getLabelFileName())

	return utils.WriteFileSafe(tmpDir, labelFile, labelData, &utils.Sha256IntegrityChecker{})
}

func getLabelFileName() string {
	return "labels.json"
}
