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

package main

/*
	#include <stdlib.h>
	typedef const char cchar_t;
*/
import "C"
import (
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/sirupsen/logrus"
)

//export etLogInit
func etLogInit(filePath *C.cchar_t) C.int {
	globalLogInstance.mutex.Lock()
	defer globalLogInstance.mutex.Unlock()

	if globalLogInstance.file != nil {
		return -1
	}

	path := C.GoString(filePath)

	if err := os.MkdirAll(path, 0o700); err != nil {
		globalLogInstance.lastError.Set(err)
		return -1
	}

	path = filepath.Join(path, internal.NewLogFileName())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		globalLogInstance.lastError.Set(err)
		return -1
	}

	globalLogInstance.clogPath = C.CString(path)

	logrus.SetOutput(file)
	logrus.SetFormatter(internal.NewLogFormatter())
	internal.LogPrelude()

	return 0
}

//export etLogClose
func etLogClose() {
	globalLogInstance.mutex.Lock()
	defer globalLogInstance.mutex.Unlock()

	if globalLogInstance.file != nil {
		logrus.SetOutput(os.Stdout)
		if err := globalLogInstance.file.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close log file")
		} else {
			globalLogInstance.file = nil
		}
	}

	if globalLogInstance.clogPath != nil {
		C.free(unsafe.Pointer(globalLogInstance.clogPath))
	}

	globalLogInstance.lastError.Close()
}

//export etLogGetLastError
func etLogGetLastError() *C.cchar_t {
	globalLogInstance.mutex.Lock()
	defer globalLogInstance.mutex.Unlock()

	return (*C.cchar_t)(globalLogInstance.lastError.GetErr())
}

//export etLogInfo
func etLogInfo(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Info(C.GoString(txt))
}

//export etLogDebug
func etLogDebug(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Debug(C.GoString(txt))
}

//export etLogWarn
func etLogWarn(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Warn(C.GoString(txt))
}

//export etLogError
func etLogError(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Error(C.GoString(txt))
}

//export etLogGetPath
func etLogGetPath() *C.cchar_t {
	globalLogInstance.mutex.Lock()
	defer globalLogInstance.mutex.Unlock()

	return globalLogInstance.clogPath
}

type globalLog struct {
	mutex     sync.Mutex
	file      *os.File
	lastError utils.CLastError
	clogPath  *C.char
}

//nolint:gochecknoglobals
var globalLogInstance globalLog
