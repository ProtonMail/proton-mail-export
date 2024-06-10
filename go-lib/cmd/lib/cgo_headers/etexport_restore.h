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

#ifndef ET_EXPORT_RESTORE_H
#define ET_EXPORT_RESTORE_H

#include "etsession.h"

typedef struct etExportRestore etExportRestore;

typedef enum etExportRestoreStatus {
	ET_EXPORT_RESTORE_STATUS_OK,
	ET_EXPORT_RESTORE_STATUS_ERROR,
	ET_EXPORT_RESTORE_STATUS_INVALID,
	ET_EXPORT_RESTORE_STATUS_CANCELLED,
} etExportRestoreStatus;

typedef enum etExportRestoreMessageType {
	ET_EXPORT_RESTORE_MESSAGE_TYPE_PROGRESS,
} etExportRestoreMessageType;

typedef struct etExportRestoreCallbacks {
    void* ptr;
    void (*onProgress)(void* ptr, float progress);
} etExportRestoreCallbacks;

#endif // ET_EXPORT_RESTORE_H



