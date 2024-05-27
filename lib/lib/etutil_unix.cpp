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
#include <wordexp.h>

namespace etcpp {

std::filesystem::path expandCLIPath(const std::filesystem::path& path) {
    auto value = path.u8string();

    wordexp_t p;

    if (wordexp(value.c_str(), &p, 0) != 0) {
        throw std::runtime_error(fmt::format("failed to expand '{}'", value));
    }

    if (p.we_wordc > 1) {
        wordfree(&p);
        throw std::runtime_error(fmt::format("'{}' expands into more than one value", value));
    }

    auto result = std::filesystem::u8path(p.we_wordv[0]);
    wordfree(&p);

    return result;
}

} // namespace etcpp
