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

#include <catch2/catch_test_macros.hpp>

#include <etsession.hpp>
#include "gpa_server.hpp"

TEST_CASE("SessionLogin") {
    GPAServer server;

    const char* userEmail = "hello@bar.com";
    const char* userPassword = "12345";

    std::string addrID;
    const auto userID = server.createUser(userEmail, userPassword, addrID);
    const auto url = server.url();

    auto session = etcpp::Session(url.c_str());
    {
        auto loginState = session.getLoginState();
        REQUIRE(loginState == etcpp::Session::LoginState::LoggedOut);
    }

    auto loginState = session.login(userEmail, userPassword);
    REQUIRE(loginState == etcpp::Session::LoginState::LoggedIn);
}
