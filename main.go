package main

import (
	"fmt"
	"os"

	"github.com/timakin/llmstxt-gen/internal/app"
)

// These variables will be set by GoReleaser at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Check if version flag is provided
	for _, arg := range os.Args {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("llmstxt-gen version %s, commit %s, built at %s\n", version, commit, date)
			return
		}
	}

	app.Run()
}
