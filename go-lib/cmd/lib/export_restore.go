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
#include "etrestore.h"
#include "etrestore_impl.h"
*/
import "C"
import (
	"context"
	"errors"
	"runtime/cgo"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
)

//export etSessionNewRestore
func etSessionNewRestore(sessionPtr *C.etSession, cRestorePath *C.cchar_t, outRestore **C.etRestore) C.etSessionStatus {
	cSession, ok := resolveSession(sessionPtr)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	defer async.HandlePanic(cSession.s.GetPanicHandler())

	if cSession.s.LoginState() != session.LoginStateLoggedIn {
		cSession.setLastError(session.ErrInvalidLoginState)
		return C.ET_SESSION_STATUS_ERROR
	}

	restorePath := C.GoString(cRestorePath)
	restoreTask, err := mail.NewRestoreTask(cSession.ctx, restorePath, cSession.s)
	if err != nil {
		cSession.setLastError(err)
		return C.ET_SESSION_STATUS_ERROR
	}

	h := internal.NewHandle(&cRestore{
		csession: cSession,
		restorer: restoreTask,
	})

	// Intentional misuse of unsafe pointer.
	//goland:noinspection GoVetUnsafePointer
	*outRestore = (*C.etRestore)(unsafe.Pointer(h)) //nolint:govet

	return C.ET_SESSION_STATUS_OK
}

//export etRestoreDelete
func etRestoreDelete(ptr *C.etRestore) C.etRestoreStatus {
	h := restorePtrToHandle(ptr)

	s, ok := h.resolve()
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(s.csession.s.GetPanicHandler())

	s.restorer.Close()
	s.lastError.Close()

	h.Delete()

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreStart
func etRestoreStart(ptr *C.etRestore, callbacks *C.etRestoreCallbacks) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	reporter := &restoreReporter{
		restorer:  ce.restorer,
		callbacks: callbacks,
	}

	ce.csession.s.GetTelemetryService().SendRestoreStart()
	startTime := time.Now()

	err := ce.restorer.Run(reporter)

	ce.csession.s.GetTelemetryService().SendRestoreFinished(
		ce.restorer.GetOperationCancelledByUser(),
		err != nil,
		int(time.Since(startTime).Seconds()),
		int(ce.restorer.GetImportableCount()),
		int(ce.restorer.GetFailedCount()),
		int(ce.restorer.GetImportedCount()),
	)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return C.ET_RESTORE_STATUS_CANCELLED
		}

		ce.lastError.Set(internal.MapError(err))
		return C.ET_RESTORE_STATUS_ERROR
	}

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreCancel
func etRestoreCancel(ptr *C.etRestore) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	ce.restorer.Cancel()

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreGetLastError
func etRestoreGetLastError(ptr *C.etRestore) *C.cchar_t {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return nil
	}

	return (*C.cchar_t)(ce.lastError.GetErr())
}

//export etRestoreGetBackupPath
func etRestoreGetBackupPath(ptr *C.etRestore, outPath **C.char) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*outPath = C.CString(ce.restorer.GetBackupPath())

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreGetImportableCount
func etRestoreGetImportableCount(ptr *C.etRestore, count *C.int64_t) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*count = C.int64_t(ce.restorer.GetImportableCount())

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreGetImportedCount
func etRestoreGetImportedCount(ptr *C.etRestore, count *C.int64_t) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*count = C.int64_t(ce.restorer.GetImportedCount())

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreGetFailedCount
func etRestoreGetFailedCount(ptr *C.etRestore, count *C.int64_t) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*count = C.int64_t(ce.restorer.GetFailedCount())

	return C.ET_RESTORE_STATUS_OK
}

//export etRestoreGetSkippedCount
func etRestoreGetSkippedCount(ptr *C.etRestore, count *C.int64_t) C.etRestoreStatus {
	ce, ok := resolveRestore(ptr)
	if !ok {
		return C.ET_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*count = C.int64_t(ce.restorer.GetSkippedCount())

	return C.ET_RESTORE_STATUS_OK
}

type cRestore struct {
	csession  *csession
	restorer  *mail.RestoreTask
	lastError utils.CLastError
}

type RestoreHandle struct {
	internal.Handle
}

func (h RestoreHandle) resolve() (*cRestore, bool) {
	return internal.ResolveHandle[cRestore](h.Handle)
}

func restorePtrToHandle(ptr *C.etRestore) RestoreHandle {
	return RestoreHandle{Handle: cgo.Handle(unsafe.Pointer(ptr))}
}

func resolveRestore(ptr *C.etRestore) (*cRestore, bool) {
	h := restorePtrToHandle(ptr)

	return h.resolve()
}

type restoreReporter struct {
	totalMessageCount   atomic.Uint64
	currentMessageCount atomic.Uint64
	callbacks           *C.etRestoreCallbacks
	restorer            *mail.RestoreTask
}

func (m *restoreReporter) SetMessageTotal(total uint64) {
	m.totalMessageCount.Store(total)
}

func (m *restoreReporter) SetMessageProcessed(total uint64) {
	m.currentMessageCount.Store(total)
}

func (m *restoreReporter) OnProgress(delta int) {
	newMessageCount := m.currentMessageCount.Add(uint64(delta))

	var progress float32
	totalMessageCount := m.totalMessageCount.Load()
	if totalMessageCount != 0 {
		progress = float32(float64(newMessageCount) / float64(totalMessageCount) * 100.0)
	} else {
		progress = float32(0.0)
	}

	C.etRestoreCallbackOnProgress(m.callbacks, C.float(progress))
}
