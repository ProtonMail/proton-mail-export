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
	"runtime"
	"time"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/hv"
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func NewReporter() reporter.Reporter {
	if !IsSentryEnabled() {
		return &reporter.NullReporter{}
	}

	return newReporter()
}

type sentryReporter struct {
	arch string
}

func (r *sentryReporter) ReportMessage(s string, context reporter.Context) {
	if err := r.scopedSentryReport(context, func() {
		if eventID := sentry.CaptureMessage(s); eventID != nil {
			logrus.WithField("msg", s).WithField("eventID", eventID).Warn("Captured message")
		}
	}); err != nil {
		logrus.WithError(err).WithField("msg", s).Error("Failed to report message to sentry")
	}
}

func (r *sentryReporter) ReportError(a any, context reporter.Context) {
	err := fmt.Errorf("error: %v", a)
	if err := r.scopedSentryReport(context, func() {
		if eventID := sentry.CaptureException(err); eventID != nil {
			logrus.WithError(err).WithField("eventID", eventID).Warn("Captured error")
		}
	}); err != nil {
		logrus.WithError(err).WithError(err).Error("Failed to report error to sentry")
	}
}

func newReporter() *sentryReporter {
	return &sentryReporter{
		arch: hv.GetHostArch(),
	}
}

func (r *sentryReporter) scopedSentryReport(context reporter.Context, do func()) error {
	tags := map[string]string{
		"OS":       runtime.GOOS,
		"Version":  internal.ETVersionString,
		"HostArch": r.arch,
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
		if len(context) != 0 {
			scope.SetContexts(
				map[string]sentry.Context{"go-export": contextToString(context)},
			)
		}

		do()
	})

	if !sentry.Flush(time.Second * 10) {
		return fmt.Errorf("failed to report sentry error")
	}

	return nil
}

func contextToString(context sentry.Context) sentry.Context {
	res := make(sentry.Context)

	for k, v := range context {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}
