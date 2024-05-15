#
# Go library which provides the common C API.
#

find_program(go_bin go)

if (NOT go_bin)
    message(FATAL_ERROR "Could not find go binary")
endif()

set(GO_LIB_PATH "${CMAKE_CURRENT_BINARY_DIR}")
set(GO_LIB_PATH "${CMAKE_CURRENT_BINARY_DIR}" PARENT_SCOPE)
set(GO_LIB "${CMAKE_CURRENT_BINARY_DIR}/proton-mail-export${CMAKE_SHARED_LIBRARY_SUFFIX}")
set(GO_LIB_WIN32_IMPL "${CMAKE_CURRENT_BINARY_DIR}/proton-mail-export.lib")
set(GO_LIB_TEST "${CMAKE_CURRENT_BINARY_DIR}/etcore-test${CMAKE_SHARED_LIBRARY_SUFFIX}")
set(GO_LIB_TEST_WIN32_IMPL "${CMAKE_CURRENT_BINARY_DIR}/etcore-test.lib")

file(GLOB_RECURSE go_files
    LIST_DIRECTORIES false
    CONFIGURE_DEPENDS
    "${CMAKE_CURRENT_SOURCE_DIR}/*.go"
)

function(build_cgo_lib)
    cmake_parse_arguments(CGO ""
            "TARGET;NAME"
            "TAGS;GO_SOURCES;GO_EXPORTS"
            ${ARGN}
    )

    set(target ${CGO_TARGET})
    set(name ${CGO_NAME})

    set(ADDITIONAL_FLAGS)
    set(lib "${CMAKE_CURRENT_BINARY_DIR}/${name}${CMAKE_SHARED_LIBRARY_SUFFIX}")
    set(win32_imp_lib "${CMAKE_CURRENT_BINARY_DIR}/${name}.lib")

    if (UNIX AND NOT APPLE)
        list(APPEND ADDITIONAL_FLAGS -ldflags '-extldflags -Wl,-soname,${name}${CMAKE_SHARED_LIBRARY_SUFFIX}')
    endif()

    set(ENV_MANIP)

    set(GO_LIB_BUILD_EXTRA)

    if (WIN32)
        set(ENV_MANIP ${CMAKE_COMMAND} -E env "CGO_ENABLED=1" "GOOS=windows" "PATH=${MINGW_PATH}")
        set(GO_LIB_BUILD_EXTRA COMMAND ${CMAKE_COMMAND} -P "${CMAKE_CURRENT_SOURCE_DIR}/cmake/patch_etcore.cmake" "${CMAKE_CURRENT_BINARY_DIR}/${name}.h")
    endif()

    if (APPLE)
        set(go_arch)
        if ("${CMAKE_OSX_ARCHITECTURES}" STREQUAL "arm64")
            set(go_arch "arm64")
        elseif ("${CMAKE_OSX_ARCHITECTURES}" STREQUAL "x86_64")
            set(go_arch "amd64")
        else()
            message(FATAL_ERROR "Unknown Mac OS arch")
        endif()

        set(ENV_MANIP ${CMAKE_COMMAND} -E env
            "CGO_ENABLED=1"
            "MACOSX_DEPLOYMENT_TARGET=${CMAKE_OSX_DEPLOYMENT_TARGET}"
            "GOARCH=${go_arch}"
        )

        list(APPEND ADDITIONAL_FLAGS -gccgoflags
            '-mmacosx-version-min=${CMAKE_OSX_DEPLOYMENT_TARGET}'
            # Due to a latest Mac os linker, we need to disable debug symbol generation for the time being (-w).
            -ldflags '-w -extldflags -mmacosx-version-min=${CMAKE_OSX_DEPLOYMENT_TARGET}'
        )
    endif()

    add_custom_command (OUTPUT ${lib}
            DEPENDS ${CGO_GO_SOURCES}
            COMMAND ${ENV_MANIP} "${go_bin}" build -a -tags="${CGO_TAGS}" ${ADDITIONAL_FLAGS} -buildmode=c-shared  -o "${lib}"
            ${CGO_GO_EXPORTS}
            ${GO_LIB_BUILD_EXTRA}
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/cmd/lib
    )

    add_custom_target("${target}-go-lib-build"
            DEPENDS ${lib}
    )

    add_library(${target} SHARED IMPORTED GLOBAL)


    if (WIN32)
        win32_gen_implib(${target} "${name}" "${CMAKE_CURRENT_BINARY_DIR}" "${target}-go-lib-build" ${lib} ${win32_imp_lib})
	else()
        add_dependencies(${target} "${target}-go-lib-build")
    endif()

    set_target_properties(${target}
            PROPERTIES
            IMPORTED_LOCATION ${lib}
            IMPORTED_RUNTIME ${lib}
            IMPORTED_IMPLIB "${win32_imp_lib}"
            INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_BINARY_DIR};${CMAKE_CURRENT_SOURCE_DIR}/cmd/lib/cgo_headers"
    )
endfunction()

build_cgo_lib(
    TARGET etcore
    NAME proton-mail-export
    GO_SOURCES ${go_files}
    GO_EXPORTS export.go export_log.go export_mail.go export_globals.go
)

build_cgo_lib(
        TARGET etgpa
        NAME etgpa
        TAGS gpa_server
        GO_SOURCES ${go_files}
        GO_EXPORTS export_gpa.go
)

configure_file(${CMAKE_CURRENT_SOURCE_DIR}/cmake/constants.go.in
    ${CMAKE_CURRENT_SOURCE_DIR}/internal/constants.go
    @ONLY
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
        -package apiclient github.com/ProtonMail/export-tool/internal/apiclient Builder,Client,RetryStrategy > "${CMAKE_BINARY_DIR}/mocks_test.go"
    COMMAND ${CMAKE_COMMAND} -E copy_if_different "${CMAKE_BINARY_DIR}/mocks_test.go" "${CMAKE_CURRENT_SOURCE_DIR}/internal/apiclient/mocks.go"
    COMMAND mockgen
        -self_package "github.com/ProtonMail/export-tool/internal/mail"
        -package mail github.com/ProtonMail/export-tool/internal/mail StageErrorReporter,MetadataFileChecker,Reporter > "${CMAKE_BINARY_DIR}/mocks_mail_test.go"
    COMMAND ${CMAKE_COMMAND} -E copy_if_different "${CMAKE_BINARY_DIR}/mocks_mail_test.go" "${CMAKE_CURRENT_SOURCE_DIR}/internal/mail/mocks_test.go"
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)
