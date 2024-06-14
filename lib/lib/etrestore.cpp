// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
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
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

#include "etrestore.hpp"

#include <proton-mail-export.h>
#include "etsession.hpp"

namespace etcpp {
inline void mapETRestoreStatusToException(etRestore* ptr, etRestoreStatus const status) {
    switch (status) {
    case ET_RESTORE_STATUS_INVALID:
        throw SessionException("Invalid instance");
    case ET_RESTORE_STATUS_ERROR:
    {
        const char* lastErr = etRestoreGetLastError(ptr);
        if (lastErr == nullptr) {
            lastErr = "unknown";
        }
        throw RestoreException(lastErr);
    }
    case ET_RESTORE_STATUS_CANCELLED:
        throw CancelledException();
    case ET_RESTORE_STATUS_OK:
        break;
    }
}

etRestoreCallbacks makeETRestoreCallback(RestoreCallback& cb) {
    auto r = etRestoreCallbacks{};
    r.ptr = &cb;
    r.onProgress = [](void* p, float progress) { reinterpret_cast<RestoreCallback*>(p)->onProgress(progress); };

    return r;
}

Restore::Restore(const etcpp::Session& session, etRestore* ptr) : mSession(session), mPtr(ptr) {}

Restore::~Restore() {
    wrapCCall([](etRestore* ptr) { return etRestoreDelete(ptr); });
}

void Restore::start(RestoreCallback& cb) {
    wrapCCall([&](etRestore* ptr) {
        auto etCb = makeETRestoreCallback(cb);
        return etRestoreStart(ptr, &etCb);
    });
}

void Restore::cancel() {
    wrapCCall([&](etRestore* ptr) { return etRestoreCancel(ptr); });
}

std::filesystem::path Restore::getBackupPath() const {
    char* outPath = nullptr;
    wrapCCall([&](etRestore* ptr) { return etRestoreGetBackupPath(ptr, &outPath); });

    auto result = std::filesystem::u8path(outPath);
    etFree(outPath);

    return result;
}

template<class F>
void Restore::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etRestoreStatus, F, etRestore*>, "invalid function/lambda signature");
    mapETRestoreStatusToException(mPtr, func(mPtr));
}

template<class F>
void Restore::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etRestoreStatus, F, etRestore*>, "invalid function/lambda signature");
    mapETRestoreStatusToException(mPtr, func(mPtr));
}
} // namespace etcpp