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

file(GLOB_RECURSE go_files
    LIST_DIRECTORIES false
    CONFIGURE_DEPENDS
    "${CMAKE_CURRENT_SOURCE_DIR}/*.go"
)


set(ADDITIONAL_FLAGS)

if (WIN32)
    # Unfortunately for us, c-shared dlls produced by go don't produce a .lib file
    # so we have to generate one ourselfs... sigh.
    set(GO_DEF ${CMAKE_CURRENT_BINARY_DIR}/etcore.def)
    find_program(gendef-bin gendef PATHS "${MINGW_PATH}")
    if (NOT gendef-bin)
        message(FATAL_ERROR "could not find gendef binary")
else()
	message(STATUS "gendef-bin=${gendef-bin}")
    endif()
    add_custom_target(go-lib-gendef
        COMMAND gendef ${GO_LIB}
        BYPRODUCTS ${GO_DEF}
    )

    find_program(dlltool-bin dlltool PATHS "${MINGW_PATH}")
    if (NOT dlltool-bin)
        message(FATAL_ERROR "could not find dlltool, please install with 'pacman -S binutils'")
    else()
	    message(STATUS "dlltool-bin=${dlltool-bin}")
    endif()
    add_custom_target(go-lib-gen-implib
        DEPENDS ${GO_DEF}
        COMMAND ${dlltool-bin} -d ${GO_DEF} -l ${GO_LIB_WIN32_IMPL}
        BYPRODUCTS ${GO_LIB_WIN32_IMPL}
    )

    add_dependencies(go-lib-gendef go-lib-build)
    add_dependencies(go-lib-gen-implib go-lib-gendef)
endif()

if (UNIX AND NOT APPLE)
    list(APPEND ADDITIONAL_FLAGS -ldflags '-extldflags -Wl,-soname,etcore${CMAKE_SHARED_LIBRARY_SUFFIX}')
endif()

set(ENV_MANIP)

set(GO_LIB_BUILD_EXTRA)

if (WIN32)	
    set(ENV_MANIP ${CMAKE_COMMAND} -E env "CGO_ENABLED=1" "GOOS=windows" "PATH=${MINGW_PATH}")
	set(GO_LIB_BUILD_EXTRA COMMAND ${CMAKE_COMMAND} -P "${CMAKE_CURRENT_SOURCE_DIR}/cmake/patch_etcore.cmake" "${CMAKE_CURRENT_BINARY_DIR}/etcore.h")
endif()

add_custom_command (OUTPUT ${GO_LIB}
        DEPENDS ${go_files}
        COMMAND ${ENV_MANIP} "${go_bin}" build -a -tags="gpa_server" ${ADDITIONAL_FLAGS} -buildmode=c-shared  -o "${GO_LIB}"
        export.go export_gpa.go export_mail.go
        ${GO_LIB_BUILD_EXTRA}
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(go-lib-build
    DEPENDS ${GO_LIB}
)

add_library(etcore SHARED IMPORTED GLOBAL)
add_dependencies(etcore go-lib-build)

if (WIN32)
	add_dependencies(etcore go-lib-gen-implib)
endif()

configure_file(${CMAKE_CURRENT_SOURCE_DIR}/cmake/constants.go.in
    ${CMAKE_CURRENT_SOURCE_DIR}/internal/constants.go
    @ONLY
)

set_target_properties(etcore
    PROPERTIES
    IMPORTED_LOCATION ${GO_LIB}
    IMPORTED_RUNTIME ${GO_LIB}
    IMPORTED_IMPLIB ${GO_LIB_WIN32_IMPL}
    INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_BINARY_DIR};${CMAKE_CURRENT_SOURCE_DIR}/cgo_headers"
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

add_custom_target(go-lib-mockgen
    COMMAND mockgen
        -self_package "github.com/ProtonMail/export-tool/internal/apiclient"
        -package apiclient github.com/ProtonMail/export-tool/internal/apiclient Builder,Client > "${CMAKE_BINARY_DIR}/mocks_test.go"
    COMMAND ${CMAKE_COMMAND} -E copy_if_different "${CMAKE_BINARY_DIR}/mocks_test.go" "${CMAKE_CURRENT_SOURCE_DIR}/internal/apiclient/mocks.go"
    COMMAND mockgen
        -self_package "github.com/ProtonMail/export-tool/internal/mail"
        -package mail github.com/ProtonMail/export-tool/internal/mail StageErrorReporter > "${CMAKE_BINARY_DIR}/mocks_mail_test.go"
    COMMAND ${CMAKE_COMMAND} -E copy_if_different "${CMAKE_BINARY_DIR}/mocks_mail_test.go" "${CMAKE_CURRENT_SOURCE_DIR}/internal/mail/mocks_test.go"
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)
