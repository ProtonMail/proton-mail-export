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

#include <memory>
#include <string>

#include "etexception.hpp"
#include "etexport_mail.hpp"

extern "C" {
struct etSession;
}

namespace etcpp {

class SessionException final : public Exception {
   public:
    explicit SessionException(std::string_view what) : Exception(what) {}
};

class SessionCallback {
   public:
    virtual ~SessionCallback() = default;
    virtual void onNetworkRestored() = 0;
    virtual void onNetworkLost() = 0;
};

class Session final {
   private:
    etSession* mPtr;
    std::shared_ptr<SessionCallback> mCallbacks;

   public:
    enum class LoginState {
        LoggedOut,
        AwaitingTOTP,
        AwaitingHV,
        AwaitingMailboxPassword,
        LoggedIn
    };

    inline explicit Session(const char* serverURL) : Session(serverURL, {}) {}
    explicit Session(const char* serverURL, const std::shared_ptr<SessionCallback>& mCallbacks);
    ~Session();
    Session(const Session&) = delete;
    Session(Session&&) noexcept;
    Session& operator=(const Session&) = delete;
    Session& operator=(Session&& rhs) noexcept;

    [[nodiscard]] LoginState login(const char* email, std::string_view password);
    [[nodiscard]] LoginState loginTOTP(const char* totp);
    [[nodiscard]] LoginState loginMailboxPassword(std::string_view password);

    [[nodiscard]] LoginState getLoginState() const;

    [[nodiscard]] ExportMail newExportMail(const char* exportPath) const;

    void cancel();

   private:
    template <class F>
    void wrapCCall(F func);

    template <class F>
    void wrapCCall(F func) const;
};
}    // namespace etcpp