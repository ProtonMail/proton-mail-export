#!/usr/bin/env bash

check_program() {
	if ! command -v $1 &> /dev/null
	then
		echo "$1 is not installed"
		exit 1
	fi
}

check_program cmake
check_program ninja

BUILD_DIR=${ET_BUILD_DIR:-cmake-build-release}

cmake -S . -B ${BUILD_DIR} -G Ninja
cmake --build ${BUILD_DIR} --config Release

