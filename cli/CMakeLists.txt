#
# CLI Executable
#

add_executable(cli)

find_package(cxxopts CONFIG REQUIRED)
find_package(fmt CONFIG REQUIRED)

target_link_libraries(cli PRIVATE etcpp cxxopts::cxxopts fmt::fmt)

target_sources(cli PRIVATE
        bin/main.cpp
)

apply_cpp_flags(cli)

copy_shared_libraries(cli etcore)

install(TARGETS cli)
