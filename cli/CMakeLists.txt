#
# CLI Executable
#

add_executable(cli)

target_link_libraries(cli PRIVATE etcpp)

target_sources(cli PRIVATE
        bin/main.cpp
)

apply_cpp_flags(cli)

find_package(ftxui CONFIG REQUIRED)

target_link_libraries(cli PRIVATE ftxui::dom ftxui::screen ftxui::component)

copy_shared_libraries(cli etcore)

install(TARGETS cli)