#!/usr/bin/env bash

set -eo pipefail

config_from_env() {
    if ! [ -f "~/.nexus_client.json" ]; then
        echo "${NEXUS_CONFIG}" | base64 -d > ~/.nexus-client.json
    fi
}

download_client() {
    case "$OSTYPE" in
        linux*)
            DEFAULT_OS=linux
            ;;
        darwin*)
            DEFAULT_OS=darwin
            ;;
        msys*)
            DEFAULT_OS=windows
            ;;
    esac

    if ! [ -f ./nexus-client ]; then
        wget "https://nexus.protontech.ch/repository/bridge-devel-builds/util/$DEFAULT_OS/nexus-client"
        chmod +x ./nexus-client
    fi
}

upload() {
    echo "nexus-client $*"
    ./nexus-client "$@"
}


config_from_env
download_client
upload "$@"

