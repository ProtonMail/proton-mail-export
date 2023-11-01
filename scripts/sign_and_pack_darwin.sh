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

    APP_PATH="$(eval realpath "$SRC_DIR/*.app")"
    FILE_TO_NOTARIZE="$(realpath "$SRC_DIR/../to_notarize.zip")"

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
    TEAMID="6UN54H93QT"
    SIGNACCOUNT="Developer ID Application: Proton Technologies AG ($TEAMID)"
    CODESIGN="codesign --force --verbose --deep --options runtime --timestamp --sign"

    echo -e "\nSigning ${APP_PATH}..."
    $CODESIGN "${SIGNACCOUNT}" "${APP_PATH}"

    echo -e "\nVerifying signature ${APP_PATH}..."
    codesign --verbose --verify "${APP_PATH}"
    codesign --verbose=4 --display "${APP_PATH}"
}


notarize() {
    ditto -ck --keepParent "${APP_PATH}" "${FILE_TO_NOTARIZE}"

    echo -e "\nSubmiting notarization ${FILE_TO_NOTARIZE}..."
    xcrun notarytool submit \
        "${FILE_TO_NOTARIZE}" \
        --apple-id "${APPLEUID}" --team-id "$TEAMID" --password "${APPLEPASSW}" \
        --wait | tee submit.log


    STATUS=$(grep status: submit.log | tail -1 | cut -d: -f2 | xargs)
    if [ "$STATUS" != "Accepted" ]; then
	    echo "ERROR during notarization submit"
	    exit 1
    fi;

    # Always check the log file, even if notarization succeeds, because it
    # might contain warnings that you can fix prior to your next submission.
    #
    # https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/customizing_the_notarization_workflow#3087732
    NOTARIZATIONID=$(grep id: submit.log | head -1 | cut -d: -f2 | xargs)
    xcrun notarytool log "${NOTARIZATIONID}" \
        --apple-id "${APPLEUID}" --team-id "$TEAMID" --password "${APPLEPASSW}" \
        notarization.json

    jq < notarization.json

    ISSUES=$(jq -r ".issues" < notarization.json)
    if [ "$ISSUES" != "null" ]; then
	    jq < notarization.json
	    echo "ERROR notarization found issues"
	    exit 1
    fi;

    echo -e "\nNotarization OK, no issues found\n"
}

staple(){
    echo -e "\nStapling $APP_PATH"
    xcrun stapler staple "$APP_PATH"

    echo -e "\nValidate staple $APP_PATH"
    stapler validate "$APP_PATH"
    spctl -a -vvv -t install "$APP_PATH"
}

target(){
    mkdir -p "$(dirname "$TARGET_ZIP")"

    ditto -ck --keepParent "${APP_PATH}" "${TARGET_ZIP}"
}


main "$@"
