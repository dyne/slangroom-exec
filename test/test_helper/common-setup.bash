#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2024-2025 Dyne.org foundation
#
# SPDX-License-Identifier: AGPL-3.0-or-later

_common_setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    load 'test_helper/bats-file/load'

    # get the containing directory of this file
    # use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
    # as those will point to the bats executable's location or the preprocessed file respectively
    PROJECT_ROOT="$( cd "$( dirname "$BATS_TEST_FILENAME" )/.." >/dev/null 2>&1 && pwd )"
    # make executables in src/ visible to PATH
    PATH="$PROJECT_ROOT/src:$PATH"
}

load_fixture() {
    encoded=$(src/slexfe -F test/fixtures/$1)
    printf -v slang_input '%s' "$encoded"
}


run_slangroom_exec() {
    bats_require_minimum_version 1.5.0
    run bats_pipe printf "$slang_input" \| ./slangroom-exec
}

run_slangroom_exec_chain() {
    bats_require_minimum_version 1.5.0
    run bats_pipe printf "$slang_input" \| ./slangroom-exec -c
}
