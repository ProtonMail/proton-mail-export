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
#include <optional>

namespace etcpp {

class GlobalScope final {
   public:
    explicit GlobalScope(const std::filesystem::path& p, void (*onRecover)());
    ~GlobalScope();

    GlobalScope(const GlobalScope&) = delete;
    GlobalScope(GlobalScope&&) = delete;
    GlobalScope& operator=(const GlobalScope&) = delete;
    GlobalScope operator=(GlobalScope&&) = delete;

    std::optional<std::filesystem::path> getLogPath() const;

    static void reportMessage(const char* tag, const char*);
    static void reportError(const char* tag, const char*);
};

}    // namespace etcpp