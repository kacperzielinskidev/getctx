package main

import (
	"fmt"
	"os"
	"sort"
)

type processedFile struct {
	Path   string
	IsText bool
}

func HandleContextBuilder(m *model, outputFilename string) error {
	if len(m.selected) == 0 {
		fmt.Printf("%s No items selected. Exiting.\n", Icons.Error)
		return nil
	}

	selectedPaths := make([]string, 0, len(m.selected))
	for path := range m.selected {
		selectedPaths = append(selectedPaths, path)
	}

	return contextBuilder(selectedPaths, outputFilename)

}

func contextBuilder(selectedPaths []string, outputFilename string) error {
	processedFiles, err := discoverFiles(selectedPaths)
	if err != nil {
		return fmt.Errorf("error expanding selected paths: %w", err)
	}

	sort.Slice(processedFiles, func(i, j int) bool {
		return processedFiles[i].Path < processedFiles[j].Path
	})

	var textFiles []processedFile
	for _, file := range processedFiles {
		if file.IsText {
			textFiles = append(textFiles, file)
		}
	}

	fmt.Printf("%s Building context file: %s\n", Icons.Building, outputFilename)

	for _, file := range processedFiles {
		if file.IsText {
			fmt.Printf("   %s Adding content from: %s\n", Icons.Cursor, file.Path)

		} else {
			fmt.Printf("   %s %sSkipped content from: %s%s\n", Icons.Cursor, Colors.Red, file.Path, Colors.Reset)
		}
	}

	if len(textFiles) == 0 {
		fmt.Printf("\n%s No text files found to include. Output file was not created.\n", Icons.Info)
		return nil
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	for _, file := range textFiles {
		content, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Printf("%s Warning: Failed to read file %s: %v\n", Icons.Warning, file.Path, err)
			continue
		}
		header := fmt.Sprintf("--- START OF FILE: %s ---\n", file.Path)
		footer := fmt.Sprintf("\n--- END OF FILE: %s ---\n\n", file.Path)
		if _, err := outputFile.WriteString(header); err != nil {
			return fmt.Errorf("error writing header for file %s: %w", file.Path, err)
		}
		if _, err := outputFile.Write(content); err != nil {
			return fmt.Errorf("error writing content for file %s: %w", file.Path, err)
		}
		if _, err := outputFile.WriteString(footer); err != nil {
			return fmt.Errorf("error writing footer for file %s: %w", file.Path, err)
		}
	}

	fmt.Printf("%s Done! All content has been combined into %s\n", Icons.Done, outputFilename)
	return nil
}
