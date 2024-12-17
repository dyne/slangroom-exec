package slangroom

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type SlangResult struct {
	Output string
	Logs   string
}

type SlangroomInput struct {
	Conf     string
	Contract string
	Data     string
	Keys     string
	Extra    string
	Context  string
}

// Exec runs the slangroom-exec command with the provided input.
func Exec(input SlangroomInput) (SlangResult, error) {
	if _, err := exec.LookPath("slangroom-exec"); err != nil {
		return SlangResult{}, fmt.Errorf(
			"slangroom-exec command not found. Please install it by running:\n\n" +
				"wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) \\\n" +
				"-O ~/.local/bin/slangroom-exec && chmod +x ~/.local/bin/slangroom-exec",
		)
	}

	// Prepare command and pipes
	execCmd := exec.Command("slangroom-exec")
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	inputs := []string{
		input.Conf, input.Contract, input.Keys,
		input.Data, input.Extra, input.Context,
	}
	for _, data := range inputs {
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		fmt.Fprintln(stdin, encoded)
	}
	stdin.Close()

	// Start command execution
	if err := execCmd.Start(); err != nil {
		return SlangResult{}, fmt.Errorf("failed to start command: %w", err)
	}

	// Capture stdout and stderr concurrently
	var wg sync.WaitGroup
	wg.Add(2)

	stdoutOutput := make(chan string, 1)
	stderrOutput := make(chan string, 1)

	go func() {
		defer wg.Done()
		captureOutput(stdout, stdoutOutput)
	}()
	go func() {
		defer wg.Done()
		captureOutput(stderr, stderrOutput)
	}()

	waitErr := execCmd.Wait()

	wg.Wait()
	close(stdoutOutput)
	close(stderrOutput)

	// Retrieve outputs
	stdoutStr := <-stdoutOutput
	stderrStr := <-stderrOutput

	return SlangResult{Output: stdoutStr, Logs: stderrStr}, waitErr
}

// Introspect runs the slangroom-exec command in introspection mode.
func Introspect(contract string) (string, error) {

	if _, err := exec.LookPath("slangroom-exec"); err != nil {
		return "", fmt.Errorf(
			"slangroom-exec command not found. Please install it by running:\n\n" +
				"wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) \\\n" +
				"-O ~/.local/bin/slangroom-exec && chmod +x ~/.local/bin/slangroom-exec",
		)
	}

	execCmd := exec.Command("slangroom-exec", "-i")
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	fmt.Fprintln(stdin, contract)
	stdin.Close()

	if err := execCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	// Capture stdout
	stdoutCh := make(chan string, 1)
	go captureOutput(stdout, stdoutCh)

	// Wait for the command to complete
	waitErr := execCmd.Wait()
	close(stdoutCh)

	// Retrieve output
	stdoutStr := <-stdoutCh
	return stdoutStr, waitErr
}

func captureOutput(pipe io.ReadCloser, output chan<- string) {

	buf := new(strings.Builder)
	_, err := io.Copy(buf, pipe)
	if err != nil {
		log.Printf("Failed to capture output: %v", err)
		return
	}
	output <- buf.String()
}
