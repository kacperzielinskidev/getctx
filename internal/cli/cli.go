package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"

	"github.com/kacperzielinskidev/getctx/internal/build"
	"github.com/kacperzielinskidev/getctx/internal/config"
	"github.com/kacperzielinskidev/getctx/internal/core"
	"github.com/kacperzielinskidev/getctx/internal/fs"
	"github.com/kacperzielinskidev/getctx/internal/logger"
)

func Run() error {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	outputFilename := flag.String("o", "context.txt", "The name of the output file.")
	debug := flag.Bool("debug", false, "Enable debug level logging.")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	var logOutput io.Writer = io.Discard
	logLevel := logger.LevelInfo

	if *debug {
		debugFile, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("could not open log file: %w", err)
		}
		defer debugFile.Close()

		logOutput = debugFile
		logLevel = logger.LevelDebug
	}
	log := logger.New(logOutput, logLevel)
	log.Info("main", "Logger initialized successfully.")

	fsys := fs.NewOSFileSystem()
	appConfig := config.NewConfig()
	contextBuilder := build.NewContextBuilder(log, fsys, *outputFilename, appConfig)

	startPath := "."
	if len(flag.Args()) > 0 {
		startPath = flag.Args()[0]
	}

	app := core.NewApp(log, contextBuilder, appConfig, fsys, startPath)
	return app.Run()
}
