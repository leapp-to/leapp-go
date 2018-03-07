package main

import (
	"fmt"

	"github.com/leapp-to/leapp-go/pkg/actors"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	_, err := actors.NewActorAPIService(true)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
