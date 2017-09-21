package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/leapp-to/leapp-go/pkg/web"
)

func main() {
	go web.RunHTTPServer()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		fmt.Println("Goodbye !")
	}

}
