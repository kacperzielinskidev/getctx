package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"

	"github.com/kacperzielinskidev/getctx/internal/logger"
)

type flagConfig struct {
	outputFilename string
	logOutput      io.Writer
	logLevel       logger.Level
	startPath      string
}

type cleanupFunc func()

var noOpCleanup = func() {}

// TODO: handle --version flag ( version should be set automatically durning release )
func setupAndParseFlags() (*flagConfig, cleanupFunc, error) {
	fs := flag.NewFlagSet("getctx", flag.ExitOnError)

	cpuprofile := fs.String("cpuprofile", "", "write cpu profile to file")
	outputFilename := fs.String("o", "context.txt", "The name of the output file.")
	debug := fs.Bool("debug", false, "Enable debug level logging.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, nil, fmt.Errorf("could not parse flags: %w", err)
	}

	config := &flagConfig{
		outputFilename: *outputFilename,
		logOutput:      io.Discard,
		logLevel:       logger.LevelInfo,
	}

	if fs.NArg() > 0 {
		config.startPath = fs.Arg(0)
	} else {
		config.startPath = "."
	}

	cleanup := noOpCleanup

	if *debug {
		debugFile, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, nil, fmt.Errorf("could not open log file: %w", err)
		}
		config.logOutput = debugFile
		config.logLevel = logger.LevelDebug

		cleanup = func() {
			debugFile.Close()
		}
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return nil, nil, fmt.Errorf("could not create CPU profile: %w", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			return nil, nil, fmt.Errorf("could not start CPU profile: %w", err)
		}

		existingCleanup := cleanup
		cleanup = func() {
			pprof.StopCPUProfile()
			f.Close()
			existingCleanup()
		}
	}

	return config, cleanup, nil

}
