// Package executor provides an abstration for executing actors using an external tool.
package executor

import (
	"bytes"
	"io"
	"os"
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
	StdoutFile string
	StderrFile string
	cmdLine    []string
	Stdin      string
}

type CommandExecutionError string

func (cee CommandExecutionError) Error() string {
	return string(cee)
}

// Execute executes a given command passing data to its stdin.
// It returns a Result struct mapping the info returned by the process executed.
func (c *Command) Execute() (*Result, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(c.cmdLine[0], c.cmdLine[1:]...)

	cmd.Stdin = strings.NewReader(c.Stdin)
	if c.StdoutFile == "" {
		cmd.Stdout = &stdout
	} else {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return nil, CommandExecutionError("Failed to create stdout pipe: " + err.Error())
		}
		f, err := os.Create(c.StdoutFile)
		if err != nil {
			return nil, CommandExecutionError("Failed to create stdout file: " + err.Error())
		}
		go io.Copy(f, io.TeeReader(stdoutPipe, &stdout))
	}
	if c.StderrFile == "" {
		cmd.Stderr = &stderr
	} else {
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return nil, CommandExecutionError("Failed to create stderr pipe: " + err.Error())
		}
		f, err := os.Create(c.StderrFile)
		if err != nil {
			return nil, CommandExecutionError("Failed to create stderr file: " + err.Error())
		}
		go io.Copy(f, io.TeeReader(stderrPipe, &stderr))
	}

	code, err := runWithExitCode(cmd)
	if err != nil {
		err = CommandExecutionError("Executing process failed with: " + err.Error())
	}

	return &Result{
		Stderr:   stderr.String(),
		Stdout:   stdout.String(),
		ExitCode: code,
	}, err
}

// runWithExitCode simply runs a *exec.Cmd and return its exit code rather than an error
func runWithExitCode(cmd *exec.Cmd) (int, error) {
	var exitCode int

	err := cmd.Run()
	if err != nil {
		exitCode = 1 // Default to 1 if err is not ExitError

		if e, ok := err.(*exec.ExitError); ok {
			s := e.Sys().(syscall.WaitStatus)
			exitCode = s.ExitStatus()
		}
	} else {
		s := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = s.ExitStatus()
	}

	return exitCode, err
}

// NewProcess initializes a new Command
func NewProcess(process string, args ...string) *Command {
	return &Command{
		cmdLine: append([]string{process}, args...),
	}
}

// New initializes a new Command that works with actorRunner
func New(actorName, stdin string) *Command {
	c := NewProcess(actorRunner, actorName)
	c.Stdin = stdin
	return c
}
