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

type SlangroomChainInput struct {
	Chain    string
	Data     string
}

type SlangroomExec struct {
	Cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Stdin  io.WriteCloser
}

func PrepareSlangroomExec(args ...string) (SlangroomExec, error) {
	if _, err := exec.LookPath("slangroom-exec"); err != nil {
		log.Println("slangroom-exec command not found")
		return SlangroomExec{}, fmt.Errorf(
			"slangroom-exec command not found. Please install it by running:\n\n" +
				"wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) \\\n" +
				"-O ~/.local/bin/slangroom-exec && chmod +x ~/.local/bin/slangroom-exec",
		)
	}
	// Prepare command and pipes
	execCmd := exec.Command("slangroom-exec", args...)
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to create stdout pipe: %v\n", err)
		return SlangroomExec{}, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return SlangroomExec{}, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return SlangroomExec{}, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	return SlangroomExec{ Cmd: execCmd, Stdout: stdout, Stderr: stderr, Stdin: stdin }, nil
}

// Exec runs the slangroom-exec command with the provided input.
func Exec(input SlangroomInput) (SlangResult, error) {
	slangroomExec, err := PrepareSlangroomExec()
	if err != nil {
		return SlangResult{}, fmt.Errorf("%w", err)
	}

	// Encode inputs and send them to stdin
	inputs := []string{
		input.Conf, input.Contract, input.Keys,
		input.Data, input.Extra, input.Context,
	}
	for _, data := range inputs {
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		fmt.Fprintln(slangroomExec.Stdin, encoded)
	}
	slangroomExec.Stdin.Close()

	// Start command execution
	if err := slangroomExec.Cmd.Start(); err != nil {
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
	go captureOutput(slangroomExec.Stdout, stdoutOutput, stdoutErr, &wg)
	go captureOutput(slangroomExec.Stderr, stderrOutput, stderrErr, &wg)

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
	waitErr := slangroomExec.Cmd.Wait()
	return SlangResult{Output: stdoutStr, Logs: stderrStr}, waitErr
}

// Introspect runs the slangroom-exec command in introspection mode.
func Introspect(contract string) (string, error) {
	slangroomExec, err := PrepareSlangroomExec("-i")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// Send the contract to stdin
	fmt.Fprintln(slangroomExec.Stdin, contract)
	slangroomExec.Stdin.Close()

	// Start the exec command
	if err := slangroomExec.Cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	// Capture stdout with a goroutine
	stdoutCh := make(chan string, 1)
	stdoutErr := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1) // Add only once here

	// Capture the output, simulating the server delay for debugging
	go func() {
		captureOutput(slangroomExec.Stdout, stdoutCh, stdoutErr, &wg)
	}()
	// Wait for the goroutine to complete
	wg.Wait()

	// Retrieve the captured stdout output
	stdoutStr := <-stdoutCh
	if err := <-stdoutErr; err != nil {
		return "", err
	}

	// Wait for the command to finish and capture any potential errors
	waitErr := slangroomExec.Cmd.Wait()

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

func ExecChain(input SlangroomChainInput) (SlangResult, error) {
	slangroomExec, err := PrepareSlangroomExec("-c")
	if err != nil {
		return SlangResult{}, fmt.Errorf("%w", err)
	}

	// Encode inputs and send them to stdin
	inputs := []string{
		input.Chain,
		input.Data,
	}
	fmt.Fprintln(slangroomExec.Stdin, "")
	for _, data := range inputs {
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		fmt.Fprintln(slangroomExec.Stdin, encoded)
	}
	slangroomExec.Stdin.Close()

	// Start command execution
	if err := slangroomExec.Cmd.Start(); err != nil {
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
	go captureOutput(slangroomExec.Stdout, stdoutOutput, stdoutErr, &wg)
	go captureOutput(slangroomExec.Stderr, stderrOutput, stderrErr, &wg)

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
	waitErr := slangroomExec.Cmd.Wait()
	return SlangResult{Output: stdoutStr, Logs: stderrStr}, waitErr
}
