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

#include "etexport_mail.hpp"

#include <etcore.h>
#include "etsession.hpp"

namespace etcpp {

inline etExportMailCallbackReply mapToETExportMailCallbackReply(ExportMailCallback::Reply r) {
    switch (r) {
        case ExportMailCallback::Reply::Continue:
            return ET_EXPORT_MAIL_CALLBACK_REPLY_CONTINUE;
        case ExportMailCallback::Reply::Cancel:
            return ET_EXPORT_MAIL_CALLBACK_REPLY_CANCEL;
        default:
            return ET_EXPORT_MAIL_CALLBACK_REPLY_CONTINUE;
    }
}

inline void mapETExportMailStatusToException(etExportMail* ptr, etExportMailStatus status) {
    switch (status) {
        case ET_EXPORT_MAIL_STATUS_INVALID:
            throw SessionException("Invalid instance");
        case ET_EXPORT_MAIL_STATUS_ERROR: {
            const char* lastErr = etExportMailGetLastError(ptr);
            if (lastErr == nullptr) {
                lastErr = "unknown";
            }
            throw ExportMailException(lastErr);
        }
        case ET_EXPORT_MAIL_STATUS_OK:
            break;
    }
}

etExportMailCallbacks makeETCallback(ExportMailCallback& cb) {
    return etExportMailCallbacks{
        .ptr = &cb,
        .onProgress =
            [](void* p, float progress) {
                return mapToETExportMailCallbackReply(
                    reinterpret_cast<ExportMailCallback*>(p)->onProgress(progress));
            },
    };
}

ExportMail::ExportMail(const etcpp::Session& session, etExportMail* ptr)
    : mSession(session), mPtr(ptr) {}

ExportMail::~ExportMail() {
    if (mPtr != nullptr) {
        wrapCCall([](etExportMail* ptr) { return etExportMailDelete(ptr); });
    }
}

void ExportMail::start(ExportMailCallback& cb) {
    wrapCCall([&](etExportMail* ptr) {
        auto etCb = makeETCallback(cb);
        return etExportMailStart(ptr, &etCb);
    });
}

template <class F>
void ExportMail::wrapCCall(F func) {
    static_assert(std::is_invocable_r_v<etExportMailStatus, F, etExportMail*>,
                  "invalid function/lambda signature");
    mapETExportMailStatusToException(mPtr, func(mPtr));
}

template <class F>
void ExportMail::wrapCCall(F func) const {
    static_assert(std::is_invocable_r_v<etExportMailStatus, F, etExportMail*>,
                  "invalid function/lambda signature");
    mapETExportMailStatusToException(mPtr, func(mPtr));
}

ExportMailException::ExportMailException(std::string_view what) : mWhat(what) {}

const char* ExportMailException::what() const noexcept {
    return mWhat.c_str();
}
}    // namespace etcpp