
add_library(etcpp STATIC include/etcpp.hpp)

target_sources(etcpp PRIVATE
        include/etcpp.hpp
        lib/etcpp.cpp
)


set(config_file_dir "${CMAKE_CURRENT_BINARY_DIR}/config")
set(config_file "${config_file_dir}/etconfig.hpp")
set(config_file_in "${CMAKE_CURRENT_SOURCE_DIR}/cmake/etconfig.hpp.in")

configure_file("${config_file_in}" "${config_file}" @ONLY)

target_link_libraries(etcpp PRIVATE etcore)

target_include_directories(etcpp PUBLIC include ${config_file_dir})

apply_cpp_flags(etcpp)

add_subdirectory(tests)