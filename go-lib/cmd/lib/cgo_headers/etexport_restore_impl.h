// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
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
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

#ifndef ET_EXPORT_RESTORE_IMPL_H
#define ET_EXPORT_RESTORE_IMPL_H

#include "etexport_restore.h"

#ifdef ET_CGO

inline void etExportRestoreCallbackOnProgress(etExportRestoreCallbacks* cb, float progress) {
    cb->onProgress(cb->ptr, progress);
}

#endif // ET_CGO

#endif // ET_EXPORT_RESTORE_IMPL_H
