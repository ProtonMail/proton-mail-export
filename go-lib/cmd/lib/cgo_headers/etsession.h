#ifndef ET_SESSION_H
#define ET_SESSION_H

#include <stdint.h>
#include <stdlib.h>

typedef const char cchar_t;

typedef struct etSession etSession;

typedef enum etSessionStatus {
	ET_SESSION_STATUS_OK,
	ET_SESSION_STATUS_ERROR,
	ET_SESSION_STATUS_INVALID,
	ET_SESSION_STATUS_CANCELLED,
} etSessionStatus;

typedef enum etSessionLoginState {
	ET_SESSION_LOGIN_STATE_LOGGED_OUT,
	ET_SESSION_LOGIN_STATE_AWAITING_TOTP,
	ET_SESSION_LOGIN_STATE_AWAITING_HV,
	ET_SESSION_LOGIN_STATE_AWAITING_MAILBOX_PASSWORD,
	ET_SESSION_LOGIN_STATE_LOGGED_IN,
} etSessionLoginState;

typedef struct etSessionCallbacks {
    void *ptr;
    void (*onNetworkLost)(void*);
    void (*onNetworkRestored)(void*);
} etSessionCallbacks;

#endif // ET_SESSION_H