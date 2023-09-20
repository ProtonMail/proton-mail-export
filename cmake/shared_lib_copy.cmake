function(copy_shared_libraries target dep)
    get_target_property(target_type ${dep} TYPE)

    if (NOT "${target_type}" STREQUAL "SHARED_LIBRARY")
        message(FATAL_ERROR "${dep} type is not a SHARED_LIBRARY, it is ${target_type}")
    endif()

    # Linux is handled with RPATH, does not require anything

    if (APPLE OR WIN32)
        add_custom_command(TARGET ${target} POST_BUILD
                COMMAND ${CMAKE_COMMAND} -E copy_if_different $<TARGET_FILE:${dep}> $<TARGET_FILE_DIR:${target}>
                COMMAND_EXPAND_LISTS
        )
    endif()

    if (APPLE)
        # Ensure the Mac os Loader can locate the shared libraries when installed.
        # Attempts at getting this to work correctly using just cmake RPATH variables did no lead to fruition.
        # Most useful info dump: https://stackoverflow.com/questions/33991581/install-name-tool-to-update-a-executable-to-search-for-dylib-in-mac-os-x
        # Rather than using CMake we can manually patch the targets ourselves.
        add_custom_command(TARGET ${target} POST_BUILD
                COMMAND install_name_tool -change "$<TARGET_FILE_NAME:${dep}>" "@loader_path/$<TARGET_FILE_NAME:${dep}>" "$<TARGET_FILE:${target}>"
                COMMENT "Patching mac os loader path for ${dep}"
                COMMAND_EXPAND_LISTS
        )
    endif()

endfunction()
