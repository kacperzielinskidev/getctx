package main

import (
	"flag"
	"fmt"
	"getctx/internal/core"
	"getctx/internal/logger"
	"os"
	"runtime/pprof"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	outputFilename := flag.String("o", "context.txt", "The name of the output file.")
	debug := flag.Bool("debug", false, "Enable debug level logging.")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	logFile, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logLevel := logger.LevelInfo
	if *debug {
		logLevel = logger.LevelDebug
	}

	log := logger.New(logFile, logLevel)
	log.Info("main", "Logger initialized successfully.")
	if *debug {
		log.Debug("main", "Debug logging is enabled.")
	}

	startPath := "."
	if len(flag.Args()) > 0 {
		startPath = flag.Args()[0]
	}

	app := core.NewApp(log, *outputFilename, startPath)
	if err := app.Run(); err != nil {
		log.Error("main.app.Run", err)
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}
