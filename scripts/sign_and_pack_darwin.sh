#!/usr/bin/env bash

set -eo pipefail

main(){
    SRC_DIR=$1
    TARGET_ZIP=$2

    FILE_TO_NOTARIZE="$(readlink -e "$SRC_DIR/../to_notarize.zip")"
    NOTARIZED_FILE="$(readlink -e "$SRC_DIR/../notarized.zip")"

    if [ -z "${APPLEUID}" ] || [ -z "${APPLEPASSW}" ]; then
        echo "ERROR: missing apple UID or password"
        exit 1
    fi

    sign
    notarize
    staple
    target
}

sign(){
    SIGNACCOUNT="Developer ID Application: Proton Technologies AG (6UN54H93QT)"
    CODESIGN="codesign --force --verbose --deep --options runtime --timestamp --sign"
    echo "Signing ${SRC_DIR}..."
    $CODESIGN "${SIGNACCOUNT}" "${SRC_DIR}"

    codesign --verify "${SRC_DIR}"
}

notarize() {
    ditto -ck "${SRC_DIR}" "${FILE_TO_NOTARIZE}"

    echo "Submiting notarization ${FILE_TO_NOTARIZE}..."
    xcrun notarytool submit \
        "${FILE_TO_NOTARIZE}" \
        --apple-id "${APPLEUID}" --password "${APPLEPASSW}" \
        --wait

    # TODO check
}


staple(){
    # While you can notarize a ZIP archive, you can’t staple to it directly.
    # Instead, run stapler against each item that you added to the archive. Then
    # create a new ZIP file containing the stapled items for distribution. Although
    # tickets are created for standalone binaries, it’s not currently possible to
    # staple tickets to them.
    #
    # https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/customizing_the_notarization_workflow#3087720
    find "${SRC_DIR}" -type f | while read -r i; do
        xcrun stapler staple "$i"
    done
    ditto -ck "${SRC_DIR}" "${NOTARIZED_FILE}"
}


target(){
    mkdir -P "$(dirname "$TARGET_ZIP")"
    cp "${NOTARIZED_FILE}" "${TARGET_ZIP}"
}


main "$@"
