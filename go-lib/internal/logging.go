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

package internal

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

func NewLogFileName() string {
	const format = "20060102_150405"
	return time.Now().Format(format) + "_export.log"
}

func NewLogFormatter() logrus.Formatter {
	return &logrus.TextFormatter{
		ForceColors:            false,
		DisableColors:          true,
		ForceQuote:             false,
		DisableTimestamp:       false,
		FullTimestamp:          false,
		TimestampFormat:        time.StampMilli,
		DisableSorting:         false,
		SortingFunc:            nil,
		DisableLevelTruncation: false,
		PadLevelText:           true,
		QuoteEmptyFields:       false,
		FieldMap:               nil,
		CallerPrettyfier:       nil,
	}
}

func LogPrelude() {
	logrus.
		WithField("appName", "Proton Export").
		WithField("version", ETVersionString).
		WithField("revision", ETRevision).
		WithField("build", ETBuildTime).
		WithField("runtime", runtime.GOOS).
		Info("Starting App")
}
