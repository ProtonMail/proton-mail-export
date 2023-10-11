# https://stackoverflow.com/a/48877892
function run {
  $exe, $argsForExe = $Args
  # Workaround: Prevents 2> redirections applied to calls to this function
  #             from accidentally triggering a terminating error.
  #             See bug report at https://github.com/PowerShell/PowerShell/issues/4002
  $ErrorActionPreference = 'Continue'
  Write-Host "$exe $argsForExe"
  try { & $exe $argsForExe } catch { Throw } # catch is triggered ONLY if $exe can't be found, never for errors reported by $exe itself
  if ($LASTEXITCODE) { Throw "$exe failed with code $LASTEXITCODE" }
}

$ErrorActionPreference = 'Stop'

# before_script
New-Item -ItemType Directory -Path $Env:VCPKG_DEFAULT_BINARY_CACHE -Force ;
$Env:Path = "C:\\\\Program Files\\CMake\\bin;" + $Env:Path
$(Get-Command cmake).Source
run cmake --version

# !reference [.script-config-git, script]
run git config --global url.https://gitlab-ci-token:$Env:CI_JOB_TOKEN@$Env:CI_SERVER_HOST.insteadOf https://$Env:CI_SERVER_HOST
run git submodule update --init --recursive

# !reference [.script-config-cmake-win, script]
run cmake `
       -S . -B "$Env:CMAKE_BUILD_DIR" `
       -G "Visual Studio 17 2022" `
       -DCMAKE_GENERATOR_PLATFORM=x64 `
       -DCMAKE_CL_64=1 `
       -A "x64" `
       -DMINGW_CACHE_PATH="$Env:VCPKG_DEFAULT_BINARY_CACHE" `
       -DVCPKG_HOST_TRIPLET=x64-windows-static `
       -DVCPKG_TARGET_TRIPLET=x64-windows-static `
       -DCMAKE_INSTALL_PREFIX="$Env:CMAKE_INSTALL_DIR" `
       $Env:CMAKE_SENTRY_OPTION

# !reference [.script-test, script]
run cmake --build $Env:CMAKE_BUILD_DIR --config $Env:CMAKE_BUILD_CONFIG
run ctest --build-config $Env:CMAKE_BUILD_CONFIG --test-dir $Env:CMAKE_BUILD_DIR -V

# !reference [.script-install, script]
run cmake --install $Env:CMAKE_BUILD_DIR --config $Env:CMAKE_BUILD_CONFIG
