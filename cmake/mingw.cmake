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

set(MINGW_ZIP_IS_NEW FALSE)

if (NOT EXISTS "${DOWNLOAD_ZIP_FILE}")
        message(STATUS "File not found, downloading mingw...")
        file(DOWNLOAD "${ZIP_URL}" "${DOWNLOAD_ZIP_FILE}" EXPECTED_HASH "SHA256=${EXPECTED_MINGW_ZIP_HASH}" SHOW_PROGRESS)
        set(MINGW_ZIP_IS_NEW TRUE)
endif()

if (EXISTS "${DOWNLOAD_ZIP_FILE}")
	message(STATUS "Found mingw.zip in cache path, validating hash")
	file(SHA256 "${DOWNLOAD_ZIP_FILE}" filehash)
	if (EXPECTED_MINGW_ZIP_HASH STREQUAL filehash)
		message(STATUS "File hash matches!")
        else()
		message(STATUS "File hash no longer matches, downloading mingw")
		file(DOWNLOAD "${ZIP_URL}" "${DOWNLOAD_ZIP_FILE}" EXPECTED_HASH "SHA256=${EXPECTED_MINGW_ZIP_HASH}" SHOW_PROGRESS)
		set(MINGW_ZIP_IS_NEW TRUE)
	endif()
endif()

if (EXISTS "${MINGW_PATH}" AND NOT MINGW_ZIP_IS_NEW)
	message(STATUS "Not extracting mingw, folder exists ${MINGW_PATH}")
else()
    message(STATUS "Extracting mingw archive")
    file(ARCHIVE_EXTRACT INPUT "${DOWNLOAD_ZIP_FILE}" DESTINATION "${MINGW_CACHE_PATH}")
endif()

message(STATUS "MINGW_PATH=${MINGW_PATH}")

find_program(gendef-bin gendef PATHS "${MINGW_PATH}")
if (NOT gendef-bin)
	message(FATAL_ERROR "could not find gendef binary")
else()
	message(STATUS "gendef-bin=${gendef-bin}")
endif()

find_program(dlltool-bin dlltool PATHS "${MINGW_PATH}")
if (NOT dlltool-bin)
	message(FATAL_ERROR "could not find dlltool, please install with 'pacman -S binutils'")
else()
	message(STATUS "dlltool-bin=${dlltool-bin}")
endif()

function(win32_gen_implib target name bin_dir go_build_target shared_lib lib_file)
	set(def_file "${bin_dir}/${name}.def")

	add_custom_target("${target}-gendef"
		DEPENDS ${shared_lib}
		COMMAND ${gendef-bin} "${shared_lib}"
		BYPRODUCTS "${def_file}"
		COMMAND_EXPAND_LISTS
		COMMENT "Generating ${def_file}"
	)

	add_custom_target("${target}-implib"
		DEPENDS ${def_file}
		COMMAND ${dlltool-bin} -d ${def_file} -l ${lib_file}
		BYPRODUCTS ${lib_file}
		COMMENT "Generating ${lib_file}"
	)

	add_dependencies("${target}-gendef" ${go_build_target})
	add_dependencies("${target}-implib" "${target}-gendef")
	add_dependencies("${target}" "${target}-implib")
endfunction()
