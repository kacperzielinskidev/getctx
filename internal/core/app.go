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

// App is the central orchestrator. It holds all major components.
type App struct {
	log            *logger.Logger
	fsys           fs.FileSystem
	config         *config.Config
	contextBuilder *build.ContextBuilder
	startPath      string
}

// NewApp is the dependency injection root. It constructs all services.
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

// Run executes the main application flow.
func (a *App) Run() error {
	// Inject dependencies into the TUI model.
	model, err := tui.NewModel(a.startPath, a.config, a.fsys)
	if err != nil {
		err = fmt.Errorf("error initializing TUI model: %w", err)
		a.log.Error("App.Run.NewModel", err)
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the TUI.
	finalModel, err := p.Run()
	if err != nil {
		err = fmt.Errorf("an error occurred while running the TUI program: %w", err)
		a.log.Error("App.Run.p.Run", err)
		return err
	}

	// After the TUI exits, check the final model's state.
	if m, ok := finalModel.(*tui.Model); ok {
		// If the user aborted (e.g., Ctrl+C), do nothing.
		if m.Aborted {
			fmt.Println("Operation cancelled.")
			a.log.Info("App.Run", "User aborted the operation.")
			return nil
		}

		// Get selected paths from the model and pass them to the builder.
		err := a.contextBuilder.Build(m.GetSelectedPaths())
		if err != nil {
			err = fmt.Errorf("a critical error occurred while creating the context file: %w", err)
			a.log.Error("App.Run.BuildContext", err)
			return err
		}
	}

	return nil
}
