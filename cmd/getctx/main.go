package main

import (
	"flag"
	"fmt"
	"os"

	"getctx/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// TODO: Move from the main.go
	outputFilename := flag.String("o", "context.txt", "The name of the output file")
	flag.Parse()

	startPath := "."
	if flag.NArg() > 0 {
		startPath = flag.Arg(0)
	}

	config := app.NewConfig()
	model, err := app.NewModel(startPath, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing model: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while running the program: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(*app.Model); ok {
		if err := app.HandleContextBuilder(m, *outputFilename); err != nil {
			fmt.Fprintf(os.Stderr, "A critical error occurred while creating the context file: %v\n", err)
			os.Exit(1)
		}
	}
}
