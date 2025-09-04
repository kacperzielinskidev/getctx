// File: file_builder.go
package main

import (
	"fmt"
	"os"
	"sort"
)

func buildContextFile(selectedPaths []string, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	fmt.Printf("🚀 Building context file: %s\n", outputFilename)

	sort.Strings(selectedPaths)

	for _, path := range selectedPaths {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to read file %s: %v\n", path, err)
			continue
		}

		fmt.Printf("   -> Adding content from: %s\n", path)

		header := fmt.Sprintf("--- START OF FILE: %s ---\n", path)
		footer := fmt.Sprintf("\n--- END OF FILE: %s ---\n\n", path)

		if _, err := outputFile.WriteString(header); err != nil {
			return fmt.Errorf("error writing header for file %s: %w", path, err)
		}
		if _, err := outputFile.Write(content); err != nil {
			return fmt.Errorf("error writing content for file %s: %w", path, err)
		}
		if _, err := outputFile.WriteString(footer); err != nil {
			return fmt.Errorf("error writing footer for file %s: %w", path, err)
		}
	}

	return nil
}
