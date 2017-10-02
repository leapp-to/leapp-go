// Package executor provides an abstration for executing actors using an external tool.
package executor

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
)

// actorRunner is an actor execution tool
// TODO: using a wrapper script for testing only, but runner.py in snactor repo should be refactored into a standalone tool so it can be used by leapp-daemon
const actorRunner = "/usr/local/bin/actor_runner"

// Result represents the outcome of a command execution.
type Result struct {
	Stderr   string `json:"stderr"`
	Stdout   string `json:"stdout"`
	ExitCode int    `json:"exit_code"`
}

// Command represents the command to be executed along its stdin.
type Command struct {
	cmdLine []string
	Stdin   string
}

// Execute executes a given command passing data to its stdin.
// It returns a Result struct mapping the info returned by the process executed.
func (c *Command) Execute() *Result {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(c.cmdLine[0], c.cmdLine[1:]...)

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

// New initializes a new Command that works with actorRunner
func New(actorName, stdin string) *Command {
	cl := append(strings.Split(actorRunner, " "), actorName)
	c := &Command{cmdLine: cl,
		Stdin: stdin,
	}

	return c
}
