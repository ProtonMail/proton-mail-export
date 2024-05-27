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

//go:build darwin

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
#import <stdio.h>
char const* getMacOSDownloadsDir() {
    @autoreleasepool{
        NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDownloadsDirectory, NSUserDomainMask, YES);
		// The buffer returned by UTF8String's lifetime does not exceed the lifetime of paths, so we must copy it.
		return strdup([[paths objectAtIndex:0] UTF8String]);// memory allocated by strudp is to be related be freed on the Go side.
	}
}
*/
import "C"
import (
	"path/filepath"
	"unsafe"
)

func getDefaultOperationFolder() (string, error) {
	var cStr *C.char = C.getMacOSDownloadsDir()
	var result string = C.GoString(cStr)
	C.free(unsafe.Pointer(cStr))
	return filepath.Join(result, "proton-mail-export-cli"), nil
}
