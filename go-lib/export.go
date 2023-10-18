// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
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
#include "etsession.h"
#include "etsession_impl.h"
*/
import "C"
import (
	"context"
	"errors"
	"sync"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/sentry"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/gluon/async"
)

type SessionHandle = internal.Handle[csession]

//export etSessionNew
func etSessionNew(apiURL *C.cchar_t, cb C.etSessionCallbacks, cErr **C.char) *C.etSession {
	goAPIURL := C.GoString(apiURL)

	cSession, err := newCSession(goAPIURL, cb)
	if err != nil {
		*cErr = C.CString(err.Error())
		return nil
	}

	h := sessionAllocator.Alloc(cSession)
	// Intentional misuse of unsafe pointer.
	p := unsafe.Pointer(uintptr(h)) //nolint:govet
	return (*C.etSession)(p)
}

//export etSessionDelete
func etSessionDelete(ptr *C.etSession) C.etSessionStatus {
	h := ptrToHandle(ptr)

	s, ok := sessionAllocator.Resolve(h)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	s.close()

	sessionAllocator.Free(h)

	return C.ET_SESSION_STATUS_OK
}

//export etSessionGetLastError
func etSessionGetLastError(ptr *C.etSession) *C.cchar_t {
	h := ptrToHandle(ptr)

	s, ok := sessionAllocator.Resolve(h)
	if !ok {
		return nil
	}

	return (*C.cchar_t)(s.lastError.GetErr())
}

//export etSessionCancel
func etSessionCancel(ptr *C.etSession) C.etSessionStatus {
	h := ptrToHandle(ptr)

	s, ok := sessionAllocator.Resolve(h)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	s.cancel()

	return C.ET_SESSION_STATUS_OK
}

//export etSessionGetLoginState
func etSessionGetLoginState(ptr *C.etSession, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionLogin
func etSessionLogin(ptr *C.etSession, email *C.cchar_t, password *C.cchar_t, passwordLen C.int, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		if err := session.Login(ctx, C.GoString(email), C.GoBytes(unsafe.Pointer(password), passwordLen)); err != nil {
			return internal.MapError(err)
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionLogout
func etSessionLogout(ptr *C.etSession) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		return internal.MapError(session.Logout(ctx))
	})
}

//export etSessionSubmitTOTP
func etSessionSubmitTOTP(ptr *C.etSession, totp *C.cchar_t, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		if err := session.SubmitTOTP(ctx, C.GoString(totp)); err != nil {
			return internal.MapError(err)
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionSubmitMailboxPassword
func etSessionSubmitMailboxPassword(ptr *C.etSession, password *C.cchar_t, passwordLen C.int, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		validator := apiclient.NewProtonMailboxPasswordValidator(session.GetUser(), session.GetUserSalts())
		if err := session.SubmitMailboxPassword(validator, C.GoBytes(unsafe.Pointer(password), passwordLen)); err != nil {
			return err
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionGetHVSolveURL
func etSessionGetHVSolveURL(ptr *C.etSession, outURL **C.char) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		hvURL, err := session.GetHVSolveURL()
		if err != nil {
			return err
		}

		*outURL = C.CString(hvURL)

		return nil
	})
}

//export etSessionMarkHVSolved
func etSessionMarkHVSolved(ptr *C.etSession, outLoginState *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		if err := session.MarkHVSolved(ctx); err != nil {
			return err
		}

		*outLoginState = mapLoginState(session.LoginState())
		return nil
	})
}

//export etFree
func etFree(ptr *C.void) {
	C.free(unsafe.Pointer(ptr))
}

func mapLoginState(s session.LoginState) C.etSessionLoginState {
	switch s {
	case session.LoginStateLoggedOut:
		return C.ET_SESSION_LOGIN_STATE_LOGGED_OUT
	case session.LoginStateAwaitingTOTP:
		return C.ET_SESSION_LOGIN_STATE_AWAITING_TOTP
	case session.LoginStateAwaitingMailboxPassword:
		return C.ET_SESSION_LOGIN_STATE_AWAITING_MAILBOX_PASSWORD
	case session.LoginStateAwaitingHV:
		return C.ET_SESSION_LOGIN_STATE_AWAITING_HV
	case session.LoginStateLoggedIn:
		return C.ET_SESSION_LOGIN_STATE_LOGGED_IN
	default:
		return C.ET_SESSION_LOGIN_STATE_LOGGED_OUT
	}
}

func withSession(ptr *C.etSession, f func(ctx context.Context, session *session.Session) error) C.etSessionStatus {
	session, ok := resolveSession(ptr)
	if !ok {
		return C.ET_SESSION_STATUS_INVALID
	}

	defer async.HandlePanic(session.s.GetPanicHandler())

	err := f(session.ctx, session.s)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return C.ET_SESSION_STATUS_CANCELLED
		}

		session.setLastError(err)
		return C.ET_SESSION_STATUS_ERROR
	}

	return C.ET_SESSION_STATUS_OK
}

type csession struct {
	s          *session.Session
	ctx        context.Context
	cancelOnce sync.Once
	ctxCancel  func()
	lastError  utils.CLastError
}

func newCSession(apiURL string, cb C.etSessionCallbacks) (*csession, error) {
	panicHandler := sentry.NewPanicHandler(GetGlobalOnRecoverCB())
	reporter := GetGlobalReporter()

	defer async.HandlePanic(panicHandler)

	sessionCb := newCSessionCallback(cb)
	builder, err := apiclient.NewProtonAPIClientBuilder(apiURL, panicHandler, sessionCb)
	if err != nil {
		return nil, err
	}

	clientBuilder := apiclient.NewAutoRetryClientBuilder(
		builder,
		&apiclient.SleepRetryStrategyBuilder{},
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &csession{
		s:         session.NewSession(clientBuilder, sessionCb, panicHandler, reporter),
		ctx:       ctx,
		ctxCancel: cancel,
	}, nil
}

func (c *csession) close() {
	defer async.HandlePanic(c.s.GetPanicHandler())

	c.s.Close(c.ctx)
	c.cancel()
	c.lastError.Close()
}

func (c *csession) cancel() {
	defer async.HandlePanic(c.s.GetPanicHandler())
	c.cancelOnce.Do(c.ctxCancel)
}

func (c *csession) setLastError(err error) {
	c.lastError.Set(err)
}

//nolint:gochecknoglobals
var sessionAllocator = internal.NewHandleMap[csession](5)

func ptrToHandle(ptr *C.etSession) SessionHandle {
	return SessionHandle(uintptr(unsafe.Pointer(ptr)))
}

func resolveSession(ptr *C.etSession) (*csession, bool) {
	h := ptrToHandle(ptr)

	return sessionAllocator.Resolve(h)
}

type csessionCallback struct {
	cb C.etSessionCallbacks
}

func newCSessionCallback(cb C.etSessionCallbacks) session.Callbacks {
	return &csessionCallback{cb: cb}
}

func (c *csessionCallback) OnNetworkRestored() {
	C.etSessionCallbackOnNetworkRestored(&c.cb) //nolint:gocritic
}

func (c *csessionCallback) OnNetworkLost() {
	C.etSessionCallbackOnNetworkLost(&c.cb) //nolint:gocritic
}

func main() {}
