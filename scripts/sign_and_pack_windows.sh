#!/usr/bin/env bash

set -eo pipefail

main(){
    SRC_DIR=$1
    TARGET_ZIP=$2

    if [ -z "${SRC_DIR}" ]; then
        echo "ERROR: missing first argument: SRC_DIR"
        exit 1
    fi;

    if [ -z "${TARGET_ZIP}" ]; then
        echo "ERROR: missing second argument: TARGET_ZIP"
        exit 1
    fi;

    AIpath="C:\\Program Files (x86)\\Caphyon\\Advanced Installer 14.2.1\\bin\\x86"
    AIsign="${AIpath}\\digisign.exe"
    PRODUCTNAME="Proton Mail Export CLI"
    TSA="http://timestamp.sectigo.com"

    pushd "$SRC_DIR"

    find . -type f | while read -r i; do
        echo "Signing ${i}..."
        "$AIsign" sign //a \
            //t "$TSA" \
            //d "$PRODUCTNAME" \
            //fd "SHA256" \
            "$(basename "$i")"
    done

    mkdir -p "$(dirname "$TARGET_ZIP")"
    7z a -tzip "$TARGET_ZIP" ./*

    popd
}

main "$@"
