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

#include "etbackup.hpp"

#include <proton-mail-export.h>
#include "etsession.hpp"

namespace etcpp {

inline void mapETBackupStatusToException(etBackup* ptr, etBackupStatus status) {
    switch (status) {
    case ET_BACKUP_STATUS_INVALID:
        throw SessionException("Invalid instance");
    case ET_BACKUP_STATUS_ERROR:
    {
        const char* lastErr = etBackupGetLastError(ptr);
        if (lastErr == nullptr) {
            lastErr = "unknown";
        }
        throw BackupException(lastErr);
    }
    case ET_BACKUP_STATUS_CANCELLED:
        throw CancelledException();
    case ET_BACKUP_STATUS_OK:
        break;
    }
}

etBackupCallbacks makeETCallback(BackupCallback& cb) {
    auto r = etBackupCallbacks{};
    r.ptr = &cb;
    r.onProgress = [](void* p, float progress) { reinterpret_cast<BackupCallback*>(p)->onProgress(progress); };

    return r;
}

Backup::Backup(const etcpp::Session& session, etBackup* ptr) : mSession(session), mPtr(ptr) {}

Backup::~Backup() {
    wrapCCall([](etBackup* ptr) { return etBackupDelete(ptr); });
}

void Backup::start(BackupCallback& cb) {
    wrapCCall([&](etBackup* ptr) {
        auto etCb = makeETCallback(cb);
        return etBackupStart(ptr, &etCb);
    });
}

void Backup::cancel() {
    wrapCCall([&](etBackup* ptr) { return etBackupCancel(ptr); });
}

std::filesystem::path Backup::getExportPath() const {
    char* outPath = nullptr;
    wrapCCall([&](etBackup* ptr) { return etBackupGetExportPath(ptr, &outPath); });

    auto result = std::filesystem::u8path(outPath);
    etFree(outPath);

    return result;
}

std::uint64_t Backup::getExpectedDiskUsage() const {
    std::uint64_t usage = 0;
    wrapCCall([&](etBackup* ptr) { return etBackupGetRequiredDiskSpaceEstimate(ptr, &usage); });
    return usage;
}

template<class F>
void Backup::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etBackupStatus, F, etBackup*>, "invalid function/lambda signature");
    mapETBackupStatusToException(mPtr, func(mPtr));
}

template<class F>
void Backup::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etBackupStatus, F, etBackup*>, "invalid function/lambda signature");
    mapETBackupStatusToException(mPtr, func(mPtr));
}

} // namespace etcpp
