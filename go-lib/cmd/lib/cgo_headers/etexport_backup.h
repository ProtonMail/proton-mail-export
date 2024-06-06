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



