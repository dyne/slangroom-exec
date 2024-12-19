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
		log.Println("slangroom-exec command not found")
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
		log.Printf("Failed to create stdout pipe: %v\n", err)
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

	// Encode inputs and send them to stdin
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
	stdoutErr := make(chan error, 1)
	stderrErr := make(chan error, 1)
	// Capture stdout and stderr in separate goroutines
	go captureOutput(stdout, stdoutOutput, stdoutErr, &wg)
	go captureOutput(stderr, stderrOutput, stderrErr, &wg)

	// Wait for both goroutines to complete capturing output
	wg.Wait()

	// Retrieve outputs and errors
	stdoutStr := <-stdoutOutput
	stderrStr := <-stderrOutput
	if err := <-stdoutErr; err != nil {
		return SlangResult{Output: stdoutStr, Logs: stderrStr}, err
	}
	if err := <-stderrErr; err != nil {
		return SlangResult{Output: stdoutStr, Logs: stderrStr}, err
	}
	waitErr := execCmd.Wait()
	return SlangResult{Output: stdoutStr, Logs: stderrStr}, waitErr
}

// Introspect runs the slangroom-exec command in introspection mode.
func Introspect(contract string) (string, error) {
	// Check if the slangroom-exec command is available
	if _, err := exec.LookPath("slangroom-exec"); err != nil {
		log.Println("slangroom-exec command not found")
		return "", fmt.Errorf(
			"slangroom-exec command not found. Please install it by running:\n\n" +
				"wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) \\\n" +
				"-O ~/.local/bin/slangroom-exec && chmod +x ~/.local/bin/slangroom-exec",
		)
	}

	// Prepare the exec command with introspection flag
	execCmd := exec.Command("slangroom-exec", "-i")
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Send the contract to stdin
	fmt.Fprintln(stdin, contract)
	stdin.Close()

	// Start the exec command
	if err := execCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	// Capture stdout with a goroutine
	stdoutCh := make(chan string, 1)
	stdoutErr := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1) // Add only once here

	// Capture the output, simulating the server delay for debugging
	go func() {
		captureOutput(stdout, stdoutCh, stdoutErr, &wg)
	}()
	// Wait for the goroutine to complete
	wg.Wait()

	// Retrieve the captured stdout output
	stdoutStr := <-stdoutCh
	if err := <-stdoutErr; err != nil {
		return "", err
	}

	// Wait for the command to finish and capture any potential errors
	waitErr := execCmd.Wait()

	return stdoutStr, waitErr
}

// captureOutput captures the output from the pipe and sends it to the output channel.
func captureOutput(pipe io.ReadCloser, output chan<- string, errCh chan<- error, wg *sync.WaitGroup) {
	defer close(output)
	defer close(errCh)
	defer wg.Done() // Signal that this goroutine is done when it finishes, this is needed for handling delays
	buf := new(strings.Builder)
	_, err := io.Copy(buf, pipe)
	if err != nil {
		errCh <- fmt.Errorf("failed to capture output: %w", err)
		return
	}
	output <- buf.String()
	errCh <- nil
}
