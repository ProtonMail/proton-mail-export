include(cmake/vcpkg_setup.cmake)
cmake_minimum_required(VERSION 3.23)
project(ExportTool CXX C)

set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}/cmake)
if (UNIX AND NOT APPLE)
    set(CMAKE_INSTALL_RPATH "$\{ORIGIN\}")
endif()

enable_testing()

include(config)
include(compile_options)
include(clang_tidy)
include(clang_format)
include(shared_lib_copy)
include(mingw)

add_subdirectory(go-lib)
add_subdirectory(lib)
add_subdirectory(cli)
