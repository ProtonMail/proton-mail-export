# Export Tool

New export tool according to [GODT-2804 RFC](https://confluence.protontech.ch/display/BRIDGE/%5BGODT-2804%5D+New+Export+Tool).


## Directories

* go-lib: Go Shared library implementation
* lib: C++ shared library over the exported C interface
* cli: CLI application


# Building

## Linux/Mac

```
cmake -S. -B $BUILD_DIR -G <Insert favorite Generator>
cmake --build $BUILD_DIR
```

## Windows (MSYS)

Sadly we need to use msys to build on windows due to Go's limitations.

Be sure to install the following packages:

```
 pacman -S mingw-w64-x86_64-cmake binutils pacman base-devel mingw-w64-x86_64-toolchain

```

Then configure CMake:

```
cmake -DVCPKG_TARGET_TRIPLET=x64-mingw-static  -S. -B $BUILD_DIR -G <[Ninja, "Ninja Multi-Config", Makefiles]  
```
