#
# clang format setup
#

find_program(clang_format_bin clang-format)

if (WIN32 OR MINGW)
    return()
endif()

if (NOT clang_format_bin)
    message(FATAL_ERROR "could not find clang-format binary")
endif ()

set(clang_format_dirs
    ${CMAKE_SOURCE_DIR}/cli
    ${CMAKE_SOURCE_DIR}/lib
)

set(clang_format_files "")

foreach (dir IN LISTS clang_format_dirs)
    file(GLOB_RECURSE files "${dir}/*.cpp" "${dir}*.h" "${dir}/*.hpp" "${dir}/*.c")
    list(APPEND clang_format_files ${files})
endforeach ()


# Target to format all cpp source files with clang-format.
add_custom_target(
    clang-format
    COMMAND "${clang_format_bin}"
    -i
    --verbose
    ${clang_format_files}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Running clang-format"
)

# Target to verify all CPP files are properly formatted.
add_custom_target(
        clang-format-check
        COMMAND "${clang_format_bin}"
        --dry-run
        --Werror
        --verbose
        ${clang_format_files}
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
        COMMENT "Running clang-format check"
)
