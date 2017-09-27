package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func captureLog(f func()) string {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestStarts(t *testing.T) {
	up := make(chan struct{})

	go Main(up)
	select {
	case <-up:
		t.Log("server is up")
	case <-time.After(5 * time.Second):
		t.Errorf("server took too long to start")
	}
}

func TestShutdown(t *testing.T) {
	// This goroutine should start correctly
	go Main(nil)

	// This should return 1
	if Main(nil) != 1 {
		t.Errorf("server should not have started")
	}

	// Check listen error
	o := captureLog(func() {
		Main(nil)
	})

	if !strings.Contains(o, "address already in use") {
		t.Errorf("did not catch bind error")
	}

}
