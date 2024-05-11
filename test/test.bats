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