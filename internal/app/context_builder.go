// Plik: internal/app/context_builder.go
package app

import (
	"fmt"
	"sort"
)

// HandleContextBuilder pobiera FileSystem z modelu i przekazuje go do logiki budowania.
func HandleContextBuilder(m *Model, outputFilename string) error {
	if len(m.selected) == 0 {
		fmt.Printf("%s No items selected. Exiting.\n", Icons.Error)
		return nil
	}

	selectedPaths := make([]string, 0, len(m.selected))
	for path := range m.selected {
		selectedPaths = append(selectedPaths, path)
	}

	return contextBuilder(selectedPaths, outputFilename, m.config, m.fsys)
}

func contextBuilder(selectedPaths []string, outputFilename string, config *Config, fsys FileSystem) error {
	acceptableFiles, warnings, err := discoverFiles(fsys, selectedPaths, config.ExcludedNames)
	if err != nil {
		return fmt.Errorf("error discovering files: %w", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("%s Some paths were skipped due to errors:\n", Icons.Warning)
		for _, warn := range warnings {
			fmt.Printf("   - %s\n", warn)
		}
	}

	var textFiles []string
	for _, path := range acceptableFiles {
		isText, err := isTextFile(fsys, path)
		if err != nil {
			fmt.Printf("%s Warning: Could not check file type for %s: %v\n", Icons.Warning, path, err)
			continue
		}
		if isText {
			textFiles = append(textFiles, path)
		}
	}

	skippedFileCount := len(acceptableFiles) - len(textFiles)

	if len(textFiles) == 0 {
		message := fmt.Sprintf("\n%s No text files found to include.", Icons.Info)
		if skippedFileCount > 0 {
			message += fmt.Sprintf(" %d file(s) were skipped (non-text or unreadable).", skippedFileCount)
		}
		message += " Output file was not created."
		fmt.Println(message)
		return nil
	}

	fmt.Printf("%s Building context file: %s\n", Icons.Building, outputFilename)

	if skippedFileCount > 0 {
		fmt.Printf("   %s Skipped %d non-text or unreadable file(s).\n", Icons.Info, skippedFileCount)
	}

	sort.Strings(textFiles)
	for _, path := range textFiles {
		fmt.Printf("   %s Adding content from: %s\n", Icons.Cursor, path)
	}

	// Zmieniono os.Create na fs.Create.
	outputFile, err := fsys.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	for _, path := range textFiles {
		// Zmieniono os.ReadFile na fs.ReadFile.
		content, err := fsys.ReadFile(path)
		if err != nil {
			fmt.Printf("%s Warning: Failed to read file %s: %v\n", Icons.Warning, path, err)
			continue
		}
		header := fmt.Sprintf("--- START OF FILE: %s ---\n", path)
		footer := fmt.Sprintf("\n--- END OF FILE: %s ---\n\n", path)
		if _, err := outputFile.Write([]byte(header)); err != nil {
			return fmt.Errorf("error writing header for file %s: %w", path, err)
		}
		if _, err := outputFile.Write(content); err != nil {
			return fmt.Errorf("error writing content for file %s: %w", path, err)
		}
		if _, err := outputFile.Write([]byte(footer)); err != nil {
			return fmt.Errorf("error writing footer for file %s: %w", path, err)
		}
	}

	fmt.Printf("%s Done! All content has been combined into %s\n", Icons.Done, outputFilename)
	return nil
}
