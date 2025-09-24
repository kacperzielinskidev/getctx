package build

import (
	"fmt"
	"getctx/internal/config"
	"getctx/internal/fs"
	"getctx/internal/logger"
	"io"
	"sort"
)

type ContextBuilder struct {
	log            *logger.Logger
	fsys           fs.FileSystem
	config         *config.Config
	outputFilename string
}

func NewContextBuilder(log *logger.Logger, fsys fs.FileSystem, outputFilename string, cfg *config.Config) *ContextBuilder {
	return &ContextBuilder{
		log:            log,
		fsys:           fsys,
		config:         cfg,
		outputFilename: outputFilename,
	}
}

func (cb *ContextBuilder) Build(selectedPaths []string) error {
	if len(selectedPaths) == 0 {
		fmt.Println("ℹ️ No items selected. Exiting.")
		cb.log.Info("BuildContext", "No items selected by user, exiting.")
		return nil
	}

	cb.log.Debug("BuildContext", map[string]any{
		"selected_path_count": len(selectedPaths),
		"selected_paths":      selectedPaths,
	})

	processableFiles, err := cb.findProcessableFiles(selectedPaths)
	if err != nil {
		cb.log.Error("BuildContext.findProcessableFiles", err)
		return err
	}

	textFiles := cb.filterTextFiles(processableFiles)

	if !cb.handleNoTextFilesFound(len(processableFiles), len(textFiles)) {
		cb.log.Info("BuildContext", "No text files found to process.")
		return nil
	}

	cb.log.Info("BuildContext", map[string]any{
		"files_to_process_count": len(textFiles),
		"output_filename":        cb.outputFilename,
	})
	return cb.writeContextFile(textFiles)
}

func (cb *ContextBuilder) findProcessableFiles(paths []string) ([]string, error) {
	files, warnings, err := fs.DiscoverFiles(cb.fsys, paths, cb.config.ExcludedNames)
	if err != nil {
		err = fmt.Errorf("error discovering files: %w", err)
		cb.log.Error("findProcessableFiles.discoverFiles", err)
		return nil, err
	}

	if len(warnings) > 0 {
		fmt.Println("⚠️ Some paths were skipped due to errors:")
		for _, warn := range warnings {
			cb.log.Warn("findProcessableFiles", warn)
			fmt.Printf("   - %s\n", warn)
		}
	}
	return files, nil
}

func (cb *ContextBuilder) filterTextFiles(files []string) []string {
	var textFiles []string
	for _, path := range files {
		isText, err := fs.IsTextFile(cb.fsys, path)
		if err != nil {
			fmt.Printf("⚠️ Warning: Could not check file type for %s: %v\n", path, err)
			cb.log.Warn("filterTextFiles", map[string]any{
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

func (cb *ContextBuilder) handleNoTextFilesFound(totalFiles, textFiles int) bool {
	if textFiles > 0 {
		return true
	}

	skippedFileCount := totalFiles - textFiles
	message := "\nℹ️ No text files found to include."
	if skippedFileCount > 0 {
		message += fmt.Sprintf(" %d file(s) were skipped (non-text or unreadable).", skippedFileCount)
	}
	message += " Output file was not created."
	fmt.Println(message)

	return false
}

func (cb *ContextBuilder) writeContextFile(files []string) error {
	fmt.Printf("🚀 Building context file: %s\n", cb.outputFilename)
	sort.Strings(files)
	for _, path := range files {
		fmt.Printf("   ❯ Adding content from: %s\n", path)
	}

	outputFile, err := cb.fsys.Create(cb.outputFilename)
	if err != nil {
		cb.log.Error("writeContextFile.Create", err)
		return fmt.Errorf("failed to create output file %s: %w", cb.outputFilename, err)
	}
	defer outputFile.Close()

	for _, path := range files {
		if err := cb.appendFileToContext(outputFile, path); err != nil {
			fmt.Printf("⚠️ Warning: Failed to append file %s: %v\n", path, err)
			cb.log.Warn("writeContextFile.append", map[string]any{
				"message": "Failed to append file to context",
				"path":    path,
				"error":   err.Error(),
			})
		}
	}
	fmt.Printf("✅ Done! All content has been combined into %s\n", cb.outputFilename)
	return nil
}

func (cb *ContextBuilder) appendFileToContext(writer io.Writer, path string) error {
	content, err := cb.fsys.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}

	chunks := []struct {
		data        []byte
		description string
	}{
		{[]byte(fmt.Appendf(nil, "--- START OF FILE: %s ---\n", path)), "header"},
		{content, "content"},
		{[]byte(fmt.Appendf(nil, "\n--- END OF FILE: %s ---\n\n", path)), "footer"},
	}

	for _, chunk := range chunks {
		if _, err := writer.Write(chunk.data); err != nil {
			return fmt.Errorf("error writing %s for file %s: %w", chunk.description, path, err)
		}
	}

	return nil
}
