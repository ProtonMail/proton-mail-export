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
#include "etexport_backup.h"
#include "etexport_backup_impl.h"
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

//export etSessionNewExportBackup
func etSessionNewExportBackup(sessionPtr *C.etSession, cExportPath *C.cchar_t, outExportBackup **C.etExportBackup) C.etSessionStatus {
	cSession, ok := resolveSession(sessionPtr)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	defer async.HandlePanic(cSession.s.GetPanicHandler())

	if cSession.s.LoginState() != session.LoginStateLoggedIn {
		cSession.setLastError(session.ErrInvalidLoginState)
		return C.ET_SESSION_STATUS_ERROR
	}

	exportPath := C.GoString(cExportPath)
	exportPath = filepath.Join(exportPath, cSession.s.GetUser().Email)

	mailExport := mail.NewExportTask(cSession.ctx, exportPath, cSession.s)

	h := internal.NewHandle(&cExportBackup{
		csession: cSession,
		exporter: mailExport,
	})

	// Intentional misuse of unsafe pointer.
	//goland:noinspection GoVetUnsafePointer
	*outExportBackup = (*C.etExportBackup)(unsafe.Pointer(h)) //nolint:govet

	return C.ET_SESSION_STATUS_OK
}

//export etExportBackupDelete
func etExportBackupDelete(ptr *C.etExportBackup) C.etExportBackupStatus {
	h := exportMailPtrToHandle(ptr)

	s, ok := h.resolve()
	if !ok {
		return C.ET_EXPORT_BACKUP_STATUS_INVALID
	}

	defer async.HandlePanic(s.csession.s.GetPanicHandler())

	s.exporter.Close()
	s.lastError.Close()

	h.Delete()

	return C.ET_EXPORT_BACKUP_STATUS_OK
}

//export etExportBackupStart
func etExportBackupStart(ptr *C.etExportBackup, callbacks *C.etExportBackupCallbacks) C.etExportBackupStatus {
	ce, ok := resolveExportBackup(ptr)
	if !ok {
		return C.ET_EXPORT_BACKUP_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	reporter := &backupReporter{
		exporter:  ce.exporter,
		callbacks: callbacks,
	}

	if err := ce.exporter.Run(ce.csession.ctx, reporter); err != nil {
		if errors.Is(err, context.Canceled) {
			return C.ET_EXPORT_BACKUP_STATUS_CANCELLED
		}

		ce.lastError.Set(internal.MapError(err))
		return C.ET_EXPORT_BACKUP_STATUS_ERROR
	}

	return C.ET_EXPORT_BACKUP_STATUS_OK
}

//export etExportBackupCancel
func etExportBackupCancel(ptr *C.etExportBackup) C.etExportBackupStatus {
	ce, ok := resolveExportBackup(ptr)
	if !ok {
		return C.ET_EXPORT_BACKUP_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	ce.exporter.Cancel()

	return C.ET_EXPORT_BACKUP_STATUS_OK
}

//export etExportBackupGetLastError
func etExportBackupGetLastError(ptr *C.etExportBackup) *C.cchar_t {
	ce, ok := resolveExportBackup(ptr)
	if !ok {
		return nil
	}

	return (*C.cchar_t)(ce.lastError.GetErr())
}

//export etExportBackupGetRequiredDiskSpaceEstimate
func etExportBackupGetRequiredDiskSpaceEstimate(ptr *C.etExportBackup, outSpace *C.uint64_t) C.etExportBackupStatus {
	ce, ok := resolveExportBackup(ptr)
	if !ok {
		return C.ET_EXPORT_BACKUP_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	space, err := ce.exporter.GetRequiredDiskSpaceEstimate(ce.csession.ctx)
	if err != nil {
		ce.lastError.Set(internal.MapError(err))
		return C.ET_EXPORT_BACKUP_STATUS_ERROR
	}

	*outSpace = C.uint64_t(space)

	return C.ET_EXPORT_BACKUP_STATUS_OK
}

//export etExportBackupGetExportPath
func etExportBackupGetExportPath(ptr *C.etExportBackup, outPath **C.char) C.etExportBackupStatus {
	ce, ok := resolveExportBackup(ptr)
	if !ok {
		return C.ET_EXPORT_BACKUP_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*outPath = C.CString(ce.exporter.GetExportPath())

	return C.ET_EXPORT_BACKUP_STATUS_OK
}

type cExportBackup struct {
	csession  *csession
	exporter  *mail.ExportTask
	lastError utils.CLastError
}

type ExportBackupHandle struct {
	internal.Handle
}

func (h ExportBackupHandle) resolve() (*cExportBackup, bool) {
	return internal.ResolveHandle[cExportBackup](h.Handle)
}

func exportMailPtrToHandle(ptr *C.etExportBackup) ExportBackupHandle {
	return ExportBackupHandle{Handle: cgo.Handle(unsafe.Pointer(ptr))}
}

func resolveExportBackup(ptr *C.etExportBackup) (*cExportBackup, bool) {
	h := exportMailPtrToHandle(ptr)

	return h.resolve()
}

type backupReporter struct {
	totalMessageCount   atomic.Uint64
	currentMessageCount atomic.Uint64
	callbacks           *C.etExportBackupCallbacks
	exporter            *mail.ExportTask
}

func (m *backupReporter) SetMessageTotal(total uint64) {
	m.totalMessageCount.Store(total)
}

func (m *backupReporter) SetMessageProcessed(total uint64) {
	m.currentMessageCount.Store(total)
}

func (m *backupReporter) OnProgress(delta int) {
	newMessageCount := m.currentMessageCount.Add(uint64(delta))

	var progress float32
	totalMessageCount := m.totalMessageCount.Load()
	if totalMessageCount != 0 {
		progress = float32(float64(newMessageCount) / float64(totalMessageCount) * 100.0)
	} else {
		progress = float32(0.0)
	}

	C.etExportBackupCallbackOnProgress(m.callbacks, C.float(progress))
}
