package main

import (
	"fmt"
	"log"
	"os"

	"myapp/config"
	"myapp/logger"
	"myapp/monitor"
	"myapp/filter"
	"myapp/upload"
	"myapp/track"
	"myapp/clean"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp [monitor|filter|upload|track|clean]")
		os.Exit(1)
	}

	// Load configuration and initialize logger.
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	logger.Init(cfg.LogLevel)

	// Dispatch to the requested component.
	command := os.Args[1]
	var runErr error
	switch command {
	case "monitor":
		runErr = monitor.Run(cfg)
	case "filter":
		runErr = filter.Run(cfg)
	case "upload":
		runErr = upload.Run(cfg)
	case "track":
		runErr = track.Run(cfg)
	case "clean":
		runErr = clean.Run(cfg)
	default:
		fmt.Println("Unknown command. Usage: myapp [monitor|filter|upload|track|clean]")
		os.Exit(1)
	}
	if runErr != nil {
		logger.Error.Println("Execution error:", runErr)
		os.Exit(1)
	}
}
