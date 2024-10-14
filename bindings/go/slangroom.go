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

func SlangroomExec(conf string, contract string, data string, keys string, extra string, context string) (SlangResult, error) {

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

	b64conf := b64.StdEncoding.EncodeToString([]byte(conf))
	fmt.Fprintln(stdin, b64conf)

	b64contract := b64.StdEncoding.EncodeToString([]byte(contract))
	fmt.Fprintln(stdin, b64contract)

	b64keys := b64.StdEncoding.EncodeToString([]byte(keys))
	fmt.Fprintln(stdin, b64keys)

	b64data := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Fprintln(stdin, b64data)

	b64extra := b64.StdEncoding.EncodeToString([]byte(extra))
	fmt.Fprintln(stdin, b64extra)

	b64context := b64.StdEncoding.EncodeToString([]byte(context))
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
