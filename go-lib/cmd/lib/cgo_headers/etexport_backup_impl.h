#ifndef ET_EXPORT_BACKUP_IMPL_H
#define ET_EXPORT_BACKUP_IMPL_H

#include "etexport_backup.h"

#ifdef ET_CGO

inline void etExportBackupCallbackOnProgress(etExportBackupCallbacks* cb, float progress) {
    cb->onProgress(cb->ptr, progress);
}

#endif // ET_CGO

#endif // ET_EXPORT_MAIL_IMPL_H
