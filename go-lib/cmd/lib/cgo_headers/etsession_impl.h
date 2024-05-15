#ifndef ET_SESSION_IMPL_H
#define ET_SESSION_IMPL_H

#include "etsession.h"

#ifdef ET_CGO

inline void etSessionCallbackOnNetworkLost(etSessionCallbacks* cb) {
    if (cb->onNetworkLost != NULL){
        cb->onNetworkLost(cb->ptr);
    }
}

inline void etSessionCallbackOnNetworkRestored(etSessionCallbacks* cb) {
    if (cb->onNetworkRestored != NULL){
        cb->onNetworkRestored(cb->ptr);
    }
}

#endif // ET_CGO

#endif //ET_SESSION_IMPL_H
