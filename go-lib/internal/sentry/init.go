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

package sentry

import (
	"fmt"
	"log"
	"sync"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/hv"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var initSentryOnce sync.Once

func InitSentry() error {
	var err error

	initSentryOnce.Do(func() {
		err = initSentryImpl()
	})

	return err
}

func IsSentryEnabled() bool {
	return len(internal.ETSentryDNS) != 0
}

func initSentryImpl() error {
	// Do not init sentry if no url is specified.
	if !IsSentryEnabled() {
		return nil
	}

	hostname, err := hv.GetProtectedHostname()
	if err != nil {
		logrus.WithError(err).Error("Failed to get hostname")
		hostname = "Unknown"
	}

	options := sentry.ClientOptions{
		Dsn:            internal.ETSentryDNS,
		Transport:      sentry.NewHTTPSyncTransport(),
		Release:        internal.ETAppIdentifier,
		MaxBreadcrumbs: 50,
		ServerName:     hostname,
	}

	if err := sentry.Init(options); err != nil {
		return fmt.Errorf("failed to init sentry: %w", err)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetFingerprint([]string{"{{ default }}"})
		scope.SetUser(sentry.User{ID: hostname})
	})

	sentry.Logger = log.New(
		logrus.WithField("sentry", "sentry-go").WriterLevel(logrus.WarnLevel),
		"",
		0,
	)

	return nil
}
