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
struct etRestore;
}

namespace etcpp {

class Session;

class RestoreException final : public Exception {
public:
    explicit RestoreException(std::string_view what) : Exception(what) {}
};

class RestoreCallback {
public:
    RestoreCallback() = default;
    virtual ~RestoreCallback() = default;

    virtual void onProgress(float progress) = 0;
};

class Restore final {
    friend class Session;

private:
    const Session& mSession;
    etRestore* mPtr;

protected:
    Restore(const Session& session, etRestore* ptr);

public:
    ~Restore();
    Restore(const Restore&) = delete;
    Restore(Restore&&) noexcept = delete;
    Restore& operator=(const Restore&) = delete;
    Restore& operator=(Restore&& rhs) noexcept = delete;

    void start(RestoreCallback& cb);

    void cancel();

    std::filesystem::path getBackupPath() const;

private:
    template<class F>
    void wrapCCall(F func);

    template<class F>
    void wrapCCall(F func) const;
};
} // namespace etcpp
