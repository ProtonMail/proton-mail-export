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

// #cgo CFLAGS: -I "cgo_headers" -D "ET_CGO=1"
/*
#include "etglobal.h"
#include "etglobal_impl.h"
*/
import "C"
import (
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/sentry"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/sirupsen/logrus"
)

//export etInit
func etInit(filePath *C.cchar_t, onRecover C.etOnRecoverFn) C.int {
	etGlobalState.mutex.Lock()
	defer etGlobalState.mutex.Unlock()

	if err := sentry.InitSentry(); err != nil {
		etGlobalState.lastError.Set(err)
		return -1
	}

	if etGlobalState.file != nil {
		return -1
	}

	path := C.GoString(filePath)

	if err := os.MkdirAll(path, 0o700); err != nil {
		etGlobalState.lastError.Set(err)
		return -1
	}

	path = filepath.Join(path, internal.NewLogFileName())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		etGlobalState.lastError.Set(err)
		return -1
	}

	etGlobalState.clogPath = C.CString(path)

	logrus.SetOutput(file)
	logrus.SetFormatter(internal.NewLogFormatter())
	internal.LogPrelude()

	if onRecover != nil {
		etGlobalState.onRecoverCB = func() {
			C.etCallOnRecover(onRecover)
		}
	} else {
		etGlobalState.onRecoverCB = func() {
			os.Exit(-200)
		}
	}

	return 0
}

//export etGetLastError
func etGetLastError() *C.cchar_t {
	etGlobalState.mutex.Lock()
	defer etGlobalState.mutex.Unlock()

	return (*C.cchar_t)(etGlobalState.lastError.GetErr())
}

//export etClose
func etClose() {
	etGlobalState.mutex.Lock()
	defer etGlobalState.mutex.Unlock()

	if etGlobalState.file != nil {
		logrus.SetOutput(os.Stdout)
		if err := etGlobalState.file.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close log file")
		} else {
			etGlobalState.file = nil
		}
	}

	if etGlobalState.clogPath != nil {
		C.free(unsafe.Pointer(etGlobalState.clogPath))
	}

	etGlobalState.lastError.Close()
}

type globalState struct {
	mutex       sync.Mutex
	file        *os.File
	lastError   utils.CLastError
	clogPath    *C.char
	onRecoverCB func()
}

//nolint:gochecknoglobals
var etGlobalState globalState

func etOnRecoverCB() func() {
	return etGlobalState.onRecoverCB
}
