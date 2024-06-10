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

#include "etexport_backup.hpp"

#include <proton-mail-export.h>
#include "etsession.hpp"

namespace etcpp {

inline void mapETExportBackupStatusToException(etExportBackup* ptr, etExportBackupStatus status) {
    switch (status) {
    case ET_EXPORT_BACKUP_STATUS_INVALID:
        throw SessionException("Invalid instance");
    case ET_EXPORT_BACKUP_STATUS_ERROR:
    {
        const char* lastErr = etExportBackupGetLastError(ptr);
        if (lastErr == nullptr) {
            lastErr = "unknown";
        }
        throw ExportBackupException(lastErr);
    }
    case ET_EXPORT_BACKUP_STATUS_CANCELLED:
        throw CancelledException();
    case ET_EXPORT_BACKUP_STATUS_OK:
        break;
    }
}

etExportBackupCallbacks makeETCallback(ExportBackupCallback& cb) {
    auto r = etExportBackupCallbacks{};
    r.ptr = &cb;
    r.onProgress = [](void* p, float progress) { reinterpret_cast<ExportBackupCallback*>(p)->onProgress(progress); };

    return r;
}

ExportBackup::ExportBackup(const etcpp::Session& session, etExportBackup* ptr) : mSession(session), mPtr(ptr) {}

ExportBackup::~ExportBackup() {
    wrapCCall([](etExportBackup* ptr) { return etExportBackupDelete(ptr); });
}

void ExportBackup::start(ExportBackupCallback& cb) {
    wrapCCall([&](etExportBackup* ptr) {
        auto etCb = makeETCallback(cb);
        return etExportBackupStart(ptr, &etCb);
    });
}

void ExportBackup::cancel() {
    wrapCCall([&](etExportBackup* ptr) { return etExportBackupCancel(ptr); });
}

std::filesystem::path ExportBackup::getExportPath() const {
    char* outPath = nullptr;
    wrapCCall([&](etExportBackup* ptr) { return etExportBackupGetExportPath(ptr, &outPath); });

    auto result = std::filesystem::u8path(outPath);
    etFree(outPath);

    return result;
}

std::uint64_t ExportBackup::getExpectedDiskUsage() const {
    std::uint64_t usage = 0;
    wrapCCall([&](etExportBackup* ptr) { return etExportBackupGetRequiredDiskSpaceEstimate(ptr, &usage); });
    return usage;
}

template<class F>
void ExportBackup::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etExportBackupStatus, F, etExportBackup*>, "invalid function/lambda signature");
    mapETExportBackupStatusToException(mPtr, func(mPtr));
}

template<class F>
void ExportBackup::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etExportBackupStatus, F, etExportBackup*>, "invalid function/lambda signature");
    mapETExportBackupStatusToException(mPtr, func(mPtr));
}

} // namespace etcpp