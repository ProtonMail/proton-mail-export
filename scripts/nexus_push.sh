#!/usr/bin/env bash

set -eo pipefail

config_from_env() {
    echo "${NEXUS_CONFIG}" | base64 -d > ~/.nexus-client.json
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
        wget \
            --user "$(jq -r .bridge.Login < ~/.nexus-client.json)" \
            --password "$(jq -r .bridge.Password < ~/.nexus-client.json)" \
            "https://nexus.protontech.ch/repository/bridge-devel-builds/util/$DEFAULT_OS/nexus-client"

        chmod +x ./nexus-client
    fi
}

upload() {
    ./nexus-client "$@"
}


config_from_env
download_client
upload "$@"

