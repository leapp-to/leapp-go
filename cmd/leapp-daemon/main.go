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
	os.Exit(Main(nil))
}

func Main(up chan<- struct{}) int {
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		return 0
	}

	// Parse options
	options := web.Options{
		ListenAddress: *flagListen,
		ReadTimeout:   time.Duration(defaultReadTimeout),
		WriteTimeout:  time.Duration(*flagTimeout),
	}

	// Set the appropriate env var for clients to connect
	os.Setenv("LEAPP_DAEMON_ADDR", options.ListenAddress)
	defer os.Unsetenv("LEAPP_DAEMON_ADDR")

	// Start HTTP server
	webHandler := web.New(&options)
	go webHandler.Run()

	// Handle shutdown under different conditions
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Received SIGTERM. Shutting down...")
	case up <- struct{}{}:
		log.Println("Up channel unblocked. Shutting down...")
	case err := <-webHandler.ErrorCh():
		log.Printf("Error starting service: %v\n", err)
		return 1
	}

	return 0
}
