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
#include "etexport_mail.h"
#include "etexport_mail_impl.h"
*/
import "C"
import (
	"github.com/ProtonMail/export-tool/internal/utils"
	"path/filepath"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/session"
)

//export etSessionNewExportMail
func etSessionNewExportMail(sessionPtr *C.etSession, cExportPath *C.cchar_t, outExportMail **C.etExportMail) C.etSessionStatus {
	csession, ok := resolveSession(sessionPtr)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	if csession.s.LoginState() != session.LoginStateLoggedIn {
		csession.setLastError(session.ErrInvalidLoginState)
		return C.ET_SESSION_STATUS_ERROR
	}

	exportPath := C.GoString(cExportPath)
	exportPath = filepath.Join(exportPath, csession.s.GetEmail())

	mailExport := mail.NewExportTask(csession.ctx, exportPath, csession.s)

	h := exportMailAllocator.Alloc(&cExportMail{
		csession: csession,
		exporter: mailExport,
	})

	*outExportMail = (*C.etExportMail)(unsafe.Pointer(uintptr(h)))

	return C.ET_SESSION_STATUS_OK
}

//export etExportMailDelete
func etExportMailDelete(ptr *C.etExportMail) C.etExportMailStatus {
	h := exportMailPtrToHandle(ptr)

	s, ok := exportMailAllocator.Resolve(h)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	s.exporter.Close()
	s.lastError.Close()

	exportMailAllocator.Free(h)

	return C.ET_EXPORT_MAIL_STATUS_OK
}

//export etExportMailStart
func etExportMailStart(ptr *C.etExportMail, callbacks *C.etExportMailCallbacks) C.etExportMailStatus {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	reporter := &mailExportReporter{
		exporter:            ce.exporter,
		totalMessageCount:   0,
		currentMessageCount: 0,
		callbacks:           callbacks,
	}

	if err := ce.exporter.Run(ce.csession.ctx, reporter); err != nil {
		ce.lastError.Set(err)
		return C.ET_EXPORT_MAIL_STATUS_ERROR
	}

	return C.ET_EXPORT_MAIL_STATUS_OK
}

//export etExportMailCancel
func etExportMailCancel(ptr *C.etExportMail) C.etExportMailStatus {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	ce.exporter.Cancel()

	return C.ET_EXPORT_MAIL_STATUS_OK
}

//export etExportMailGetLastError
func etExportMailGetLastError(ptr *C.etExportMail) *C.cchar_t {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return nil
	}

	return (*C.cchar_t)(ce.lastError.GetErr())
}

type cExportMail struct {
	csession  *csession
	exporter  *mail.ExportTask
	lastError utils.CLastError
}

var exportMailAllocator = internal.NewHandleMap[cExportMail](5)

type ExportMailHandle = internal.Handle[cExportMail]

func exportMailPtrToHandle(ptr *C.etExportMail) ExportMailHandle {
	return ExportMailHandle(uintptr(unsafe.Pointer(ptr)))
}

func resolveExportMail(ptr *C.etExportMail) (*cExportMail, bool) {
	h := exportMailPtrToHandle(ptr)

	return exportMailAllocator.Resolve(h)
}

type mailExportReporter struct {
	totalMessageCount   uint64
	currentMessageCount uint64
	callbacks           *C.etExportMailCallbacks
	exporter            *mail.ExportTask
}

func (m *mailExportReporter) SetMessageTotal(total uint64) {
	m.totalMessageCount = total
}

func (m *mailExportReporter) OnProgress(delta int) {
	m.currentMessageCount += uint64(delta)

	var progress float32
	if m.totalMessageCount != 0 {
		progress = float32(float64(m.currentMessageCount) / float64(m.totalMessageCount) * 100.0)
	} else {
		progress = float32(0.0)
	}

	C.etExportMailCallbackOnProgress(m.callbacks, C.float(progress))
}
