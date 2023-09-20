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
#include <string_view>

#include <fmt/format.h>
#include <optional>
#include "etexception.hpp"

#if defined(_WIN32)
#define ET_LOG_DECL __declspec(dllexport)
#else
#define ET_LOG_DECL
#endif

extern "C" {
void ET_LOG_DECL etLogInfo(const char*, const char*);
void ET_LOG_DECL etLogDebug(const char*, const char*);
void ET_LOG_DECL etLogWarn(const char*, const char*);
void ET_LOG_DECL etLogError(const char*, const char*);
}

namespace etcpp {

std::string& getThreadLocalLogBuffer();

constexpr const char* kLogTag = "etcpp";

class LogException : public Exception {
   public:
    explicit LogException(std::string_view w) : Exception(w) {}
};

class LogScope final {
   public:
    explicit LogScope(const std::filesystem::path& p);
    ~LogScope();

    LogScope(const LogScope&) = delete;
    LogScope(LogScope&&) = delete;
    LogScope& operator=(const LogScope&) = delete;
    LogScope operator=(LogScope&&) = delete;

    std::optional<std::filesystem::path> getLogPath() const;
};

template <typename... T>
inline void logInfo(fmt::format_string<T...> fmt, T&&... args) {
    auto& buffer = getThreadLocalLogBuffer();
    fmt::format_to(std::back_inserter(buffer), fmt, args...);
    etLogInfo(kLogTag, buffer.c_str());
}

template <typename... T>
inline void logDebug(fmt::format_string<T...> fmt, T&&... args) {
    auto& buffer = getThreadLocalLogBuffer();
    fmt::format_to(std::back_inserter(buffer), fmt, args...);
    etLogDebug(kLogTag, buffer.c_str());
}

template <typename... T>
inline void logError(fmt::format_string<T...> fmt, T&&... args) {
    auto& buffer = getThreadLocalLogBuffer();
    fmt::format_to(std::back_inserter(buffer), fmt, args...);
    etLogError(kLogTag, buffer.c_str());
}

template <typename... T>
inline void logWarn(fmt::format_string<T...> fmt, T&&... args) {
    auto& buffer = getThreadLocalLogBuffer();
    fmt::format_to(std::back_inserter(buffer), fmt, args...);
    etLogWarn(kLogTag, buffer.c_str());
}

}    // namespace etcpp