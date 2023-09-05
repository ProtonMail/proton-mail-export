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

#include <etcpp.hpp>

#include <etcore.h>

namespace etcpp {

inline void mapETStatusToException(etSession* ptr, etSessionStatus status) {
    switch (status) {
        case ET_SESSION_STATUS_INVALID:
            throw Exception("Invalid instance");
        case ET_SESSION_STATUS_ERROR: {
            const char* lastErr = etSessionGetLastError(ptr);
            if (lastErr == nullptr) {
                lastErr = "unknown";
            }
            throw Exception(lastErr);
        }
        case ET_SESSION_STATUS_OK:
            break;
    }
}

Session::Session() : mPtr(etSessionNew()) {}

Session::~Session() {
    if (mPtr != nullptr) {
        wrapCCall([](etSession* ptr) -> etSessionStatus { return etSessionDelete(ptr); });
    }
}

Session::Session(Session&& rhs) noexcept : mPtr(rhs.mPtr) {
    rhs.mPtr = nullptr;
}

Session& Session::operator=(Session&& rhs) noexcept {
    if (this != &rhs) {
        if (mPtr != nullptr) {
            wrapCCall([](etSession* ptr) -> etSessionStatus { return etSessionDelete(ptr); });
        }

        mPtr = rhs.mPtr;
        rhs.mPtr = nullptr;
    }

    return *this;
}

std::string Session::hello() const {
    std::string result;
    wrapCCallOut(result, [](etSession* ptr, std::string& out) {
        char* cOut;
        auto status = etSessionHello(ptr, &cOut);
        if (status != ET_SESSION_STATUS_OK) {
            return status;
        }

        out.assign(cOut);
        free(cOut);

        return status;
    });

    return result;
}

void Session::helloError() const {
    wrapCCall([](etSession* ptr) { return etSessionHelloError(ptr); });
}

const char* Exception::what() const noexcept {
    return mWhat.c_str();
}

Exception::Exception(std::string_view what) : mWhat(what) {}

template <class F>
void Session::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*>,
                  "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr));
}

template <class F>
void Session::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*>,
                  "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr));
}

template <class F, class OUT>
void Session::wrapCCallOut(OUT& out, F func) {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*, OUT&>,
                  "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr, out));
}

template <class F, class OUT>
void Session::wrapCCallOut(OUT& out, F func) const {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*, OUT&>,
                  "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr, out));
}
}    // namespace etcpp
