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

#include <etsession.hpp>
#include <string>
#include <type_traits>

#include "tasks/task.hpp"

template<class R>
class SessionTask : public Task<R> {
protected:
    std::string mDesc;
    etcpp::Session& mSession;

protected:
    SessionTask(etcpp::Session& session, std::string_view desc) : mDesc(desc), mSession(session) {}

public:
    virtual ~SessionTask() override = default;

    void cancel() override { mSession.cancel(); }

    std::string_view description() const override { return mDesc; }
};

template<class F>
class LoginSessionTask final : public SessionTask<etcpp::Session::LoginState> {
    static_assert(std::is_invocable_r_v<etcpp::Session::LoginState, F, etcpp::Session&>);
    F mFn;

public:
    LoginSessionTask(etcpp::Session& session, std::string_view desc, F&& f) : SessionTask<etcpp::Session::LoginState>(session, desc), mFn(f) {}

    virtual ~LoginSessionTask() override = default;

    etcpp::Session::LoginState run() override { return mFn(mSession); }
};
