#
# CLI Executable
#

add_executable(proton-export-cli)

find_package(cxxopts CONFIG REQUIRED)
find_package(fmt CONFIG REQUIRED)

target_link_libraries(proton-export-cli PRIVATE etcpp cxxopts::cxxopts)

target_include_directories(proton-export-cli PRIVATE "${CMAKE_CURRENT_SOURCE_DIR}/bin")

target_sources(proton-export-cli PRIVATE
        bin/main.cpp
        bin/tui_util.cpp
        bin/tui_util.hpp
        bin/task_runner.hpp
        bin/tasks/task.hpp
        bin/tasks/mail_task.cpp
        bin/tasks/mail_task.hpp
        bin/tasks/session_task.hpp
)


apply_cpp_flags(proton-export-cli)

copy_shared_libraries(proton-export-cli etcore)

if (WIN32)
    # For Win32 RC file
    set(ET_CLI_FULL_NAME "Proton Export CLI")
    set(ET_CLI_NAME "Proton Export")

    set(CLI_RC_FILE "${CMAKE_CURRENT_BINARY_DIR}/resources.rc")

    configure_file(
        ${CMAKE_CURRENT_SOURCE_DIR}/cmake/resources.rc.cmake
        ${CLI_RC_FILE}
        @ONLY
    )

    target_sources(proton-export-cli PRIVATE ${CLI_RC_FILE})

endif()

install(TARGETS proton-export-cli)
