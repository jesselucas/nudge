package main

import (
	"log"
	"os"
)

func main() {
	// Read command arguments
	args := os.Args[1:]

	// Load from environment variable
	config := os.Getenv("NUDGEPATH")
	if config == "" {
		log.Fatal("NUDGEPATH environment variable isn't set")
	}
	n, err := CreateNudge(config)
	if err != nil {
		log.Fatal("Unabled to open nudge toml", err)
	}

	if err := n.Parse(args); err != nil {
		log.Fatal(err)
	}
}
