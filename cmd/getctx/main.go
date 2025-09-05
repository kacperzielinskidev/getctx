package main

import (
	"flag"
	"log"

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

	p := tea.NewProgram(app.NewModel(startPath))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("An error occurred while running the program: %v", err)
	}

	if m, ok := finalModel.(*app.Model); ok {
		if err := app.HandleContextBuilder(m, *outputFilename); err != nil {
			log.Fatalf("A critical error occurred while creating the context file: %v", err)
		}
	}
}
