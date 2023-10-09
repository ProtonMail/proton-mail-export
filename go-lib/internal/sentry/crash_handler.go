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
	"runtime/pprof"

	"github.com/ProtonMail/gluon/async"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func NewPanicHandler(onRecover func()) async.PanicHandler {
	if !IsSentryEnabled() {
		return &async.NoopPanicHandler{}
	}

	return &sentryPanicHandler{
		onRecover: onRecover,
	}
}

type sentryPanicHandler struct {
	onRecover func()
	reporter  *sentryReporter
}

func (s *sentryPanicHandler) HandlePanic(r interface{}) {
	if r == nil {
		return
	}

	recoverErr := fmt.Errorf("recover: %v", r)
	if err := s.reporter.scopedSentryReport(nil, func() {
		if eventID := sentry.CaptureException(recoverErr); eventID != nil {
			logrus.WithError(recoverErr).WithField("reportID", eventID).Warn("Captured exception")
		}
	}); err != nil {
		logrus.WithError(err).Error("Failed to publish sentry crash report")
	}

	if err := pprof.Lookup("goroutine").WriteTo(logrus.StandardLogger().Writer(), 2); err != nil {
		logrus.WithError(err).Error("Failed to write crash report")
	}

	s.onRecover()
}
