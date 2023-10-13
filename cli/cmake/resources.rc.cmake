IDI_ICON1 ICON DISCARDABLE "@ET_REPO_ROOT@/res/icon.ico"

1 VERSIONINFO
FILEVERSION     @ET_VERSION_STR_COMMA@,0
PRODUCTVERSION  @ET_VERSION_STR_COMMA@,0
BEGIN
    BLOCK "StringFileInfo"
    BEGIN
        BLOCK "040904b0"
        BEGIN
        VALUE "Comments", "Proton Export CLI is an application that allows you to export your Proton data."
            VALUE "CompanyName", "@ET_VENDOR@"
            VALUE "FileDescription", "@ET_CLI_FULL_NAME@"
            VALUE "FileVersion", "@ET_VERSION_STR_COMMA@,0"
            VALUE "InternalName", "@ET_CLI_NAME@.exe"
            VALUE "LegalCopyright", "(C) @ET_BUILD_YEAR@ @ET_VENDOR@"
            VALUE "OriginalFilename", "@ET_CLI_NAME@.exe"
            VALUE "ProductName", "@ET_CLI_FULL_NAME@ for Windows"
            VALUE "ProductVersion", "@ET_VERSION_STR@"
        END
    END
    BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x0409, 0x04B0
    END
END
