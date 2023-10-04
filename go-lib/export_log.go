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

/*
	#include <stdlib.h>
	typedef const char cchar_t;
*/
import "C"
import (
	"github.com/sirupsen/logrus"
)

//export etLogInfo
func etLogInfo(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Info(C.GoString(txt))
}

//export etLogDebug
func etLogDebug(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Debug(C.GoString(txt))
}

//export etLogWarn
func etLogWarn(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Warn(C.GoString(txt))
}

//export etLogError
func etLogError(tag *C.cchar_t, txt *C.cchar_t) {
	logrus.WithField("tag", C.GoString(tag)).Error(C.GoString(txt))
}

//export etLogGetPath
func etLogGetPath() *C.cchar_t {
	etGlobalState.mutex.Lock()
	defer etGlobalState.mutex.Unlock()

	return etGlobalState.clogPath
}
