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
#include <string>
#include <string_view>

namespace etcpp {

class Exception : public std::exception {
   protected:
    std::string mWhat;

   public:
    inline explicit Exception(std::string_view what) : mWhat(what) {}
    ~Exception() override = default;

    Exception(const Exception&) = default;
    Exception(Exception&&) = default;
    Exception& operator=(const Exception&) = default;
    Exception& operator=(Exception&&) = default;

    [[nodiscard]] const char* what() const noexcept override { return mWhat.c_str(); }
};

class CancelledException final : public Exception {
   public:
    inline CancelledException() : Exception("Operation Cancelled") {}
};
}    // namespace etcpp