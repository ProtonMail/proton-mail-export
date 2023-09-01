
add_library(etcpp STATIC include/etcpp.hpp)

target_sources(etcpp PRIVATE
        include/etcpp.hpp
        lib/etcpp.cpp
)

target_link_libraries(etcpp PRIVATE etcore)

target_include_directories(etcpp PUBLIC include)

apply_cpp_flags(etcpp)

add_subdirectory(tests)