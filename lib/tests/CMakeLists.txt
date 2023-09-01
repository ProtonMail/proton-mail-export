#
# etcpp unit tests
#

find_package(Catch2 CONFIG REQUIRED)

add_executable(etcpp_test)

target_sources(etcpp_test PRIVATE
    main.cpp
)

target_link_libraries(etcpp_test PRIVATE etcpp Catch2::Catch2 Catch2::Catch2WithMain)

add_test(NAME etcpp-test
    COMMAND etcpp_test
)

copy_shared_libraries(etcpp_test etcore)
