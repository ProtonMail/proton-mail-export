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

#include <mach-o/dyld.h>
#include <climits>

namespace etcpp {
std::filesystem::path getExecutablePath() {
    char rawPathName[PATH_MAX];
    char realPathName[PATH_MAX];
    auto rawPathSize = (uint32_t)sizeof(rawPathName);

    if (_NSGetExecutablePath(rawPathName, &rawPathSize) != 0) {
        throw std::runtime_error("failed to extract nsexecutable path");
    }
    if (realpath(rawPathName, realPathName) == nullptr) {
        throw std::runtime_error(strerror(errno));
    }

    return std::filesystem::u8path(realPathName);
}
}    // namespace etcpp