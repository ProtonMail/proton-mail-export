#!/usr/bin/env bash

set -eo pipefail

main(){
    TARGET_ZIP=$1
    APP_NAME=$2
    META_FOLDER=$3
    BIN_DIR_ARM=$4
    BIN_DIR_x86=$5


    PACKDIR="$PWD/tmp_pkg"
    APP_PATH="$PACKDIR/${APP_NAME}"
    BIN_DIR="$APP_PATH/Contents/MacOS"

    setup_tmp_dir
    merge_binaries
    add_meta_files
    target

    rm -rf "$PACKDIR"
}

setup_tmp_dir(){
    rm -rf "$PACKDIR"
    mkdir -p "$BIN_DIR"
}

merge_binaries(){
    find "${BIN_DIR_ARM}" -type f | while read -r i; do
        f="$(basename "$i")"
        echo -e "\nMerging binary $i..."
        lipo -create \
            -output "$BIN_DIR/$f" \
            "$BIN_DIR_ARM/$f" \
            "$BIN_DIR_x86/$f"
    done
}

add_meta_files(){
    echo -e "\nAdding meta files..."
    cp "$META_FOLDER/Info.plist" "$APP_PATH/Contents/Info.plist"

    resources="$APP_PATH/Contents/Resources"
    mkdir -p "$resources"
    cp "$META_FOLDER/icon.icns" "$resources/icon.icns"

    cp "$META_FOLDER/launcher.sh" "$BIN_DIR/launcher.sh"
    chmod +x "$BIN_DIR/launcher.sh"

    run_script="$(basename "${APP_NAME}" .app).sh"
    cp "$META_FOLDER/$run_script" "$BIN_DIR/$run_script"
    chmod +x "$BIN_DIR/$run_script"

    ls -lR "${APP_PATH}"
}

target(){
    mkdir -p "$(dirname "$TARGET_ZIP")"
    ditto -ck --keepParent "$APP_PATH" "$TARGET_ZIP"
}


main "$@"
