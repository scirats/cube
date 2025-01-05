package main

import (
	"bufio"
	"strings"
	"bytes"
	"os"
	"os/exec"
)

func preExecCmd(command []string, buffer *string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)

	if buffer != nil {
		cmd.Stdin = bytes.NewReader([]byte(*buffer))
	} else {
		cmd.Stdin = os.Stdin
	}

	cmd.Stderr = os.Stderr

	return cmd
}

func postExecCmd(capture bool, cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer

	if capture {
		cmd.Stdout = &stdout
		err := cmd.Run()
		return strings.TrimSpace(stdout.String()), err
	} else {
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		return "", err
	}
}

func ExecCmdWithBuffer(command []string, buffer string, capture bool) (string, error) {
	cmd := preExecCmd(command, &buffer)
	return postExecCmd(capture, cmd)
}

func ExecCmd(command []string, capture bool) (string, error) {
	cmd := preExecCmd(command, nil)
	return postExecCmd(capture, cmd)
}

func GetResults(out string) ([]string, error) {
	reader := strings.NewReader(string(out))
	scanner := bufio.NewScanner(reader)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

