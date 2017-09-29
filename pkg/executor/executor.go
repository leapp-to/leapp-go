// Package executor provides an abstration for executing commands in the OS.
package executor

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
)

// Result represents the outcome of a command execution.
type Result struct {
	Stderr   string `json:"stderr"`
	Stdout   string `json:"stdout"`
	ExitCode int    `json:"exit_code"`
}

// Command represents the command to be executed along its stdin.
type Command struct {
	CmdLine []string
	Stdin   string
}

// Execute executes a given command passing data to its stdin.
// It returns a Result struct mapping the info returned by the process executed.
func (c *Command) Execute() *Result {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(c.CmdLine[0], c.CmdLine[1:]...)

	cmd.Stdin = strings.NewReader(c.Stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	code := runWithExitCode(cmd)

	return &Result{
		Stderr:   stderr.String(),
		Stdout:   stdout.String(),
		ExitCode: code,
	}
}

// runWithExitCode simply runs a *exec.Cmd and return its exit code rather than an error
func runWithExitCode(cmd *exec.Cmd) int {
	var exitCode int

	if err := cmd.Run(); err != nil {
		exitCode = 1 // Default to 1 if err is not ExitError

		if e, ok := err.(*exec.ExitError); ok {
			s := e.Sys().(syscall.WaitStatus)
			exitCode = s.ExitStatus()
		}
	} else {
		s := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = s.ExitStatus()
	}

	return exitCode
}
