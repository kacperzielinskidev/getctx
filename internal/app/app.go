package app

import (
	"flag"
	"fmt"
	"getctx/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	Config         *Config
	FileSystem     FileSystem
	OutputFilename string
	StartPath      string
}

func NewApp() (*App, error) {
	outputFilename := flag.String("o", "context.txt", "The name of the output file")
	flag.Parse()

	startPath := "."
	if flag.NArg() > 0 {
		startPath = flag.Arg(0)
	}

	logger.Info("NewApp", map[string]string{
		"outputFile": *outputFilename,
		"startPath":  startPath,
	})

	return &App{
		Config:         NewConfig(),
		FileSystem:     NewOSFileSystem(),
		OutputFilename: *outputFilename,
		StartPath:      startPath,
	}, nil
}

func (a *App) Run() error {
	model, err := NewModel(a.StartPath, a.Config, a.FileSystem)
	if err != nil {
		err = fmt.Errorf("error initializing TUI model: %w", err)
		logger.Error("App.Run.NewModel", err)
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		err = fmt.Errorf("an error occurred while running the TUI program: %w", err)
		logger.Error("App.Run.p.Run", err)
		return err
	}

	if m, ok := finalModel.(*Model); ok {
		err := BuildContext(m, a.OutputFilename)
		if err != nil {
			err = fmt.Errorf("a critical error occurred while creating the context file: %w", err)
			logger.Error("App.Run.BuildContext", err)
			return err
		}
	}

	return nil
}
