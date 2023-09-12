#ifndef ET_EXPORT_MAIL_H
#define ET_EXPORT_MAIL_H

#include "etsession.h"

typedef struct etExportMail etExportMail;

typedef enum etExportMailStatus {
	ET_EXPORT_MAIL_STATUS_OK,
	ET_EXPORT_MAIL_STATUS_ERROR,
	ET_EXPORT_MAIL_STATUS_INVALID,
} etExportMailStatus;

typedef enum etExportMailMessageType {
	ET_EXPORT_MAIL_MESSAGE_TYPE_PROGRESS,
} etExportMailMessageType;

typedef enum etExportMailCallbackReply {
    ET_EXPORT_MAIL_CALLBACK_REPLY_CONTINUE,
    ET_EXPORT_MAIL_CALLBACK_REPLY_CANCEL,
} etExportMailCallbackReply;

typedef struct etExportMailCallbacks {
    void* ptr;
    etExportMailCallbackReply (*onProgress)(void* ptr, float progress);
} etExportMailCallbacks;

#endif // ET_EXPORT_MAIL_H



