package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leapp-to/leapp-go/pkg/web"
)

var (
	defaultReadTimeout = 5 // seconds

	flagHelp    = flag.Bool("help", false, "show usage")
	flagListen  = flag.String("listen", "127.0.0.1:8000", "host:port to listen on.")
	flagTimeout = flag.Int64("timeout", 10, "time range in which daemon has to send a response to the client.")
)

func main() {
	os.Exit(Main())
}

func Main() int {
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Parse options
	options := web.Options{
		ListenAddress: *flagListen,
		ReadTimeout:   time.Duration(defaultReadTimeout),
		WriteTimeout:  time.Duration(*flagTimeout),
	}

	// Set the appropriate env var for clients to connect
	os.Setenv("LEAPP_DAEMON_ADDR", options.ListenAddress)

	// Start server
	webHandler := web.New(&options)
	go webHandler.Run()

	// Shutdown conditions
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Received SIGTERM. Shutting down...")
		return 0
	case err := <-webHandler.ErrorCh():
		log.Printf("Error starting service: %v\n", err)
		return 1
	}
	return 0
}
