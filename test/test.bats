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
    run ./slangroom-exec < test/fixtures/simple_zencode.zen
    assert_output '{"output":["Welcome_to_slangroom-exec_🥳"]}'
    assert_success
}

@test "should execute simple slangroom" {
    run ./slangroom-exec < test/fixtures/simple_slangroom.zen
    assert_output --partial timestamp

    # check that ts is a number
    ts=$(echo $output | jq '.timestamp')
    assert_regex "$ts" '^[0-9]+$'

    assert_success
}
