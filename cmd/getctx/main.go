package main

import (
	"flag"
	"fmt"
	"getctx/internal/core"
	"getctx/internal/logger"
	"os"
)

func main() {
	// Initialize the logger first.
	logFile, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Create a logger instance instead of using a global one.
	log := logger.New(logFile, logger.LevelDebug)
	log.Info("main", "Logger initialized successfully.")

	// --- Flag Parsing ---
	outputFilename := flag.String("o", "context.txt", "The name of the output file.")
	flag.Parse()

	startPath := "."
	if len(flag.Args()) > 0 {
		startPath = flag.Args()[0]
	}

	// Create the core App, injecting the logger and other dependencies.
	app := core.NewApp(log, *outputFilename, startPath)

	// Run the application.
	if err := app.Run(); err != nil {
		log.Error("main.app.Run", err)
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}
