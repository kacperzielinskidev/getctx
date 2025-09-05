package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	outputFilename := flag.String("o", "context.txt", "The name of the output file")
	flag.Parse()

	startPath := "."
	if flag.NArg() > 0 {
		startPath = flag.Arg(0)
	}

	p := tea.NewProgram(newModel(startPath))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("An error occurred while running the program: %v", err)
	}

	if m, ok := finalModel.(*model); ok {

		selectedPaths := make([]string, 0, len(m.selected))
		for path := range m.selected {
			selectedPaths = append(selectedPaths, path)
		}

		if err := fileContextBuilder(selectedPaths, *outputFilename); err != nil {
			log.Fatalf("A critical error occurred while creating the context file: %v", err)
		}

	}
}
