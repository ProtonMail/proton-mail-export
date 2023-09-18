#
# CLI Executable
#

add_executable(cli)

find_package(cxxopts CONFIG REQUIRED)
find_package(fmt CONFIG REQUIRED)

target_link_libraries(cli PRIVATE etcpp cxxopts::cxxopts)

target_include_directories(cli PRIVATE "${CMAKE_CURRENT_SOURCE_DIR}/bin")

target_sources(cli PRIVATE
        bin/main.cpp
        bin/tui_util.cpp
        bin/tui_util.hpp
        bin/tasks/task.hpp
        bin/tasks/mail_task.cpp
        bin/tasks/mail_task.hpp
)

apply_cpp_flags(cli)

copy_shared_libraries(cli etcore)

install(TARGETS cli)
