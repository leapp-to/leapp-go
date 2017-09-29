package executor

import (
	"strings"
	"testing"
)

func TestSuccessExec(t *testing.T) {
	var testData = []struct {
		Cmd      []string
		Stdin    string
		Stdout   string
		ExitCode int
	}{
		{[]string{"cat"}, "hello", "hello", 0},
		{[]string{"tr", "'[:lower:]'", "'[:upper:]'"}, "abc", "ABC", 0},
		{[]string{"tail", "-n", "1"}, "abc\nxyz", "xyz", 0},
	}

	// Test if stdout and exit code are correct
	for _, td := range testData {

		c := Command{
			CmdLine: td.Cmd,
			Stdin:   td.Stdin,
		}

		r := c.Execute()

		if r.Stdout != td.Stdout {
			t.Errorf("unexpected stdout: got=%s expected=%s\n", r.Stdout, td.Stdout)
		}

		if r.ExitCode != td.ExitCode {
			t.Errorf("unexpected exit code: got=%d expected=%d\n", r.ExitCode, td.ExitCode)
		}
	}
}

func TestFailExec(t *testing.T) {
	var testData = []struct {
		Cmd      []string
		Stdin    string
		Stderr   string
		ExitCode int
	}{
		{[]string{"cp"}, "invalid_input", "missing file operand", 1},
		{[]string{"mkdir"}, "invalid_input", "missing operand", 1},
	}

	// Test if stderr and exit code are correct
	for _, td := range testData {

		c := Command{
			CmdLine: td.Cmd,
			Stdin:   td.Stdin,
		}

		r := c.Execute()

		if !strings.Contains(r.Stderr, td.Stderr) {
			t.Errorf("unexpected stderr: got=%s expected to contain=%s\n", r.Stderr, td.Stderr)
		}

		if r.ExitCode != td.ExitCode {
			t.Errorf("unexpected exit code: got=%d expected=%d\n", r.ExitCode, td.ExitCode)
		}
	}
}
