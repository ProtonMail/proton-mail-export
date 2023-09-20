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
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

func NewLogFileName() string {
	const format = "20060102_150405"
	return time.Now().Format(format) + "_export.log"
}

func NewLogFormatter() logrus.Formatter {
	return &logFormatter{}
}

func LogPrelude() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.
		WithField("appName", "Proton Export").
		WithField("version", ETVersionString).
		WithField("revision", ETRevision).
		WithField("build", ETBuildTime).
		WithField("runtime", runtime.GOOS).
		Info("Starting App")
}

type logFormatter struct{}

func (l logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// write time.
	if _, err := b.WriteString(entry.Time.Format(time.StampMilli)); err != nil {
		return nil, err
	}

	b.WriteByte('|')

	// write level.
	if _, err := b.Write(levelToBytes(entry.Level)); err != nil {
		return nil, err
	}

	b.WriteByte('|')
	b.WriteByte(' ')

	// write message.
	if _, err := b.WriteString(entry.Message); err != nil {
		return nil, err
	}
	b.WriteByte('\n')

	maxKeyLen := 0

	for _, f := range maps.Keys(entry.Data) {
		l := len(f)
		if l > maxKeyLen {
			maxKeyLen = l
		}
	}

	for f, v := range entry.Data {
		b.WriteByte('\t')
		if _, err := b.WriteString(fmt.Sprintf("%*s=", maxKeyLen, f)); err != nil {
			return nil, err
		}

		if _, err := b.WriteString(fmt.Sprint(v)); err != nil {
			return nil, err
		}

		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}

func levelToBytes(level logrus.Level) []byte {
	switch level {
	case logrus.TraceLevel:
		return []byte("trace")
	case logrus.DebugLevel:
		return []byte("debug")
	case logrus.InfoLevel:
		return []byte("info ")
	case logrus.WarnLevel:
		return []byte("warn ")
	case logrus.ErrorLevel:
		return []byte("error")
	case logrus.FatalLevel:
		return []byte("fatal")
	case logrus.PanicLevel:
		return []byte("panic")
	default:
		return []byte("unkno")
	}
}
