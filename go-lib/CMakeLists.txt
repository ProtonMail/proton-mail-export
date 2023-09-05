#
# Go library which provides the common C API.
#

find_program(go_bin go)

if (NOT go_bin)
    message(FATAL_ERROR "Could not find go binary")
endif()

set(GO_LIB_PATH "${CMAKE_CURRENT_BINARY_DIR}")
set(GO_LIB_PATH "${CMAKE_CURRENT_BINARY_DIR}" PARENT_SCOPE)
set(GO_LIB "${CMAKE_CURRENT_BINARY_DIR}/etcore${CMAKE_SHARED_LIBRARY_SUFFIX}")
set(GO_LIB_WIN32_IMPL "${CMAKE_CURRENT_BINARY_DIR}/etcore.lib")

file(GLOB go_files
    LIST_DIRECTORIES false
    CONFIGURE_DEPENDS
    "${CMAKE_CURRENT_SOURCE_DIR}/*.go"
)

set(ADDITIONAL_FLAGS)

if (MINGW)
    # Unfortunately for us, c-shared dlls produced by go don't produce a .lib file
    # so we have to generate one ourselfs... sigh.
    set(GO_DEF ${CMAKE_CURRENT_BINARY_DIR}/etcore.def)
    find_program(gendef-bin gendef)
    if (NOT gendef-bin)
        message(FATAL_ERROR "could not find gendef binary")
    endif()
    add_custom_target(go-lib-gendef
        COMMAND gendef ${GO_LIB}
        BYPRODUCTS ${GO_DEF}
    )

    find_program(dlltool-bin dlltool)
    if (NOT dlltool-bin)
        message(FATAL_ERROR "could not find dlltool, please install with 'pacman -S binutils'")
    endif()
    add_custom_target(go-lib-gen-implib
        DEPENDS ${GO_DEF}
        COMMAND ${dlltool-bin} -d ${GO_DEF} -l ${GO_LIB_WIN32_IMPL}
        BYPRODUCTS ${GO_LIB_WIN32_IMPL}
    )

    add_dependencies(go-lib-gendef go-lib-build)
    
endif()

if (UNIX AND NOT APPLE)
    list(APPEND ADDITIONAL_FLAGS -ldflags '-extldflags -Wl,-soname,etcore${CMAKE_SHARED_LIBRARY_SUFFIX}')
endif()

add_custom_target(go-lib-build
    DEPENDS ${go_files}
    COMMAND "${go_bin}" build ${ADDITIONAL_FLAGS} -buildmode=c-shared  -o "${GO_LIB}"  export.go
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    BYPRODUCTS ${GO_LIB}
)

add_library(etcore SHARED IMPORTED GLOBAL)
add_dependencies(etcore go-lib-build)

set_target_properties(etcore
    PROPERTIES
    IMPORTED_LOCATION ${GO_LIB}
    IMPORTED_RUNTIME ${GO_LIB}
    IMPORTED_IMPLIB ${GO_LIB_WIN32_IMPL}
    INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR}
)

install(IMPORTED_RUNTIME_ARTIFACTS etcore DESTINATION bin)

add_test(NAME go-lib-test
    COMMAND go test -v ./...
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(go-lib-lint
    COMMAND golangci-lint run ./...
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)
