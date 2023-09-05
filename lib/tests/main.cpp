
#include <catch2/catch_test_macros.hpp>

#include <etcpp.hpp>
#include "gpa_server.hpp"

TEST_CASE("SessionLogin") {
    GPAServer server;

    const char* userEmail = "hello@bar.com";
    const char* userPassword = "12345";

    const auto userID = server.createUser(userEmail, userPassword);
    const auto url = server.url();

    auto session = etcpp::Session(url.c_str());
    {
        auto loginState = session.getLoginState();
        REQUIRE(loginState == etcpp::Session::LoginState::LoggedOut);
    }

    auto loginState = session.login(userEmail, userPassword);
    REQUIRE(loginState == etcpp::Session::LoginState::LoggedIn);
}