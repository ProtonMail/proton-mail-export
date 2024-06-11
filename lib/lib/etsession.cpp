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

#include <etsession.hpp>

#include <etconfig.hpp>
#include <proton-mail-export.h>

namespace etcpp {

Session::LoginState mapLoginState(etSessionLoginState s);

inline void mapETStatusToException(etSession* ptr, etSessionStatus status) {
    switch (status) {
    case ET_SESSION_STATUS_INVALID:
        throw SessionException("Invalid instance");
    case ET_SESSION_STATUS_ERROR:
    {
        const char* lastErr = etSessionGetLastError(ptr);
        if (lastErr == nullptr) {
            lastErr = "unknown";
        }
        throw SessionException(lastErr);
    }
    case ET_SESSION_STATUS_CANCELLED:
        throw CancelledException();
    case ET_SESSION_STATUS_OK:
        break;
    }
}

etSessionCallbacks makeCCallback(SessionCallback* ptr) {
    etSessionCallbacks cb{};

    if (ptr == nullptr) {
        return cb;
    }

    cb.ptr = ptr;
    cb.onNetworkLost = [](void* p) { reinterpret_cast<SessionCallback*>(p)->onNetworkLost(); };
    cb.onNetworkRestored = [](void* p) { reinterpret_cast<SessionCallback*>(p)->onNetworkRestored(); };

    return cb;
}

Session::Session(const char* serverURL, const std::shared_ptr<SessionCallback>& callbacks) : mCallbacks(callbacks) {
    char* outErr = nullptr;
    mPtr = etSessionNew(serverURL, makeCCallback(mCallbacks.get()), &outErr);
    if (mPtr == nullptr) {
        auto ex = SessionException(outErr);
        etFree(outErr);
        throw std::move(ex);
    }
}

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

Session::LoginState Session::login(const char* email, std::string_view password) {
    LoginState ls = LoginState::LoggedOut;
    wrapCCall([&](etSession* ptr) {
        etSessionLoginState els = ET_SESSION_LOGIN_STATE_LOGGED_OUT;
        auto status = etSessionLogin(ptr, email, password.data(), int(password.length()), &els);
        if (status == ET_SESSION_STATUS_OK) {
            ls = mapLoginState(els);
        }
        return status;
    });

    return ls;
}

Session::LoginState Session::loginTOTP(const char* totp) {
    LoginState ls = LoginState::LoggedOut;
    wrapCCall([&](etSession* ptr) {
        etSessionLoginState els = ET_SESSION_LOGIN_STATE_LOGGED_OUT;
        auto status = etSessionSubmitTOTP(ptr, totp, &els);
        if (status == ET_SESSION_STATUS_OK) {
            ls = mapLoginState(els);
        }
        return status;
    });

    return ls;
}

Session::LoginState Session::loginMailboxPassword(std::string_view password) {
    LoginState ls = LoginState::LoggedOut;
    wrapCCall([&](etSession* ptr) {
        etSessionLoginState els = ET_SESSION_LOGIN_STATE_LOGGED_OUT;
        auto status = etSessionSubmitMailboxPassword(ptr, password.data(), int(password.length()), &els);
        if (status == ET_SESSION_STATUS_OK) {
            ls = mapLoginState(els);
        }
        return status;
    });

    return ls;
}

Session::LoginState Session::getLoginState() const {
    LoginState ls = LoginState::LoggedOut;
    wrapCCall([&](etSession* ptr) {
        etSessionLoginState els = ET_SESSION_LOGIN_STATE_LOGGED_OUT;
        auto status = etSessionGetLoginState(ptr, &els);
        if (status == ET_SESSION_STATUS_OK) {
            ls = mapLoginState(els);
        }
        return status;
    });

    return ls;
}

ExportBackup Session::newExportBackup(const char* exportPath) const {
    etExportBackup* exportPtr = nullptr;
    wrapCCall([&](etSession* ptr) -> etSessionStatus { return etSessionNewExportBackup(ptr, exportPath, &exportPtr); });

    return ExportBackup(*this, exportPtr);
}

ExportRestore Session::newExportRestore(const char* backupPath) const {
    etExportRestore* restorePtr = nullptr;
    wrapCCall([&](etSession* ptr) -> etSessionStatus { return etSessionNewExportRestore(ptr, backupPath, &restorePtr); });

    return ExportRestore(*this, restorePtr);
}

void Session::cancel() {
    wrapCCall([](etSession* ptr) { return etSessionCancel(ptr); });
}
std::string Session::getEmail() const {
    char* outEmail = nullptr;
    wrapCCall([&](etSession* ptr) -> etSessionStatus { return etSessionGetEmail(ptr, &outEmail); });

    auto result = std::string(outEmail);
    etFree(outEmail);

    return result;
}

std::string Session::getHVSolveURL() const {
    char* outURL = nullptr;
    wrapCCall([&](etSession* ptr) -> etSessionStatus { return etSessionGetHVSolveURL(ptr, &outURL); });

    auto result = std::string(outURL);
    etFree(outURL);

    return result;
}

Session::LoginState Session::markHVSolved() {
    LoginState ls = LoginState::LoggedOut;
    wrapCCall([&](etSession* ptr) -> etSessionStatus {
        etSessionLoginState els = ET_SESSION_LOGIN_STATE_LOGGED_OUT;
        auto status = etSessionMarkHVSolved(ptr, &els);
        if (status == ET_SESSION_STATUS_OK) {
            ls = mapLoginState(els);
        }
        return status;
    });

    return ls;
}

template<class F>
void Session::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*>, "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr));
}

template<class F>
void Session::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etSessionStatus, F, etSession*>, "invalid function/lambda signature");
    mapETStatusToException(mPtr, func(mPtr));
}
Session::LoginState mapLoginState(etSessionLoginState s) {
    switch (s) {
    case ET_SESSION_LOGIN_STATE_LOGGED_OUT:
        return Session::LoginState::LoggedOut;
    case ET_SESSION_LOGIN_STATE_AWAITING_HV:
        return Session::LoginState::AwaitingHV;
    case ET_SESSION_LOGIN_STATE_AWAITING_MAILBOX_PASSWORD:
        return Session::LoginState::AwaitingMailboxPassword;
    case ET_SESSION_LOGIN_STATE_LOGGED_IN:
        return Session::LoginState::LoggedIn;
    case ET_SESSION_LOGIN_STATE_AWAITING_TOTP:
        return Session::LoginState::AwaitingTOTP;
    }

    return Session::LoginState::LoggedOut;
}

} // namespace etcpp
