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

#include "etexport_restore.hpp"

#include <proton-mail-export.h>
#include "etsession.hpp"

namespace etcpp {
inline void mapETExportRestoreStatusToException(etExportRestore* ptr, etExportRestoreStatus const status) {
    switch (status) {
    case ET_EXPORT_RESTORE_STATUS_INVALID:
        throw SessionException("Invalid instance");
    case ET_EXPORT_RESTORE_STATUS_ERROR:
    {
        const char* lastErr = etExportRestoreGetLastError(ptr);
        if (lastErr == nullptr) {
            lastErr = "unknown";
        }
        throw ExportRestoreException(lastErr);
    }
    case ET_EXPORT_RESTORE_STATUS_CANCELLED:
        throw CancelledException();
    case ET_EXPORT_RESTORE_STATUS_OK:
        break;
    }
}

etExportRestoreCallbacks makeETRestoreCallback(ExportRestoreCallback& cb) {
    auto r = etExportRestoreCallbacks{};
    r.ptr = &cb;
    r.onProgress = [](void* p, float progress) { reinterpret_cast<ExportRestoreCallback*>(p)->onProgress(progress); };

    return r;
}

ExportRestore::ExportRestore(const etcpp::Session& session, etExportRestore* ptr) :
    mSession(session), mPtr(ptr) {
}

ExportRestore::~ExportRestore() {
    wrapCCall([](etExportRestore* ptr) { return etExportRestoreDelete(ptr); });
}

void ExportRestore::start(ExportRestoreCallback& cb) {
    wrapCCall([&](etExportRestore* ptr) {
        auto etCb = makeETRestoreCallback(cb);
        return etExportRestoreStart(ptr, &etCb);
    });
}

void ExportRestore::cancel() {
    wrapCCall([&](etExportRestore* ptr) { return etExportRestoreCancel(ptr); });
}

std::filesystem::path ExportRestore::getBackupPath() const {
    char* outPath = nullptr;
    wrapCCall([&](etExportRestore* ptr) { return etExportRestoreGetBackupPath(ptr, &outPath); });

    auto result = std::filesystem::u8path(outPath);
    etFree(outPath);

    return result;
}

template<class F>
void ExportRestore::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etExportRestoreStatus, F, etExportRestore*>, "invalid function/lambda signature");
    mapETExportRestoreStatusToException(mPtr, func(mPtr));
}

template<class F>
void ExportRestore::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etExportRestoreStatus, F, etExportRestore*>, "invalid function/lambda signature");
    mapETExportRestoreStatusToException(mPtr, func(mPtr));
}
} // namespace etcpp