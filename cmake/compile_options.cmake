#
# CPP compiler options
#

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")

function(apply_cpp_flags target)
    if (NOT MSVC)
        target_compile_options(${target} PUBLIC
                -Wall
                -Wextra
                -Werror
                -Wdouble-promotion
                -Wshadow
                -Wformat=2
                -Wcast-align
                -Wsign-compare
                -Wno-float-equal
                -Wreturn-type
                -Wunused-variable
                -Wno-error=attributes
        )
    else()
        target_compile_definitions(${target} PUBLIC
                /WX
                /permissive-
                /EHsc
                /utf8
        )

        target_compile_definitions(${target} PUBLIC _UNICODE UNICODE)
    endif()

endfunction()
