function(copy_shared_libraries target dep)
    #if (WIN32)
    #    add_custom_command(target ${TARGET} POST_BUILD
    #            COMMAND ${CMAKE_COMMAND} -E copy_if_different $<target_RUNTIME_DLLS:$TARGET> $<TARGET_FILE_DIR:$TARGET>
    #            COMMAND_EXPAND_LISTS
    #    )
    #endif ()

    # Linux is handled with RPATH, does not require anything

    if (APPLE OR WIN32)
        add_custom_command(TARGET ${target} POST_BUILD
                COMMAND ${CMAKE_COMMAND} -E copy_if_different $<TARGET_FILE:${dep}> $<TARGET_FILE_DIR:${target}>
                COMMAND_EXPAND_LISTS
        )
    endif()

endfunction()
