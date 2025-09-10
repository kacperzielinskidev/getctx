// Plik: cmd/getctx/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"getctx/internal/app"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	application, err := app.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialization error: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}

func setupLogging() *os.File {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println("--- Application Start ---")
	return f
}
