#!/usr/bin/env bash

set -eo pipefail

# RELEASE_ENV is used as dotenv artifact for artifactlift pipelines.
RELEASE_ENV=./release.env
# DEPLOY_DIR is used in artifactlift pipelines as place where all artifacts should be stored.
DEPLOY_DIR=./build
# NEXUS_CLIENT assuming we are running from project root
nexus_client=./scripts/nexus_push.sh

magenta='\033[35;1m'
red='\033[31;1m'
nocolor='\033[0m'

error(){
    echo 1>&2
    FORMAT="${red}Error: $1${nocolor}\n"
    shift;
    # shellcheck disable=SC2059
    printf "$FORMAT" "$@" 1>&2
}

main() {
    parse_and_validate_input "$@"

    get_version

    # Clean file which contains list of deploy files.
    # Files are expected to be in DEPLOY_DIR so names are relative to this
    # folders.
    printf "" >tmp.files
    printf "" >tmp.metadata

    mkdir -p "$DEPLOY_DIR"

    if [ "$DEPLOY_OS" == "linux" ]; then
        prepare_metadata
    fi

    prepare_installer_files

    release_env

    ls "$DEPLOY_DIR"
    cat "$RELEASE_ENV"
}

parse_and_validate_input() {
    if [ "$#" -ne 2 ]; then
        error "deploy script takes exactly 2 arguments: et-tag generate-job-name"
        return 1
    fi

    ET_TAG=$1
    if ! [[ "$ET_TAG" =~ ^et-[0-9]{3}$ ]]; then
        error "arg 1: expected et tag, but have '$IDA_TAG'"
        return 1
    fi

    export ET_TAG
    export BUILD_NEXUS_PATH=":bridge:/export-tool/${ET_TAG}/deploy"


    generate_job="$(echo "$2" | cut -d: -f1)"
    case "$generate_job" in
    windows | linux | macos)
        DEPLOY_OS="$generate_job"
        ;;
    *)
        error "arg 2: unknown OS '$generate_job'"
        return 1
        ;;
    esac
    export DEPLOY_OS
}

get_version(){
    $nexus_client cp "$BUILD_NEXUS_PATH/version.json" ./version.json

    jq < ./version.json

    BUILD_VERSION=$(jq -r .version < version.json)
    if ! [[ "$BUILD_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+ ]]; then
        error "expected version from json file, but have '$BUILD_VERSION'"
        return 1
    fi
}

prepare_metadata() {
    version_json=imex_version.json
    cp ./version.json  "${DEPLOY_DIR}/${version_json}"
    printf "%s " "$version_json" >>tmp.metadata

    sign_file "${DEPLOY_DIR}/${version_json}"
    printf "%s.sig " "$version_json" >>tmp.metadata
}

prepare_installer_files() {
    case "${DEPLOY_OS}" in
    windows)
        download_installer_file "proton-mail-export-cli-windows_x86_64.zip"
        ;;
    macos)
        download_installer_file "proton-mail-export-cli-macos.zip"
        ;;
    linux)
        download_installer_file "proton-mail-export-cli-linux_x86_64.tar.gz"
        ;;
    esac
}

download_installer_file() {
    $nexus_client cp "${BUILD_NEXUS_PATH}/${1}" "${DEPLOY_DIR}/${1}"
    printf "%s " "${1// /\\ }" >>tmp.files
    printf "%s " "${1// /\\ }" >>tmp.metadata

    sign_file "${DEPLOY_DIR}/${1}"
    printf "%s.sig " "${1// /\\ }" >>tmp.files
    printf "%s.sig " "${1// /\\ }" >>tmp.metadata
}


release_env() {
    {
        echo "ARTIFACT_LIST=$(cat tmp.files)"
        echo "METADATA_LIST=$(cat tmp.metadata)"
        echo "RELEASE_VERSION=${BUILD_VERSION}"
        echo "VERSION=${BUILD_VERSION}"
    } >${RELEASE_ENV}
}

sign_file() {
    gpg --local-user E2C75D68E6234B07 --detach-sign "$1"
}



main "$@"
