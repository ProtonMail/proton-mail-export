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

#include <filesystem>
#include <string>

namespace etcpp {

/// Get the current executable path, may throw on failure.
std::filesystem::path getExecutablePath();

std::filesystem::path getExecutableDir();

/// May throw runtime_error if the expanded arguments result in multiple values or if there was an
/// error;
std::filesystem::path expandCLIPath(const std::filesystem::path& path);
} // namespace etcpp
