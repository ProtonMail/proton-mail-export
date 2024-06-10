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

#ifndef ET_EXPORT_BACKUP_H
#define ET_EXPORT_BACKUP_H

#include "etsession.h"

typedef struct etExportBackup etExportBackup;

typedef enum etExportBackupStatus {
	ET_EXPORT_BACKUP_STATUS_OK,
	ET_EXPORT_BACKUP_STATUS_ERROR,
	ET_EXPORT_BACKUP_STATUS_INVALID,
	ET_EXPORT_BACKUP_STATUS_CANCELLED,
} etExportBackupStatus;

typedef enum etExportBackupMessageType {
	ET_EXPORT_BACKUP_MESSAGE_TYPE_PROGRESS,
} etExportBackupMessageType;

typedef struct etExportBackupCallbacks {
    void* ptr;
    void (*onProgress)(void* ptr, float progress);
} etExportBackupCallbacks;

#endif // ET_EXPORT_BACKUP_H



