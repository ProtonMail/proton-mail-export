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

*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/ProtonMail/export-tool/internal"
)

type SessionHandle = internal.Handle[csession]

//export etSessionNew
func etSessionNew() *C.etSession {
	h := sessionAllocator.Alloc(newCSession())
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

//export etSessionHello
func etSessionHello(ptr *C.etSession, out **C.char) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *internal.Session) error {
		str := session.Hello()
		*out = C.CString(str)
		return nil
	})
}

//export etSessionHelloError
func etSessionHelloError(ptr *C.etSession) C.etSessionStatus {
	return withSession(ptr, func(ctx context.Context, session *internal.Session) error {
		return session.HelloError()
	})
}

func withSession(ptr *C.etSession, f func(ctx context.Context, session *internal.Session) error) C.etSessionStatus {
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
	s         *internal.Session
	ctx       context.Context
	ctxCancel func()
	lastError *C.char
}

func newCSession() *csession {
	ctx, cancel := context.WithCancel(context.Background())

	return &csession{
		s:         &internal.Session{},
		ctx:       ctx,
		ctxCancel: cancel,
		lastError: nil,
	}
}

func (c *csession) close() {
	c.s.Close()
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
