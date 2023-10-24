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

    FILE_TO_NOTARIZE="$(realpath "$SRC_DIR/../to_notarize.zip")"

    if [ -z "${APPLEUID}" ] || [ -z "${APPLEPASSW}" ]; then
        echo "ERROR: missing apple UID or password"
        exit 1
    fi

    sign
    notarize
    # No staple when there is no .app bundle
    target
}

sign(){
    TEAMID="6UN54H93QT"
    SIGNACCOUNT="Developer ID Application: Proton Technologies AG ($TEAMID)"
    CODESIGN="codesign --force --verbose --deep --options runtime --timestamp --sign"
    find "${SRC_DIR}" -type f | while read -r i; do
        echo "Signing ${i}..."
        $CODESIGN "${SIGNACCOUNT}" "${i}"

        echo "Verifying signature ${i}..."
	codesign --verbose --verify "${i}"
    done
}

notarize() {
    ditto -ck "${SRC_DIR}" "${FILE_TO_NOTARIZE}"

    echo "Submiting notarization ${FILE_TO_NOTARIZE}..."
    xcrun notarytool submit \
        "${FILE_TO_NOTARIZE}" \
        --apple-id "${APPLEUID}" --team-id "$TEAMID" --password "${APPLEPASSW}" \
        --wait | tee submit.log

    STATUS=$(cat submit.log | grep status: | tail -1 | cut -d: -f2 | xargs)
    if [ "$STATUS" != "Accepted" ]; then
	    echo "ERROR during notarization submit"
	    exit 1
    fi;

    # Always check the log file, even if notarization succeeds, because it
    # might contain warnings that you can fix prior to your next submission.
    #
    # https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/customizing_the_notarization_workflow#3087732
    NOTARIZATIONID=$(cat submit.log | grep id: | head -1 | cut -d: -f2 | xargs)
    xcrun notarytool log "${NOTARIZATIONID}" \
        --apple-id "${APPLEUID}" --team-id "$TEAMID" --password "${APPLEPASSW}" \
        notarization.json

    ISSUES=$(jq -r ".issues" < notarization.json)
    if [ "$ISSUES" != "null" ]; then
	    jq < notarization.json
	    echo "ERROR notarization found issues"
	    exit 1
    fi;

    echo "Notarization OK, no issues found"
}

target(){
    mkdir -p "$(dirname "$TARGET_ZIP")"
    cp "${FILE_TO_NOTARIZE}" "${TARGET_ZIP}"
}


main "$@"
