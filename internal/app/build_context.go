package app

import (
	"fmt"
	"getctx/internal/logger"
	"io"
	"sort"
)

func BuildContext(m *Model, outputFilename string) error {
	if len(m.selected) == 0 {
		fmt.Printf("%s No items selected. Exiting.\n", Icons.Error)
		logger.Info("BuildContext", "No items selected by user, exiting.")
		return nil
	}

	selectedPaths := extractSelectedPaths(m)
	logger.Debug("BuildContext", map[string]any{
		"selected_path_count": len(selectedPaths),
		"selected_paths":      selectedPaths,
	})

	processableFiles, err := findProcessableFiles(m.fsys, selectedPaths, m.config)
	if err != nil {
		logger.Error("BuildContext.findProcessableFiles", err)
		return err
	}

	textFiles := filterTextFiles(m.fsys, processableFiles)

	if !handleNoTextFilesFound(len(processableFiles), len(textFiles)) {
		logger.Info("BuildContext", "No text files found to process.")
		return nil
	}

	logger.Info("BuildContext", map[string]any{
		"files_to_process_count": len(textFiles),
		"output_filename":        outputFilename,
	})
	return writeContextFile(m.fsys, outputFilename, textFiles)
}

func extractSelectedPaths(m *Model) []string {
	paths := make([]string, 0, len(m.selected))
	for path := range m.selected {
		paths = append(paths, path)
	}
	return paths
}

func findProcessableFiles(fsys FileSystem, paths []string, config *Config) ([]string, error) {
	files, warnings, err := discoverFiles(fsys, paths, config.ExcludedNames)
	if err != nil {
		err = fmt.Errorf("error discovering files: %w", err)
		logger.Error("findProcessableFiles.discoverFiles", err)
		return nil, err
	}

	if len(warnings) > 0 {
		fmt.Printf("%s Some paths were skipped due to errors:\n", Icons.Warning)
		for _, warn := range warnings {
			logger.Warn("findProcessableFiles", warn)
			fmt.Printf("   - %s\n", warn)
		}
	}
	return files, nil

}

func filterTextFiles(fsys FileSystem, files []string) []string {
	var textFiles []string
	for _, path := range files {
		isText, err := isTextFile(fsys, path)
		if err != nil {
			fmt.Printf("%s Warning: Could not check file type for %s: %v\n", Icons.Warning, path, err)
			logger.Warn("filterTextFiles", map[string]any{
				"message": "Could not check file type",
				"path":    path,
				"error":   err.Error(),
			})
			continue
		}
		if isText {
			textFiles = append(textFiles, path)
		}
	}
	return textFiles
}

func handleNoTextFilesFound(totalFiles, textFiles int) bool {
	if textFiles > 0 {
		return true
	}

	skippedFileCount := totalFiles - textFiles
	message := fmt.Sprintf("\n%s No text files found to include.", Icons.Info)
	if skippedFileCount > 0 {
		message += fmt.Sprintf(" %d file(s) were skipped (non-text or unreadable).", skippedFileCount)
	}
	message += " Output file was not created."
	fmt.Println(message)

	return false
}

func writeContextFile(fsys FileSystem, filename string, files []string) error {
	fmt.Printf("%s Building context file: %s\n", Icons.Building, filename)
	sort.Strings(files)
	for _, path := range files {
		fmt.Printf("   %s Adding content from: %s\n", Icons.Cursor, path)
	}

	outputFile, err := fsys.Create(filename)
	if err != nil {
		logger.Error("writeContextFile.Create", err)
		return fmt.Errorf("failed to create output file %s: %w", filename, err)
	}
	defer outputFile.Close()

	for _, path := range files {
		if err := appendFileToContext(outputFile, fsys, path); err != nil {
			fmt.Printf("%s Warning: Failed to append file %s: %v\n", Icons.Warning, path, err)
			logger.Warn("writeContextFile.append", map[string]any{
				"message": "Failed to append file to context",
				"path":    path,
				"error":   err.Error(),
			})
		}
	}
	fmt.Printf("%s Done! All content has been combined into %s\n", Icons.Done, filename)
	return nil
}

func appendFileToContext(writer io.Writer, fsys FileSystem, path string) error {
	content, err := fsys.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}

	chunks := []struct {
		data        []byte
		description string
	}{
		{fmt.Appendf(nil, "--- START OF FILE: %s ---\n", path), "header"},
		{content, "content"},
		{fmt.Appendf(nil, "\n--- END OF FILE: %s ---\n\n", path), "footer"},
	}

	for _, chunk := range chunks {
		if _, err := writer.Write(chunk.data); err != nil {
			return fmt.Errorf("error writing %s for file %s: %w", chunk.description, path, err)
		}
	}

	return nil

}
