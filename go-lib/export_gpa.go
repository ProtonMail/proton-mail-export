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

//go:build gpa_server

package main

/*
#include <stdint.h>
#include <stdlib.h>

typedef const char cchar_t;

typedef struct gpaServer gpaServer;

typedef enum gpaServerStatus {
	GPA_SERVER_STATUS_OK,
	GPA_SERVER_STATUS_FAILED,
	GPA_SERVER_STATUS_INVALID
} gpaServerStatus;

*/
import "C"
import (
	"context"
	"github.com/ProtonMail/export-tool/internal"
	"github.com/sirupsen/logrus"
	"unsafe"
)

//nolint:gochecknoglobals
var gpaServerAllocator = internal.HandleMap[internal.GPAServer]{}

type GPAHandle = internal.Handle[internal.GPAServer]

//export gpaServerNew
func gpaServerNew() *C.gpaServer {
	h := gpaServerAllocator.Alloc(internal.NewGPAServer(context.Background()))

	p := unsafe.Pointer(uintptr(h))
	return (*C.gpaServer)(p)
}

//export gpaServerDelete
func gpaServerDelete(ptr *C.gpaServer) C.gpaServerStatus {
	h := gpaPtrToHandle(ptr)

	s, ok := gpaServerAllocator.Resolve(h)
	if !ok {
		return C.GPA_SERVER_STATUS_INVALID
	}

	s.Close()

	gpaServerAllocator.Free(h)

	return C.GPA_SERVER_STATUS_OK
}

//export gpaServerCreateUser
func gpaServerCreateUser(ptr *C.gpaServer, email *C.cchar_t, password *C.cchar_t, outID **C.char, outAddrID **C.char) C.gpaServerStatus {
	s, ok := resolveGPAServer(ptr)
	if !ok {
		return C.GPA_SERVER_STATUS_INVALID
	}

	goEmail := C.GoString(email)
	goPassword := C.GoString(password)

	userID, addrID, err := s.CreateUser(goEmail, goPassword)
	if err != nil {
		return C.GPA_SERVER_STATUS_FAILED
	}

	*outID = C.CString(userID)
	*outAddrID = C.CString(addrID)

	return C.GPA_SERVER_STATUS_OK
}

//export gpaServerGetURL
func gpaServerGetURL(ptr *C.gpaServer, outURL **C.char) C.gpaServerStatus {
	s, ok := resolveGPAServer(ptr)
	if !ok {
		return C.GPA_SERVER_STATUS_INVALID
	}

	*outURL = C.CString(s.GetURL())

	return C.GPA_SERVER_STATUS_OK
}

//export gpaServerCreateMessages
func gpaServerCreateMessages(ptr *C.gpaServer,
	userID *C.cchar_t,
	addrID *C.cchar_t,
	email *C.cchar_t,
	password *C.cchar_t,
	count C.size_t,
	outIDS **C.char,
) C.gpaServerStatus {
	s, ok := resolveGPAServer(ptr)
	if !ok {
		return C.GPA_SERVER_STATUS_INVALID
	}

	goUserID := C.GoString(userID)
	goAddrID := C.GoString(addrID)
	goEmail := C.GoString(email)
	goPassword := C.GoString(password)

	messageIDs, err := s.CreateTestMessages(goUserID, goAddrID, goEmail, goPassword, int(count))
	if err != nil {
		logrus.WithError(err).Error("Failed to create test messages")
		return C.GPA_SERVER_STATUS_FAILED
	}

	cMessageIDs := unsafe.Slice(outIDS, int(count))

	for i, v := range messageIDs {
		cMessageIDs[i] = C.CString(v)
	}

	return C.GPA_SERVER_STATUS_OK
}

//export gpaFree
func gpaFree(ptr *C.void) {
	C.free(unsafe.Pointer(ptr))
}

func gpaPtrToHandle(ptr *C.gpaServer) GPAHandle {
	return GPAHandle(uintptr(unsafe.Pointer(ptr)))
}

func resolveGPAServer(ptr *C.gpaServer) (*internal.GPAServer, bool) {
	h := gpaPtrToHandle(ptr)

	return gpaServerAllocator.Resolve(h)
}

func main() {}
