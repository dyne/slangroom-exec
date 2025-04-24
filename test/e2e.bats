# SPDX-FileCopyrightText: 2024-2025 Dyne.org foundation
#
# SPDX-License-Identifier: AGPL-3.0-or-later

setup() {
    load 'test_helper/common-setup'

    _common_setup
}

@test "slangroom-exec exists and is executable" {
    assert_file_exists ./slangroom-exec
    assert_file_executable ./slangroom-exec
    assert_size_not_zero ./slangroom-exec
}

@test "should execute simple zencode" {
    load_fixture "simple_zencode"
    run_slangroom_exec
    assert_output '{"output":["Welcome_to_slangroom-exec_🥳"]}'
    assert_success
}

@test "should execute simple slangroom" {
    load_fixture "simple_slangroom"
    run_slangroom_exec
    assert_output --partial 'timestamp'

    # check that ts is a number
    ts=$(echo $output | jq '.timestamp')
    assert_regex "$ts" '^[0-9]+$'
    assert_success
}

@test "should read data correctly" {
    export FILES_DIR="."
    load_fixture "read_file"
    run_slangroom_exec
    assert_output --partial "Do you know who greets you? 🥒"
    assert_success
}

@test "should fail on empty or broken contract" {
    load_fixture "broken_conf"
    run_slangroom_exec
    assert_output "Malformed input: Slangroom contract is empty"
    assert_failure 1
}

@test "should fail on broken slangroom" {
    load_fixture "broken_slangroom"
    run_slangroom_exec
    assert_output --partial "Gibberish may be given or then"
    assert_failure 1
}

@test "should fail on empty contract" {
    load_fixture "empty"
    run_slangroom_exec
    assert_output "Malformed input: Slangroom contract is empty"
    assert_failure 1
}

@test "should show the version and header" {
    run ./slangroom-exec -v
    assert_output --partial "License AGPL-3.0-or-later: GNU AGPL version 3 <https://www.gnu.org/licenses/agpl-3.0.html>"
    assert_output --partial "Copyright (C) 2024-2025 Dyne.org foundation"
    assert_output --partial "slangroom-exec"
    run ./slangroom-exec --version
    assert_output --partial "License AGPL-3.0-or-later: GNU AGPL version 3 <https://www.gnu.org/licenses/agpl-3.0.html>"
    assert_output --partial "Copyright (C) 2024-2025 Dyne.org foundation"
    assert_output --partial "slangroom-exec"
}

@test "simple chain execution" {
    encoded=$(src/slexfe -s test/fixtures/chain.yaml -d test/fixtures/chain.data.json)
    printf -v slang_input '%s' "$encoded"
    run_slangroom_exec_chain
    assert_output '{"hello":"hello from data!","output":["Hello,_world!"]}'
    assert_success
}

@test "execute chain that read contracts and keys from file" {
    encoded=$(src/slexfe -s test/fixtures/chain_sl_from_file.yaml -d test/fixtures/chain_sl_from_file.data.json)
    printf -v slang_input '%s' "$encoded"
    run_slangroom_exec_chain
    assert_output '{"hello_from_data":"hello from data!","hello_from_keys":"hello from keys!","output":["Hello,_world!"]}'
    assert_success
}
