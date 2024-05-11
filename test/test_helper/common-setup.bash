#!/usr/bin/env bash

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
    name=$1

    conf_file="test/fixtures/${name}.conf"
    zencode_file="test/fixtures/${name}.slang"
    keys_file="test/fixtures/${name}.keys.json"
    data_file="test/fixtures/${name}.data.json"
    extra_file="test/fixtures/${name}.extra.json"
    context_file="test/fixtures/${name}.context.json"

    conf=""
    zencode=""
    keys=""
    data=""
    extra=""
    context=""

    if [ -f "$conf_file" ]; then
        conf=$(cat "$conf_file")
    fi
    if [ -f "$zencode_file" ]; then
        zencode=$(base64 -w 0 "$zencode_file")
    fi
    if [ -f "$keys_file" ]; then
        keys=$(jq -c . "$keys_file" | base64 -w 0)
    fi
    if [ -f "$data_file" ]; then
        data=$(jq -c . "$data_file" | base64 -w 0)
    fi
    if [ -f "$extra_file" ]; then
        extra=$(jq -c . "$extra_file" | base64 -w 0)
    fi
    if [ -f "$context_file" ]; then
        context=$(jq -c . "$context_file" | base64 -w 0)
    fi

    printf -v slang_input '%s\n%s\n%s\n%s\n%s\n%s' "$conf" "$zencode" "$keys" "$data" "$extra" "$context"
}


run_slangroom_exec() {
    bats_require_minimum_version 1.5.0
    run -0 bats_pipe echo "$slang_input" \| ./slangroom-exec
}