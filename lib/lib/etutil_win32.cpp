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

#include "etutil.hpp"

#include <fmt/format.h>
#include <userenv.h>
#include <windows.h>

namespace etcpp {
constexpr int WIN32_MAX_PATH = 4096;

std::filesystem::path getExecutablePath() {
    wchar_t rawPathName[WIN32_MAX_PATH];
    if (GetModuleFileNameW(NULL, rawPathName, WIN32_MAX_PATH) == 0) {
        throw std::runtime_error(fmt::format("failed to get executable path {:x}", GetLastError()));
    }
    return std::filesystem::path(rawPathName);
}

std::filesystem::path expandCLIPath(const std::filesystem::path& path) {
    wchar_t outBuffer[WIN32_MAX_PATH];

    const auto value = path.wstring();

    if (ExpandEnvironmentStringsForUserW(NULL, value.c_str(), outBuffer, WIN32_MAX_PATH) == FALSE) {
        throw std::runtime_error(fmt::format("failed to expand '{}'", path.u8string()));
    }

    return std::filesystem::path(outBuffer);
}

}    // namespace etcpp
