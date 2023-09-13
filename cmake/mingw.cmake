#
# Setup MINGW for CGO on windows.
# We download the the standalone MINGW from https://winlibs.com/ and use that to compile
# the go library on windows without having to switch to a msys or cygwing shell.
#

if (NOT WIN32)
    return()
endif()


if (NOT MINGW_CACHE_PATH)
	set(MINGW_CACHE_PATH ${CMAKE_BINARY_DIR}/mingw)
	message(STATUS "MINGW_CACHE_PATH not provided, using ${MINGW_CACHE_PATH}")
endif()


set(ZIP_URL "https://github.com/brechtsanders/winlibs_mingw/releases/download/13.2.0mcf-16.0.6-11.0.1-ucrt-r2/winlibs-x86_64-mcf-seh-gcc-13.2.0-mingw-w64ucrt-11.0.1-r2.zip")
set(EXPECTED_MINGW_ZIP_HASH "247cce3632f4543275081428501d73034d2308b1b1d219fc4fee2d855d92a3ce")
set(DOWNLOAD_ZIP_FILE "${MINGW_CACHE_PATH}/mingw.zip")

set(MINGW_PATH "${MINGW_CACHE_PATH}/mingw64/bin")

if (EXISTS "${DOWNLOAD_ZIP_FILE}")
	message(STATUS "Found mingw.zip in cache path, validating hash")
	file(SHA256 "${DOWNLOAD_ZIP_FILE}" filehash)
	if (EXPECTED_MINGW_ZIP_HASH STREQUAL filehash)
		message(STATUS "File hash matches, no need to download")
        else()
		message(STATUS "File not found or hash no longer matches, downloading mingw")
		file(DOWNLOAD "${ZIP_URL}" "${DOWNLOAD_ZIP_FILE}" EXPECTED_HASH "SHA256=${EXPECTED_MINGW_ZIP_HASH}" SHOW_PROGRESS)
	endif()
endif()

if (EXISTS "${MINGW_PATH}")
	message(STATUS "Not extracting migw, folder exists ${MINGW_PATH}")
else()
    message(STATUS "Extracting mingw archive")
    file(ARCHIVE_EXTRACT INPUT "${DOWNLOAD_ZIP_FILE}" DESTINATION "${MINGW_CACHE_PATH}")
endif()

message(STATUS "MINGW_PATH=${MINGW_PATH}")
