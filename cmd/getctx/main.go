package main

import (
	"flag"
	"fmt"
	"os"

	"getctx/internal/app"
	"getctx/internal/logger"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug logging to debug.log")
	outputFilename := flag.String("o", "context.txt", "The name of the output file")
	flag.Parse()

	debugEnv := os.Getenv("GETCTX_DEBUG")
	enableLogging := *debug || debugEnv == "true" || debugEnv == "1"

	if enableLogging {
		logFile, err := logger.InitGlobalLogger("debug.log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "CRITICAL: Failed to initialize logger: %v\n", err)
		} else {
			defer logFile.Close()
		}
	}

	application, err := app.NewApp(*outputFilename, flag.Args())
	if err != nil {
		logger.Error("main.NewApp", err)
		fmt.Fprintf(os.Stderr, "Initialization error: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		logger.Error("main.Run", err)
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}

	logger.Info("main", "Application finished successfully.")

}
