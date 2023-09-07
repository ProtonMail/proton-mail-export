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

/*
#include <stdint.h>
#include <stdlib.h>

typedef const char cchar_t;

typedef struct etSession etSession;

typedef enum etSessionStatus {
	ET_SESSION_STATUS_OK,
	ET_SESSION_STATUS_ERROR,
	ET_SESSION_STATUS_INVALID,
} etSessionStatus;

typedef enum etSessionLoginState {
	ET_SESSION_LOGIN_STATE_LOGGED_OUT,
	ET_SESSION_LOGIN_STATE_AWAITING_TOTP,
	ET_SESSION_LOGIN_STATE_AWAITING_HV,
	ET_SESSION_LOGIN_STATE_AWAITING_MAILBOX_PASSWORD,
	ET_SESSION_LOGIN_STATE_LOGGED_IN,
} etSessionLoginState;

*/
import "C"
import (
	"context"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/gluon/async"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
)

type SessionHandle = internal.Handle[csession]

//export etSessionNew
func etSessionNew(apiURL *C.cchar_t) *C.etSession {
	goAPIURL := C.GoString(apiURL)
	h := sessionAllocator.Alloc(newCSession(goAPIURL))
	p := unsafe.Pointer(uintptr(h))
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

	return s.lastError
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
			return err
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionLogout
func etSessionLogout(ptr *C.etSession) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		return session.Logout(ctx)
	})
}

//export etSessionSubmitTOTP
func etSessionSubmitTOTP(ptr *C.etSession, totp *C.cchar_t, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		if err := session.SubmitTOTP(ctx, C.GoString(totp)); err != nil {
			return err
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
}

//export etSessionSubmitMailboxPassword
func etSessionSubmitMailboxPassword(ptr *C.etSession, password *C.cchar_t, passwordLen C.int, outStatus *C.etSessionLoginState) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *session.Session) error {
		if err := session.SubmitMailboxPassword(C.GoBytes(unsafe.Pointer(password), passwordLen)); err != nil {
			return err
		}

		*outStatus = mapLoginState(session.LoginState())
		return nil
	})
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

	err := f(session.ctx, session.s)
	if err != nil {
		session.setLastError(err)
		return C.ET_SESSION_STATUS_ERROR
	}

	return C.ET_SESSION_STATUS_OK
}

type csession struct {
	s         *session.Session
	ctx       context.Context
	ctxCancel func()
	lastError *C.char
}

func newCSession(apiURL string) *csession {
	clientBuilder := apiclient.NewProtonAPIClientBuilder(apiURL, &async.NoopPanicHandler{})
	return &csession{
		s:         session.NewSession(clientBuilder),
		lastError: nil,
	}
}

func (c *csession) close() {
	c.s.Close(c.ctx)
	c.ctxCancel()
	if c.lastError != nil {
		C.free(unsafe.Pointer(c.lastError))
		c.lastError = nil
	}
}

func (c *csession) setLastError(err error) {
	if c.lastError != nil {
		C.free(unsafe.Pointer(c.lastError))
	}

	c.lastError = C.CString(err.Error())
}

var sessionAllocator = internal.NewHandleMap[csession](5)

func ptrToHandle(ptr *C.etSession) SessionHandle {
	return SessionHandle(uintptr(unsafe.Pointer(ptr)))
}

func resolveSession(ptr *C.etSession) (*csession, bool) {
	h := ptrToHandle(ptr)

	return sessionAllocator.Resolve(h)
}

func main() {}
