package slangroom

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	b64 "encoding/base64"
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

func Exec(input SlangroomInput) (SlangResult, error) {

	if _, err := exec.LookPath("slangroom-exec"); err != nil {
		return SlangResult{}, fmt.Errorf(
			"slangroom-exec command not found. Please install it by running:\n\n" +
				"wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) \\\n" +
				"-O ~/.local/bin/slangroom-exec && chmod +x ~/.local/bin/slangroom-exec",
		)
	}
	execCmd := exec.Command("slangroom-exec")
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to create stderr pipe: %v", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %v", err)
	}

	b64conf := b64.StdEncoding.EncodeToString([]byte(input.Conf))
	fmt.Fprintln(stdin, b64conf)

	b64contract := b64.StdEncoding.EncodeToString([]byte(input.Contract))
	fmt.Fprintln(stdin, b64contract)

	b64keys := b64.StdEncoding.EncodeToString([]byte(input.Keys))
	fmt.Fprintln(stdin, b64keys)

	b64data := b64.StdEncoding.EncodeToString([]byte(input.Data))
	fmt.Fprintln(stdin, b64data)

	b64extra := b64.StdEncoding.EncodeToString([]byte(input.Extra))
	fmt.Fprintln(stdin, b64extra)

	b64context := b64.StdEncoding.EncodeToString([]byte(input.Context))
	fmt.Fprintln(stdin, b64context)

	stdin.Close()

	err = execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	stdoutOutput := make(chan string)
	stderrOutput := make(chan string)
	go captureOutput(stdout, stdoutOutput)
	go captureOutput(stderr, stderrOutput)

	err = execCmd.Wait()

	stdoutStr := <-stdoutOutput
	stderrStr := <-stderrOutput

	return SlangResult{Output: stdoutStr, Logs: stderrStr}, err
}
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
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %v", err)
	}

	fmt.Fprintln(stdin, contract)
	stdin.Close()

	err = execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	stdoutOutput := make(chan string)
	go captureOutput(stdout, stdoutOutput)

	err = execCmd.Wait()

	stdoutStr := <-stdoutOutput

	return stdoutStr, err

}
func captureOutput(pipe io.ReadCloser, output chan<- string) {
	defer close(output)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, pipe)
	if err != nil {
		log.Printf("Failed to capture output: %v", err)
		return
	}
	output <- buf.String()
}
