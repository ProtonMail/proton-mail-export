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

#pragma once

#include <exception>
#include <filesystem>
#include <string>

#include "etexception.hpp"

extern "C" {
struct etExportRestore;
}

namespace etcpp {

class Session;

class ExportRestoreException final : public Exception {
public:
    explicit ExportRestoreException(std::string_view what) : Exception(what) {}
};

class ExportRestoreCallback {
public:
    ExportRestoreCallback() = default;
    virtual ~ExportRestoreCallback() = default;

    virtual void onProgress(float progress) = 0;
};

class ExportRestore final {
    friend class Session;

private:
    const Session& mSession;
    etExportRestore* mPtr;

protected:
    ExportRestore(const Session& session, etExportRestore* ptr);

public:
    ~ExportRestore();
    ExportRestore(const ExportRestore&) = delete;
    ExportRestore(ExportRestore&&) noexcept = delete;
    ExportRestore& operator=(const ExportRestore&) = delete;
    ExportRestore& operator=(ExportRestore&& rhs) noexcept = delete;

    void start(ExportRestoreCallback& cb);

    void cancel();

    std::filesystem::path getBackupPath() const;

    std::uint64_t getExpectedDiskUsage() const;

private:
    template<class F>
    void wrapCCall(F func);

    template<class F>
    void wrapCCall(F func) const;
};
} // namespace etcpp
