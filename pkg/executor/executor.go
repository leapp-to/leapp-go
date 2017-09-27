package executor

import (
	"bytes"
	"os/exec"
	"log"
)


func Execute(args []string) interface{} {
	var out bytes.Buffer

	cmd := exec.Command("python", "runner.py")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("stdout %q", out.String())

	// pass args to stdin
	// get stdout and stderr
	return nil
}
