package core

import (
	"fmt"
	"getctx/internal/build"
	"getctx/internal/config"
	"getctx/internal/fs"
	"getctx/internal/logger"
	"getctx/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	log            *logger.Logger
	fsys           fs.FileSystem
	config         *config.Config
	contextBuilder *build.ContextBuilder
	startPath      string
}

func NewApp(log *logger.Logger, outputFilename string, startPath string) *App {
	fsys := fs.NewOSFileSystem()
	cfg := config.NewConfig()
	builder := build.NewContextBuilder(log, fsys, outputFilename, cfg)

	log.Info("NewApp", map[string]string{
		"outputFile": outputFilename,
		"startPath":  startPath,
	})

	return &App{
		log:            log,
		fsys:           fsys,
		config:         cfg,
		contextBuilder: builder,
		startPath:      startPath,
	}
}

func (a *App) Run() error {
	model, err := tui.NewModel(a.startPath, a.config, a.fsys)
	if err != nil {
		err = fmt.Errorf("error initializing TUI model: %w", err)
		a.log.Error("App.Run.NewModel", err)
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		err = fmt.Errorf("an error occurred while running the TUI program: %w", err)
		a.log.Error("App.Run.p.Run", err)
		return err
	}

	if m, ok := finalModel.(*tui.Model); ok {
		if m.Aborted {
			fmt.Println("Operation cancelled.")
			a.log.Info("App.Run", "User aborted the operation.")
			return nil
		}

		err := a.contextBuilder.Build(m.GetSelectedPaths())
		if err != nil {
			err = fmt.Errorf("a critical error occurred while creating the context file: %w", err)
			a.log.Error("App.Run.BuildContext", err)
			return err
		}
	}

	return nil
}
