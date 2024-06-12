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
    assert_output --partial "Invalid Zencode prefix 1: 'Gibberish'"
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
		assert_output --partial "Copyright (C) 2024 undefined"
		assert_output --partial "slangroom-exec"
		run ./slangroom-exec --version
    assert_output --partial "License AGPL-3.0-or-later: GNU AGPL version 3 <https://www.gnu.org/licenses/agpl-3.0.html>"
		assert_output --partial "Copyright (C) 2024 undefined"
		assert_output --partial "slangroom-exec"
}
