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

#ifndef ET_RESTORE_H
#define ET_RESTORE_H

#include "etsession.h"
#include <stdint.h>

typedef struct etRestore etRestore;

typedef enum etRestoreStatus {
	ET_RESTORE_STATUS_OK,
	ET_RESTORE_STATUS_ERROR,
	ET_RESTORE_STATUS_INVALID,
	ET_RESTORE_STATUS_CANCELLED,
} etRestoreStatus;

typedef enum etRestoreMessageType {
	ET_RESTORE_MESSAGE_TYPE_PROGRESS,
} etRestoreMessageType;

typedef struct etRestoreCallbacks {
    void* ptr;
    void (*onProgress)(void* ptr, float progress);
} etRestoreCallbacks;

#endif // ET_RESTORE_H



