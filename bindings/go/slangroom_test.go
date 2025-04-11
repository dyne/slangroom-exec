package slangroom

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteSimpleZencode(t *testing.T) {
	contract := `Given nothing
Then print the string 'Welcome to slangroom-exec 🥳'`
	res, success := Exec(SlangroomInput{Contract: contract})
	assert.JSONEq(t, `{"output":["Welcome_to_slangroom-exec_🥳"]}`, res.Output)
	assert.Nil(t, success, "Expected success but got failure")
}

func TestExecuteSimpleSlangroom(t *testing.T) {
	contract := `Rule unknown ignore
Given I fetch the local timestamp in seconds and output into 'timestamp'
Given I have a 'number' named 'timestamp'
Then print the 'timestamp'`
	res, success := Exec(SlangroomInput{Contract: contract})
	assert.Contains(t, res.Output, "timestamp")
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Output), &result); err == nil {
		ts, ok := result["timestamp"].(float64)
		assert.True(t, ok, "Expected timestamp to be present")
		assert.True(t, ts == float64(int(ts)), "Expected timestamp to be a number")
	} else {
		t.Errorf("Failed to unmarshal output: %v", err)
	}
	assert.Nil(t, success, "Expected success but got failure")
}

func TestFailOnBrokenSlangroom(t *testing.T) {
	contract := `Gibberish`
	res, success := Exec(SlangroomInput{Contract: contract})
	assert.Contains(t, res.Logs, "Zencode Invalid Statement @slangroom/ignored@1.45.0 Error: Maybe missing: \x1b[35mRule unknown ignore\x1b[0m")
	assert.NotNil(t, success, "Expected failure but got success")
}

func TestFailOnEmptyContract(t *testing.T) {
	contract := ``
	res, success := Exec(SlangroomInput{Contract: contract})
	assert.Equal(t, "Malformed input: Slangroom contract is empty\n", res.Logs)
	assert.NotNil(t, success, "Expected failure but got success")
}

func TestReadDataCorrectly(t *testing.T) {
	os.Setenv("FILES_DIR", ".")
	filePath := "test/test.txt"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Test file does not exist: %v", err)
	}
	contract := `Rule unknown ignore
Given I send path 'filename' and read verbatim file content and output into 'content'
Given I have a 'string' named 'filename'
Given I have a 'string' named 'content'
Then print data`
	data := `{
    "filename": "` + filePath + `"
}`
	res, success := Exec(SlangroomInput{Contract: contract, Data: data})
	assert.Contains(t, res.Output, "Do you know who greets you? 🥒")
	assert.Nil(t, success, "Expected success but got failure")
}

func TestFailOnEmptyOrBrokenContract(t *testing.T) {
	contract := ``
	conf := `error`
	res, success := Exec(SlangroomInput{Conf: conf, Contract: contract})
	assert.Equal(t, "Malformed input: Slangroom contract is empty\n", res.Logs)
	assert.NotNil(t, success, "Expected failure but got success")
}

func TestIntrsopection(t *testing.T) {
	contract := `Rule unknown ignore
Given I fetch the local timestamp in seconds and output into 'timestamp'
Given I have a 'number' named 'timestamp'
Then print the 'timestamp'`
	res, err := Introspect(contract)
	assert.Contains(t, res, "encoding")
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res), &result); err == nil {
		ts, ok := result["timestamp"].(map[string]interface{})
		assert.True(t, ok, "Expected timestamp to be present")
		_, ok = ts["encoding"]
		assert.True(t, ok, "Expected encoding to be present")
	} else {
		t.Errorf("Failed to unmarshal output: %v", err)
	}
	assert.Nil(t, err, "Expected success but got failure")
}
