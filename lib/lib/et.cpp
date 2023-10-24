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

#include "et.hpp"
#include "etexception.hpp"

#include <proton-mail-export.h>

namespace etcpp {
GlobalScope::GlobalScope(const std::filesystem::path& path, void (*onRecover)()) {
    auto cpath = path.u8string();
    if (etInit(cpath.c_str(), onRecover) != 0) {
        const char* lastErr = etGetLastError();
        if (lastErr == nullptr) {
            lastErr = "unknown error";
        }

        throw Exception(lastErr);
    }
}

GlobalScope::~GlobalScope() {
    etClose();
}

std::optional<std::filesystem::path> GlobalScope::getLogPath() const {
    const char* clogPath = etLogGetPath();
    if (clogPath == nullptr) {
        return {};
    }

    return std::filesystem::u8path(clogPath);
}

void GlobalScope::reportMessage(const char* tag, const char* msg) {
    etReportMessage(tag, msg);
}

void GlobalScope::reportError(const char* tag, const char* msg) {
    etReportError(tag, msg);
}

bool GlobalScope::newVersionAvailable() const {
    return etNewVersionAvailable() == 1;
}

}    // namespace etcpp