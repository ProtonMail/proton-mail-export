#
# Configuration options.
#

set(ET_VERSION_MAJOR 0)
set(ET_VERSION_MINOR 1)
set(ET_VERSION_PATCH 0)

set(ET_VERSION_STR "${ET_VERSION_MAJOR}.${ET_VERSION_MINOR}.${ET_VERSION_PATCH}")

set(ET_DEFAULT_API_URL "https://mail-api.proton.me")

string(TIMESTAMP ET_BUILD_TIME "%Y-%m-%dT%H:%M:%SZ" UTC)

# Get git commit hash.
execute_process(
    COMMAND git rev-parse --short=10 HEAD
    OUTPUT_VARIABLE ET_REVISION
)
string(REPLACE "\n" "" ET_REVISION "${ET_REVISION}")