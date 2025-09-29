package cli

import (
	"fmt"

	"github.com/kacperzielinskidev/getctx/internal/build"
	"github.com/kacperzielinskidev/getctx/internal/config"
	"github.com/kacperzielinskidev/getctx/internal/core"
	"github.com/kacperzielinskidev/getctx/internal/fs"
	"github.com/kacperzielinskidev/getctx/internal/logger"
)

func Run() error {
	cfg, cleanup, err := setupAndParseFlags()

	if err != nil {
		return fmt.Errorf("failed to setup flags: %w", err)
	}
	defer cleanup()

	log := logger.New(cfg.logOutput, cfg.logLevel)
	log.Info("main", "Logger initialized successfully.")

	fsys := fs.NewOSFileSystem()
	appConfig := config.NewConfig()
	contextBuilder := build.NewContextBuilder(log, fsys, cfg.outputFilename, appConfig)

	app := core.NewApp(log, contextBuilder, appConfig, fsys, cfg.startPath)
	return app.Run()
}
