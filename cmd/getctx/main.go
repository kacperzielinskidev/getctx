// Plik: cmd/getctx/main.go
package main

import (
	"fmt"
	"os"

	"getctx/internal/app"
	"getctx/internal/logger"
)

func main() {

	logFile, err := logger.InitGlobalLogger("debug.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CRITICAL: Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	application, err := app.NewApp()
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
