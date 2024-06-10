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
#include "etexport_restore.h"
#include "etexport_restore_impl.h"
*/
import "C"
import (
	"context"
	"errors"
	"path/filepath"
	"runtime/cgo"
	"sync/atomic"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
)

//export etSessionNewExportRestore
func etSessionNewExportRestore(sessionPtr *C.etSession, cRestorePath *C.cchar_t, outExportRestore **C.etExportRestore) C.etSessionStatus {
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
	restorePath = filepath.Join(restorePath, cSession.s.GetUser().Email)

	restoreTask, err := mail.NewRestoreTask(cSession.ctx, restorePath, cSession.s)
	if err != nil {
		cSession.setLastError(err)
		return C.ET_SESSION_STATUS_ERROR
	}

	h := internal.NewHandle(&cExportRestore{
		csession: cSession,
		restorer: restoreTask,
	})

	// Intentional misuse of unsafe pointer.
	//goland:noinspection GoVetUnsafePointer
	*outExportRestore = (*C.etExportRestore)(unsafe.Pointer(h)) //nolint:govet

	return C.ET_SESSION_STATUS_OK
}

//export etExportRestoreDelete
func etExportRestoreDelete(ptr *C.etExportRestore) C.etExportRestoreStatus {
	h := exportRestorePtrToHandle(ptr)

	s, ok := h.resolve()
	if !ok {
		return C.ET_EXPORT_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(s.csession.s.GetPanicHandler())

	s.restorer.Close()
	s.lastError.Close()

	h.Delete()

	return C.ET_EXPORT_RESTORE_STATUS_OK
}

//export etExportRestoreStart
func etExportRestoreStart(ptr *C.etExportRestore, callbacks *C.etExportRestoreCallbacks) C.etExportRestoreStatus {
	ce, ok := resolveExportRestore(ptr)
	if !ok {
		return C.ET_EXPORT_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	reporter := &restoreReporter{
		restorer:  ce.restorer,
		callbacks: callbacks,
	}

	if err := ce.restorer.Run(reporter); err != nil {
		if errors.Is(err, context.Canceled) {
			return C.ET_EXPORT_RESTORE_STATUS_CANCELLED
		}

		ce.lastError.Set(internal.MapError(err))
		return C.ET_EXPORT_RESTORE_STATUS_ERROR
	}

	return C.ET_EXPORT_RESTORE_STATUS_OK
}

//export etExportRestoreCancel
func etExportRestoreCancel(ptr *C.etExportRestore) C.etExportRestoreStatus {
	ce, ok := resolveExportRestore(ptr)
	if !ok {
		return C.ET_EXPORT_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	ce.restorer.Cancel()

	return C.ET_EXPORT_RESTORE_STATUS_OK
}

//export etExportRestoreGetLastError
func etExportRestoreGetLastError(ptr *C.etExportRestore) *C.cchar_t {
	ce, ok := resolveExportRestore(ptr)
	if !ok {
		return nil
	}

	return (*C.cchar_t)(ce.lastError.GetErr())
}

//export etExportRestoreGetBackupPath
func etExportRestoreGetBackupPath(ptr *C.etExportRestore, outPath **C.char) C.etExportRestoreStatus {
	ce, ok := resolveExportRestore(ptr)
	if !ok {
		return C.ET_EXPORT_RESTORE_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*outPath = C.CString(ce.restorer.GetBackupPath())

	return C.ET_EXPORT_RESTORE_STATUS_OK
}

type cExportRestore struct {
	csession  *csession
	restorer  *mail.RestoreTask
	lastError utils.CLastError
}

type ExportRestoreHandle struct {
	internal.Handle
}

func (h ExportRestoreHandle) resolve() (*cExportRestore, bool) {
	return internal.ResolveHandle[cExportRestore](h.Handle)
}

func exportRestorePtrToHandle(ptr *C.etExportRestore) ExportRestoreHandle {
	return ExportRestoreHandle{Handle: cgo.Handle(unsafe.Pointer(ptr))}
}

func resolveExportRestore(ptr *C.etExportRestore) (*cExportRestore, bool) {
	h := exportRestorePtrToHandle(ptr)

	return h.resolve()
}

type restoreReporter struct {
	totalMessageCount   atomic.Uint64
	currentMessageCount atomic.Uint64
	callbacks           *C.etExportRestoreCallbacks
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

	C.etExportRestoreCallbackOnProgress(m.callbacks, C.float(progress))
}
