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
	"context"
	"path/filepath"
	"regexp"

	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/sirupsen/logrus"
)

var mailFolderRegExp = regexp.MustCompile(`^mail_\d{8}_\d{6}$`)

type RestoreTask struct {
	ctx          context.Context
	ctxCancel    func()
	backupDir    string
	session      *session.Session
	log          *logrus.Entry
	labelMapping map[string]string // map of [backup labelIDs] to remoteLabelIDs
}

func NewRestoreTask(ctx context.Context, backupDir string, session *session.Session) (*RestoreTask, error) {
	absPath, err := filepath.Abs(backupDir)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &RestoreTask{
		ctx:          ctx,
		ctxCancel:    cancel,
		backupDir:    absPath,
		session:      session,
		log:          logrus.WithField("backup", "mail").WithField("userID", session.GetUser().ID),
		labelMapping: make(map[string]string),
	}, nil
}

func (r *RestoreTask) Run(_ Reporter) error {
	defer r.log.Info("Finished")
	r.log.WithField("backupDir", r.backupDir).Info("Starting")

	if err := r.validateBackupDir(); err != nil {
		return err
	}

	if err := r.restoreLabels(); err != nil {
		return err
	}

	return r.importMails()
}
