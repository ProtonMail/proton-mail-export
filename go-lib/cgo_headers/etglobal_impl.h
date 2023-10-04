#ifndef ET_GLOBAL_IMPL_H
#define ET_GLOBAL_IMPL_H

#include "etglobal.h"

#ifdef ET_CGO

inline void etCallOnRecover(etOnRecoverFn cb) {
    cb();
}

#endif // ET_CGO

#endif // ET_GLOBAL_IMPL_H
