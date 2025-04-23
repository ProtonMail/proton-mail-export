#
# Configuration options.
#

set(ET_VERSION_MAJOR 1)
set(ET_VERSION_MINOR 0)
set(ET_VERSION_PATCH 5)
set(ET_VENDOR "Proton AG")
set(ET_BUILD_YEAR "2024")
set(ET_REPO_ROOT "${PROJECT_SOURCE_DIR}")

set(ET_VERSION_STR "${ET_VERSION_MAJOR}.${ET_VERSION_MINOR}.${ET_VERSION_PATCH}")
set(ET_VERSION_STR_COMMA "${ET_VERSION_MAJOR},${ET_VERSION_MINOR},${ET_VERSION_PATCH}")

set(ET_DEFAULT_API_URL "https://mail-api.proton.me")

string(TIMESTAMP ET_BUILD_TIME "%Y-%m-%dT%H:%M:%SZ" UTC)

# Get git commit hash.
execute_process(
    COMMAND git rev-parse --short=10 HEAD
    OUTPUT_VARIABLE ET_REVISION
)
string(REPLACE "\n" "" ET_REVISION "${ET_REVISION}")


if (WIN32)
    set(ET_APP_IDENTIFIER "windows-export")
elseif(APPLE)
    set(ET_APP_IDENTIFIER "macos-export")
elseif(UNIX)
    set(ET_APP_IDENTIFIER "linux-export")
else()
    message(FATAL_ERROR "Unknown platform")
endif()

set(ET_APP_IDENTIFIER "${ET_APP_IDENTIFIER}@${ET_VERSION_STR}")
set(ET_SENTRY_DNS "${SENTRY_DNS}")

if (ET_SENTRY_DNS)
    message(STATUS "Sentry Reporting is enabled for this build")
endif()

if (UNIX AND NOT APPLE)
    configure_file(
        "${PROJECT_SOURCE_DIR}/cmake/version.json.in"
        "${PROJECT_BINARY_DIR}/version.json"
        @ONLY
    )

    install(FILES "${PROJECT_BINARY_DIR}/version.json" DESTINATION "meta")
endif()
