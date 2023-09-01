#
# clang tidy setup
#

find_program(clang_tidy_bin clang-tidy)

if (WIN32 OR MINGW)
    return()
endif()

if (NOT clang_tidy_bin)
  message(FATAL_ERROR "Could not find clang-tidy on your system")
endif()

if (clang_tidy_bin)

set(CMAKE_CXX_CLANG_TIDY
   ${clang_tidy_bin};-checks=-*,performance-*,portability-*,-readability-*,-performance-unnecessary-value-param,clang-analyzer-*;
)

endif()
