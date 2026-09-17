#!/usr/bin/env bash

# Usage: test_cpp.sh {linux|windows|darwin}

set -eo pipefail

: "${CMAKE_BUILD_DIR:=cmake-build-release}"
: "${CMAKE_BUILD_CONFIG:=Release}"

JUNIT_FILE="$PWD/lib/test-cpp-$1.xml"

cmake --build "$CMAKE_BUILD_DIR" --config "$CMAKE_BUILD_CONFIG" -t etcpp_test

ctest --test-dir "$CMAKE_BUILD_DIR" \
    --build-config "$CMAKE_BUILD_CONFIG" \
    --tests-regex '^etcpp-test$' \
    --repeat until-pass:2 \
    --timeout 900 \
    --output-on-failure \
    --output-junit "$JUNIT_FILE"
