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
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/sirupsen/logrus"
)

var mailFolderRegExp = regexp.MustCompile(`^mail_\d{8}_\d{6}$`)

type RestoreTask struct {
	ctx           context.Context
	ctxCancel     func()
	addrKR        *apiclient.UnlockedKeyRing
	addrID        string
	backupDir     string
	session       *session.Session
	log           *logrus.Entry
	labelMapping  map[string]string // map of [backup labelIDs] to remoteLabelIDs
	importLabelID string
}

func NewRestoreTask(ctx context.Context, backupDir string, session *session.Session) (*RestoreTask, error) {
	absPath, err := filepath.Abs(backupDir)
	if err != nil {
		return nil, err
	}

	addrID, err := getPrimaryAddressID(ctx, session)
	if err != nil {
		return nil, err
	}

	log := logrus.WithField("backup", "mail").WithField("userID", session.GetUser().ID)

	addrKR, err := getUnlockedAddressKeyRing(ctx, session)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &RestoreTask{
		ctx:          ctx,
		ctxCancel:    cancel,
		addrKR:       addrKR,
		addrID:       addrID,
		backupDir:    absPath,
		session:      session,
		log:          log,
		labelMapping: make(map[string]string),
	}, nil
}

func (r *RestoreTask) Teardown() {
	r.addrKR.Close()
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

	if err := r.createImportLabel(); err != nil {
		return err
	}

	return r.importMails()
}

func getPrimaryAddressID(ctx context.Context, session *session.Session) (string, error) {
	addresses, err := session.GetClient().GetAddresses(ctx)

	if err != nil {
		return "", err
	}

	if len(addresses) == 0 {
		return "", errors.New("address list is empty")
	}

	return addresses[0].ID, nil
}

func getUnlockedAddressKeyRing(ctx context.Context, session *session.Session) (*apiclient.UnlockedKeyRing, error) {
	client := session.GetClient()
	user := session.GetUser()
	salts := session.GetUserSalts()

	saltedKeyPass, err := salts.SaltForKey(session.GetMailboxPassword(), user.Keys.Primary().ID)
	if err != nil {
		return nil, fmt.Errorf("failed to salt key password: %w", err)
	}

	if userKR, err := user.Keys.Unlock(saltedKeyPass, nil); err != nil {
		return nil, fmt.Errorf("failed to unlock user keys: %w", err)
	} else if userKR.CountDecryptionEntities() == 0 {
		return nil, fmt.Errorf("failed to unlock user keys")
	}

	// Get User addresses
	addresses, err := client.GetAddresses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	addrKR, err := apiclient.NewUnlockedKeyRing(user, addresses, saltedKeyPass)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock user keyring:%w", err)
	}

	return addrKR, nil
}
