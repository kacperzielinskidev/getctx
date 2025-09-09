// Plik: internal/app/context_builder.go
package app

import (
	"fmt"
	"sort"
)

// BuildContext tworzy plik kontekstowy na podstawie zaznaczonych elementów w modelu.
// Funkcja ta jest jedynym publicznym punktem wejścia dla logiki budowania kontekstu.
func BuildContext(m *Model, outputFilename string) error {
	// Krok 1: Sprawdzenie, czy cokolwiek zostało zaznaczone.
	if len(m.selected) == 0 {
		fmt.Printf("%s No items selected. Exiting.\n", Icons.Error)
		return nil
	}

	// Krok 2: Wyodrębnienie ścieżek z mapy w modelu.
	selectedPaths := make([]string, 0, len(m.selected))
	for path := range m.selected {
		selectedPaths = append(selectedPaths, path)
	}

	// Krok 3: Odkrycie wszystkich plików w podanych ścieżkach.
	// Używamy FileSystem i Config bezpośrednio z modelu.
	fs := m.fsys
	config := m.config
	acceptableFiles, warnings, err := discoverFiles(fs, selectedPaths, config.ExcludedNames)
	if err != nil {
		return fmt.Errorf("error discovering files: %w", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("%s Some paths were skipped due to errors:\n", Icons.Warning)
		for _, warn := range warnings {
			fmt.Printf("   - %s\n", warn)
		}
	}

	// Krok 4: Filtrowanie, aby zostawić tylko pliki tekstowe.
	var textFiles []string
	for _, path := range acceptableFiles {
		isText, err := isTextFile(fs, path)
		if err != nil {
			fmt.Printf("%s Warning: Could not check file type for %s: %v\n", Icons.Warning, path, err)
			continue
		}
		if isText {
			textFiles = append(textFiles, path)
		}
	}

	skippedFileCount := len(acceptableFiles) - len(textFiles)

	// Krok 5: Obsługa przypadku, gdy nie znaleziono żadnych plików tekstowych.
	if len(textFiles) == 0 {
		message := fmt.Sprintf("\n%s No text files found to include.", Icons.Info)
		if skippedFileCount > 0 {
			message += fmt.Sprintf(" %d file(s) were skipped (non-text or unreadable).", skippedFileCount)
		}
		message += " Output file was not created."
		fmt.Println(message)
		return nil
	}

	// Krok 6: Zapisywanie zawartości do pliku wyjściowego.
	fmt.Printf("%s Building context file: %s\n", Icons.Building, outputFilename)
	if skippedFileCount > 0 {
		fmt.Printf("   %s Skipped %d non-text or unreadable file(s).\n", Icons.Info, skippedFileCount)
	}

	sort.Strings(textFiles)
	for _, path := range textFiles {
		fmt.Printf("   %s Adding content from: %s\n", Icons.Cursor, path)
	}

	outputFile, err := fs.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	for _, path := range textFiles {
		content, err := fs.ReadFile(path)
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
