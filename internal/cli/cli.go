package cli

import (
	"errors"
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
	fsys := fs.NewOSFileSystem()
	appConfig := config.NewConfig()

	contextBuilder := build.NewContextBuilder(log, fsys, appConfig)
	app := core.NewApp(log, contextBuilder, appConfig, fsys, cfg.startPath, cfg.outputFilename)

	result, err := app.Run()
	if err != nil {
		if errors.Is(err, core.ErrAbortedByUser) {
			return nil
		}
		return err
	}

	return presentResults(result, cfg.outputFilename)

}

func presentResults(result *build.BuildResult, outputFilename string) error {
	// This case will occur if the user has not selected any files.
	if result.FilesProcessed == 0 && len(result.PathsWithErr) == 0 {
		fmt.Println("ℹ️ No files selected. The output file was not created.")
		return nil
	}

	if len(result.PathsWithErr) > 0 {
		fmt.Println("⚠️ Some paths were skipped due to errors:")
		for _, warn := range result.PathsWithErr {
			fmt.Printf("   - %s\n", warn)
		}
		fmt.Println() // Additional empty line for readability
	}

	if result.FilesProcessed == 0 {
		message := "ℹ️ No text files found to include."
		if result.FilesSkipped > 0 {
			message += fmt.Sprintf(" %d file(s) were skipped (non-text or unreadable).", result.FilesSkipped)
		}
		message += " The output file was not created."
		fmt.Println(message)
		return nil
	}

	fmt.Printf("🚀 Processing finished. Found %d files to process.\n", result.FilesProcessed)
	fmt.Printf("✅ Done! All content has been combined into the file %s\n", outputFilename)

	return nil
}
