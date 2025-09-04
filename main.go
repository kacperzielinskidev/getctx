package main

import (
	"flag"
	"fmt"
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
		if len(m.selected) == 0 {
			fmt.Println("❌ No files selected. Exiting.")
			return
		}

		selectedPaths := make([]string, 0, len(m.selected))
		for path := range m.selected {
			selectedPaths = append(selectedPaths, path)
		}

		if err := buildContextFile(selectedPaths, *outputFilename); err != nil {
			log.Fatalf("A critical error occurred while creating the context file: %v", err)
		}

		fmt.Printf("✅ Done! All content has been combined into %s\n", *outputFilename)
	}
}
