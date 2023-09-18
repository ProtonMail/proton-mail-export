#ifndef ET_EXPORT_MAIL_IMPL_H
#define ET_EXPORT_MAIL_IMPL_H

#include "etexport_mail.h"

#ifdef ET_CGO

inline void etExportMailCallbackOnProgress(etExportMailCallbacks* cb, float progress) {
    cb->onProgress(cb->ptr, progress);
}

#endif // ET_CGO

#endif // ET_EXPORT_MAIL_IMPL_H
