#!/usr/bin/env bash

set -eo pipefail

main(){
    SRC_DIR=$1
    TARGET_ZIP=$2
    PROGRAM_NAME=$3
    BUNDLE_NAME=$4

    if [ -z "${SRC_DIR}" ]; then
        echo "ERROR: missing first argument: SRC_DIR"
	exit 1
    fi;

    if [ -z "${TARGET_ZIP}" ]; then
        echo "ERROR: missing second argument: TARGET_ZIP"
	exit 1
    fi;

    if [ -z "${PROGRAM_NAME}" ]; then
        echo "ERROR: missing second argument: PROGRAM_NAME"
	exit 1
    fi;

    if [ -z "${BUNDLE_NAME}" ]; then
        echo "ERROR: missing second argument: BUNDLE_NAME"
	exit 1
    fi;


    tmproot="$(realpath "./tmp_pkg/")"
    APP_PATH="$tmproot/${BUNDLE_NAME}"
    VERSION="$(get_version)"

    FILE_TO_NOTARIZE="$(realpath "$SRC_DIR/../to_notarize.zip")"

    if [ -z "${APPLEUID}" ] || [ -z "${APPLEPASSW}" ]; then
        echo "ERROR: missing apple UID or password"
        exit 1
    fi

    make_app
    sign
    notarize
    staple
    target

    rm -rf "$tmproot"
}

get_version(){
    major="$(grep "ET_VERSION_MAJOR " < cmake/config.cmake | cut -d")" -f1 | cut -d" " -f2)"
    minor="$(grep "ET_VERSION_MINOR " < cmake/config.cmake | cut -d")" -f1 | cut -d" " -f2)"
    patch="$(grep "ET_VERSION_PATCH " < cmake/config.cmake | cut -d")" -f1 | cut -d" " -f2)"
    echo -n "$major.$minor.$patch"
}

make_app(){
    rm -rf "$tmproot"

    bindir="$APP_PATH/Contents/MacOS"
    mkdir -p "$bindir"

    # EXES
    cp "$SRC_DIR"/* "$bindir/"
    launch_script > "$bindir/launcher.sh"
    chmod +x "$bindir/launcher.sh"

    # Info.plist
    info_plist > "$APP_PATH/Contents/Info.plist"
}

launch_script(){
    cat << EOF
#!/bin/sh
open -a "Terminal" "\$(dirname \$0)/${PROGRAM_NAME}"
EOF
}

info_plist(){
    APP_NAME="$(basename "${BUNDLE_NAME}" .app)"
    IDENTIFIER="ch.protonmail.export-tool"
    EXECUTABLE="launcher.sh"
    YEAR="$(date +%Y)"
    COPYRIGHT="Copyright Proton AG $YEAR"
    VENDOR="Proton AG, © $YEAR"

    cat << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>CFBundleName</key>
<string>${APP_NAME}</string>
<key>CFBundleIdentifier</key>
<string>${IDENTIFIER}</string>
<key>CFBundleVersion</key>
<string>${VERSION}</string>
<key>CFBundleExecutable</key>
<string>${EXECUTABLE}</string>
<key>CFBundleShortVersionString</key>
<string>${VERSION}</string>
<key>CFBundleAllowMixedLocalizations</key>
<string>true</string>
<key>CFBundleDevelopmentRegion</key>
<string>English</string>
<key>CFBundlePackageType</key>
<string>APPL</string>
<key>CFBundleInfoDictionaryVersion</key>
<string>6.0</string>
<key>NSHumanReadableCopyright</key>
<string>${COPYRIGHT}</string>
<key>CFBundleGetInfoString</key>
<string>${VENDOR}</string>
<key>CFBundleDisplayName</key>
<string>${APP_NAME}</string>
<key>NSHighResolutionCapable</key>
<true/>
<key>NSPrincipalClass</key>
<string>NSApplication</string>
</dict>
</plist>
EOF
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
