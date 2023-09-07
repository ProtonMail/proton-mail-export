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
	"context"
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/parallel"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type WriteStage struct {
	tempPath         string
	dirPath          string
	panicHandler     async.PanicHandler
	log              *logrus.Entry
	progressReporter StageProgressReporter
	parallelWriters  int
}

func NewWriteStage(
	tempPath string,
	dirPath string,
	parallelWriters int,
	log *logrus.Entry,
	progressReporter StageProgressReporter,
	panicHandler async.PanicHandler,
) *WriteStage {
	return &WriteStage{
		tempPath:         tempPath,
		dirPath:          dirPath,
		panicHandler:     panicHandler,
		parallelWriters:  parallelWriters,
		progressReporter: progressReporter,
		log:              log.WithField("stage", "write"),
	}
}

func (w *WriteStage) Run(ctx context.Context, inputs <-chan BuildStageOutput, errReporter StageErrorReporter) {
	w.log.Debug("Starting")
	defer w.log.Debug("Exiting")

	for input := range inputs {
		if ctx.Err() != nil {
			return
		}

		if err := parallel.DoContext(ctx, w.parallelWriters, len(input.messages), func(ctx context.Context, i int) error {
			emlBytes := input.result[i].eml.Bytes()
			// GODT-2915 handle messages that failed to decrypt.
			// GODT-2916 handle messages that failed to build.
			if len(emlBytes) == 0 {
				w.log.Errorf("Not yet implemented")
				return nil
			}

			filePath := filepath.Join(w.dirPath, input.messages[i].ID)
			metadataPath := filePath + "metadata.json"
			filePath += ".eml"

			metadataBytes, err := generateMetadataBytes(&input.messages[i])
			if err != nil {
				w.log.WithField("msg-id", input.messages[i].ID).WithError(err).Error("Failed to generate metadata")
				return fmt.Errorf("failed to generate message metadata: %w", err)
			}

			integrityChecker := &utils.Sha256IntegrityChecker{}

			if err := utils.WriteFileSafe(w.tempPath, filePath, emlBytes, integrityChecker); err != nil {
				w.log.WithField("msg-id", input.messages[i].ID).WithError(err).Errorf("Failed to write %v", filePath)
				return fmt.Errorf("failed to write '%v': %w", filePath, err)
			}

			if err := utils.WriteFileSafe(w.tempPath, metadataPath, metadataBytes, integrityChecker); err != nil {
				w.log.WithField("msg-id", input.messages[i].ID).WithError(err).Errorf("Failed to write metadata file %v", metadataPath)
				return fmt.Errorf("failed to write metadata '%v': %w", metadataPath, err)
			}

			return nil
		}); err != nil {
			errReporter.ReportStageError(err)
			return
		}

		w.progressReporter.OnProgress(len(input.messages))
	}
}

func generateMetadataBytes(msg *proton.FullMessage) ([]byte, error) {
	// GODT-2925 version metadata.
	return json.Marshal(msg.MessageMetadata)
}
