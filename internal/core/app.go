package core

import (
	"fmt"

	"github.com/kacperzielinskidev/getctx/internal/build"
	"github.com/kacperzielinskidev/getctx/internal/config"
	"github.com/kacperzielinskidev/getctx/internal/fs"
	"github.com/kacperzielinskidev/getctx/internal/logger"
	"github.com/kacperzielinskidev/getctx/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	log            *logger.Logger
	contextBuilder *build.ContextBuilder
	config         *config.Config
	fsys           fs.FileSystem
	startPath      string
}

func NewApp(
	log *logger.Logger,
	builder *build.ContextBuilder,
	cfg *config.Config,
	fsys fs.FileSystem,
	startPath string,
) *App {
	return &App{
		log:            log,
		contextBuilder: builder,
		config:         cfg,
		fsys:           fsys,
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
