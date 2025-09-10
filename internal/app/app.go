package app

import (
	"flag"
	"fmt"

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
		return fmt.Errorf("error initializing TUI model: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("an error occurred while running the program: %w", err)
	}
	if m, ok := finalModel.(*Model); ok {
		err := BuildContext(m, a.OutputFilename)
		if err != nil {
			return fmt.Errorf("a critical error occurred while creating the context file: %w", err)
		}
	}

	return nil
}
