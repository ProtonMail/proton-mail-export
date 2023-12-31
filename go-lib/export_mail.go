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

//export etSessionNewExportMail
func etSessionNewExportMail(sessionPtr *C.etSession, cExportPath *C.cchar_t, outExportMail **C.etExportMail) C.etSessionStatus {
	csession, ok := resolveSession(sessionPtr)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	defer async.HandlePanic(csession.s.GetPanicHandler())

	if csession.s.LoginState() != session.LoginStateLoggedIn {
		csession.setLastError(session.ErrInvalidLoginState)
		return C.ET_SESSION_STATUS_ERROR
	}

	exportPath := C.GoString(cExportPath)
	exportPath = filepath.Join(exportPath, csession.s.GetUser().Email)

	mailExport := mail.NewExportTask(csession.ctx, exportPath, csession.s)

	h := internal.NewHandle(&cExportMail{
		csession: csession,
		exporter: mailExport,
	})

	// Intentional misuse of unsafe pointer.
	*outExportMail = (*C.etExportMail)(unsafe.Pointer(h)) //nolint:govet

	return C.ET_SESSION_STATUS_OK
}

//export etExportMailDelete
func etExportMailDelete(ptr *C.etExportMail) C.etExportMailStatus {
	h := exportMailPtrToHandle(ptr)

	s, ok := h.resolve()
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	defer async.HandlePanic(s.csession.s.GetPanicHandler())

	s.exporter.Close()
	s.lastError.Close()

	h.Delete()

	return C.ET_EXPORT_MAIL_STATUS_OK
}

//export etExportMailStart
func etExportMailStart(ptr *C.etExportMail, callbacks *C.etExportMailCallbacks) C.etExportMailStatus {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	reporter := &mailExportReporter{
		exporter:  ce.exporter,
		callbacks: callbacks,
	}

	if err := ce.exporter.Run(ce.csession.ctx, reporter); err != nil {
		if errors.Is(err, context.Canceled) {
			return C.ET_EXPORT_MAIL_STATUS_CANCELLED
		}

		ce.lastError.Set(internal.MapError(err))
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

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

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

//export etExportMailGetRequiredDiskSpaceEstimate
func etExportMailGetRequiredDiskSpaceEstimate(ptr *C.etExportMail, outSpace *C.uint64_t) C.etExportMailStatus {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	space, err := ce.exporter.GetRequiredDiskSpaceEstimate(ce.csession.ctx)
	if err != nil {
		ce.lastError.Set(internal.MapError(err))
		return C.ET_EXPORT_MAIL_STATUS_ERROR
	}

	*outSpace = C.uint64_t(space)

	return C.ET_EXPORT_MAIL_STATUS_OK
}

//export etExportMailGetExportPath
func etExportMailGetExportPath(ptr *C.etExportMail, outPath **C.char) C.etExportMailStatus {
	ce, ok := resolveExportMail(ptr)
	if !ok {
		return C.ET_EXPORT_MAIL_STATUS_INVALID
	}

	defer async.HandlePanic(ce.csession.s.GetPanicHandler())

	*outPath = C.CString(ce.exporter.GetExportPath())

	return C.ET_EXPORT_MAIL_STATUS_OK
}

type cExportMail struct {
	csession  *csession
	exporter  *mail.ExportTask
	lastError utils.CLastError
}

type ExportMailHandle struct {
	internal.Handle
}

func (h ExportMailHandle) resolve() (*cExportMail, bool) {
	return internal.ResolveHandle[cExportMail](h.Handle)
}

func exportMailPtrToHandle(ptr *C.etExportMail) ExportMailHandle {
	return ExportMailHandle{Handle: cgo.Handle(unsafe.Pointer(ptr))}
}

func resolveExportMail(ptr *C.etExportMail) (*cExportMail, bool) {
	h := exportMailPtrToHandle(ptr)

	return h.resolve()
}

type mailExportReporter struct {
	totalMessageCount   atomic.Uint64
	currentMessageCount atomic.Uint64
	callbacks           *C.etExportMailCallbacks
	exporter            *mail.ExportTask
}

func (m *mailExportReporter) SetMessageTotal(total uint64) {
	m.totalMessageCount.Store(total)
}

func (m *mailExportReporter) SetMessageDownloaded(total uint64) {
	m.currentMessageCount.Store(total)
}

func (m *mailExportReporter) OnProgress(delta int) {
	newMessageCount := m.currentMessageCount.Add(uint64(delta))

	var progress float32
	totalMessageCount := m.totalMessageCount.Load()
	if totalMessageCount != 0 {
		progress = float32(float64(newMessageCount) / float64(totalMessageCount) * 100.0)
	} else {
		progress = float32(0.0)
	}

	C.etExportMailCallbackOnProgress(m.callbacks, C.float(progress))
}
