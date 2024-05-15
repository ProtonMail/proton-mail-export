#
# CLI Executable
#

add_executable(proton-mail-export-cli)

find_package(cxxopts CONFIG REQUIRED)
find_package(fmt CONFIG REQUIRED)

target_link_libraries(proton-mail-export-cli PRIVATE etcpp cxxopts::cxxopts)

target_include_directories(proton-mail-export-cli PRIVATE "${CMAKE_CURRENT_SOURCE_DIR}/bin")

target_sources(proton-mail-export-cli PRIVATE
        bin/main.cpp
        bin/operation.h
        bin/operation.cpp
        bin/tui_util.cpp
        bin/tui_util.hpp
        bin/task_runner.hpp
        bin/tasks/task.hpp
        bin/tasks/mail_task.cpp
        bin/tasks/mail_task.hpp
        bin/tasks/session_task.hpp
        bin/tasks/global_task.cpp
        bin/tasks/global_task.hpp
)

if (APPLE)
    target_sources(proton-mail-export-cli PRIVATE
        bin/macos.hpp
        bin/macos.mm
    )

    target_link_libraries(proton-mail-export-cli PRIVATE "-framework Foundation")
endif()


apply_cpp_flags(proton-mail-export-cli)

copy_shared_libraries(proton-mail-export-cli etcore)

if (WIN32)
    # For Win32 RC file
    set(ET_CLI_FULL_NAME "Proton Mail Export CLI")
    set(ET_CLI_NAME "Proton Mail Export")

    set(CLI_RC_FILE "${CMAKE_CURRENT_BINARY_DIR}/resources.rc")

    configure_file(
        ${CMAKE_CURRENT_SOURCE_DIR}/cmake/resources.rc.cmake
        ${CLI_RC_FILE}
        @ONLY
    )

    target_sources(proton-mail-export-cli PRIVATE ${CLI_RC_FILE})

endif()

install(TARGETS proton-mail-export-cli)

function(gen_meta_file input output)
    configure_file(
        "${input}"
        "${output}.gen"
        @ONLY
    )

    file(GENERATE
        OUTPUT "${output}"
        INPUT "${output}.gen"
    )

    install(FILES ${output} DESTINATION "meta")
endfunction()

if (APPLE)
    set(BUNDLE_APP_NAME "$<TARGET_FILE_NAME:proton-mail-export-cli>")
    set(EXE_TO_LAUNCH "${BUNDLE_APP_NAME}")
    set(SCRIPT_TO_LAUNCH "${BUNDLE_APP_NAME}.sh")
    set(BUNDLE_IDENTIFIER "ch.protonmail.export-cli")
    set(BUNDLE_EXEC "launcher.sh")

    gen_meta_file(
        "${PROJECT_SOURCE_DIR}/cli/cmake/Info.plist.in"
        "${PROJECT_BINARY_DIR}/Info.plist"
    )

    gen_meta_file(
        "${PROJECT_SOURCE_DIR}/cli/cmake/${BUNDLE_EXEC}.in"
        "${PROJECT_BINARY_DIR}/${BUNDLE_EXEC}"
    )

    gen_meta_file(
        "${PROJECT_SOURCE_DIR}/cli/cmake/run.sh.in"
        "${PROJECT_BINARY_DIR}/${SCRIPT_TO_LAUNCH}"
    )

    install(FILES ${PROJECT_SOURCE_DIR}/res/icon.icns DESTINATION "meta")
endif()
