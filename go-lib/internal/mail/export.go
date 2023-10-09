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
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/xslices"
	"github.com/pbnjay/memory"
	"github.com/sirupsen/logrus"
)

const NumParallelDownloads = 10
const NumParallelBuilders = 4
const NumParallelWriters = 4
const MetadataPageSize = 64
const MB = 1024 * 1024
const MinDownloadMemMB = 128 * MB
const MinBuildMemMB = 128 * MB
const MaxDownloadMemMB = 1024 * MB
const MaxBuildMemMB = 512 * MB

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
	exportPath string,
	session *session.Session,
) *ExportTask {
	exportPath = filepath.Join(exportPath, "mail")

	// Tmp dir needs to be next to export path to as os.rename doesn't work if export path is on a different volume.
	tmpDir := filepath.Join(exportPath, "temp")

	return &ExportTask{
		group:     async.NewGroup(ctx, session.GetPanicHandler()),
		tmpDir:    tmpDir,
		exportDir: exportPath,
		session:   session,
		log:       logrus.WithField("export", "mail").WithField("email", session.GetEmail()),
	}
}

type Reporter interface {
	StageProgressReporter
}

func (e *ExportTask) Close() {
	e.group.CancelAndWait()

	if err := os.RemoveAll(e.tmpDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			e.log.WithError(err).Error("Failed to remove temp directory")
		}
	}
}

func (e *ExportTask) Cancel() {
	e.group.Cancel()
}

func (e *ExportTask) GetRequiredDiskSpaceEstimate(ctx context.Context) (uint64, error) {
	user, err := e.session.GetClient().GetUser(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to load user info: %w", err)
	}

	return approximateDiskUsage(user.ProductUsedSpace.Mail), nil
}

func (e *ExportTask) Run(ctx context.Context, reporter Reporter) error {
	defer e.log.Info("Finished")
	e.log.WithFields(logrus.Fields{"tmp-dir": e.tmpDir, "export-dir": e.exportDir}).Info("Starting")

	reporter.OnProgress(0)

	e.log.Debug("Preparing export dir")

	if err := os.MkdirAll(e.exportDir, 0o700); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	if err := os.MkdirAll(e.tmpDir, 0o700); err != nil {
		return fmt.Errorf("failed to create export tmp directory: %w", err)
	}

	e.log.Debug("Getting user info")
	client := e.session.GetClient()
	// Get user info
	user, err := client.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to load user info: %w", err)
	}

	e.log.Infof(
		"Reported space usage %v MB, estimated disk uage %v MB",
		toMB(user.ProductUsedSpace.Mail),
		toMB(approximateDiskUsage(user.ProductUsedSpace.Mail)),
	)

	salts, err := client.GetSalts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get key salts: %w", err)
	}

	saltedKeyPass, err := salts.SaltForKey(e.session.GetMailboxPassword(), user.Keys.Primary().ID)
	if err != nil {
		return fmt.Errorf("failed to salt key password: %w", err)
	}

	e.log.Debug("Unlocking decryption key")
	if userKR, err := user.Keys.Unlock(saltedKeyPass, nil); err != nil {
		return fmt.Errorf("failed to unlock user keys: %w", err)
	} else if userKR.CountDecryptionEntities() == 0 {
		return fmt.Errorf("failed to unlock user keys")
	}

	e.log.Debug("Getting addresses")
	// Get User addresses
	addresses, err := client.GetAddresses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user addresses: %w", err)
	}

	e.log.Debug("Unlocking address keys")
	keyRing, err := apiclient.NewUnlockedKeyRing(&user, addresses, saltedKeyPass)
	if err != nil {
		return fmt.Errorf("failed to unlock user keyring:%w", err)
	}
	defer keyRing.Close()

	// Create required folders
	if err := e.WriteLabelMetadata(ctx, e.tmpDir, e.exportDir); err != nil {
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

	totalMemory := memory.TotalMemory()

	var (
		buildMemMB    uint64
		downloadMemMb uint64
	)

	if totalMemory >= 4096*MB {
		buildMemMB = MaxBuildMemMB
		downloadMemMb = MaxDownloadMemMB
	} else {
		buildMemMB = MinBuildMemMB
		downloadMemMb = MinDownloadMemMB
	}

	// Build stages
	metaStage := NewMetadataStage(client, e.log, MetadataPageSize)
	downloadStage := NewDownloadStage(client, NumParallelDownloads, e.log, downloadMemMb, e.session.GetPanicHandler())
	buildStage := NewBuildStage(NumParallelBuilders, e.log, buildMemMB, e.session.GetPanicHandler(), e.session.GetReporter(), user.ID)
	writeStage := NewWriteStage(e.tmpDir, e.exportDir, NumParallelWriters, e.log, reporter, e.session.GetPanicHandler())

	e.log.Debug("Starting message download")
	errReporter := &exportErrReporter{
		export: e,
		lock:   sync.Mutex{},
		errors: nil,
	}

	// start pipeline.
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

	e.log.Debug("Message download finished")

	// collect errors.
	exportError := errReporter.getErrors()
	if len(exportError) == 0 {
		return nil
	}

	e.log.Error("Export task ran into the following errors")
	for i, err := range exportError {
		e.log.WithError(err).Errorf("Error %v", i)
	}

	return exportError[0]
}

const LabelMetadataVersion = 1

func (e *ExportTask) WriteLabelMetadata(ctx context.Context, tmpDir, exportPath string) error {
	e.log.Debug("Writing root label metadata")
	apiLabels, err := e.session.GetClient().GetLabels(ctx, proton.LabelTypeSystem, proton.LabelTypeFolder, proton.LabelTypeLabel)
	if err != nil {
		return fmt.Errorf("failed to retrieve labels: %w", err)
	}

	apiLabels = xslices.Filter(apiLabels, wantLabel)

	labelData, err := utils.GenerateVersionedJSON(LabelMetadataVersion, apiLabels)
	if err != nil {
		return fmt.Errorf("failed to json encode labels: %w", err)
	}

	labelFile := filepath.Join(exportPath, getLabelFileName())

	return utils.WriteFileSafe(tmpDir, labelFile, labelData, &utils.Sha256IntegrityChecker{})
}

func (e *ExportTask) GetExportPath() string {
	return e.exportDir
}

func getLabelFileName() string {
	return "labels.json"
}

type exportErrReporter struct {
	export *ExportTask
	lock   sync.Mutex
	errors []error
}

func (e *exportErrReporter) ReportStageError(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if len(e.errors) == 0 {
		e.export.log.Debug("Cancelling context due to error")
		e.export.group.Cancel()
	}
	e.errors = append(e.errors, err)
}

func (e *exportErrReporter) getErrors() []error {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.errors
}

func approximateDiskUsage(v uint64) uint64 {
	// add another 30% of to current usage estimate due to variance in the decrypted message sizes and metadata.
	return uint64(math.Ceil(float64(v) * 1.3))
}

func toMB(v uint64) uint64 {
	return v / 1024 / 1024
}
