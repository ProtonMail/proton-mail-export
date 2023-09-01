
#include <catch2/catch_test_macros.hpp>

#include <etcpp.hpp>

TEST_CASE("TestHello") {
    auto session = etcpp::Session();

    REQUIRE(session.hello() == "Hello world");
}

TEST_CASE("TestHelloError") {
    auto session = etcpp::Session();

    REQUIRE_THROWS(session.helloError());
}