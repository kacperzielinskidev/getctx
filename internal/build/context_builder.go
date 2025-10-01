package build

import (
	"fmt"
	"io"
	"sort"

	"github.com/kacperzielinskidev/getctx/internal/config"
	"github.com/kacperzielinskidev/getctx/internal/fs"
	"github.com/kacperzielinskidev/getctx/internal/logger"
)

const (
	fileHeaderFormat = "--- START OF FILE: %s ---\n"
	fileFooterFormat = "\n--- END OF FILE: %s ---\n\n"
)

type ContextBuilder struct {
	log    *logger.Logger
	fsys   fs.FileSystem
	config *config.Config
}

type BuildResult struct {
	FilesProcessed int
	FilesSkipped   int
	PathsWithErr   []string
}

func NewContextBuilder(log *logger.Logger, fsys fs.FileSystem, cfg *config.Config) *ContextBuilder {
	return &ContextBuilder{
		log:    log,
		fsys:   fsys,
		config: cfg,
	}
}

func (cb *ContextBuilder) Build(selectedPaths []string, outputFilename string) (*BuildResult, error) {
	if len(selectedPaths) == 0 {
		cb.log.Info("BuildContext", "No items selected by user, exiting.")
		return &BuildResult{}, nil
	}

	cb.log.Debug("BuildContext", map[string]any{
		"selected_path_count": len(selectedPaths),
		"selected_paths":      selectedPaths,
	})

	allFiles, warnings, err := fs.DiscoverFiles(cb.fsys, selectedPaths, cb.config.ExcludedNames)
	if err != nil {
		cb.log.Error("BuildContext.DiscoverFiles", err)
		return nil, fmt.Errorf("error discovering files: %w", err)
	}

	textFiles := cb.filterTextFiles(allFiles)

	result := &BuildResult{
		FilesProcessed: len(textFiles),
		FilesSkipped:   len(allFiles) - len(textFiles),
		PathsWithErr:   warnings,
	}

	if len(textFiles) == 0 {
		cb.log.Info("BuildContext", "No text files found to process.")
		return result, nil
	}

	cb.log.Info("BuildContext", map[string]any{
		"files_to_process_count": len(textFiles),
		"output_filename":        outputFilename,
	})

	err = cb.writeContextFile(outputFilename, textFiles)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cb *ContextBuilder) filterTextFiles(files []string) []string {
	var textFiles []string
	for _, path := range files {
		isText, err := fs.IsTextFile(cb.fsys, path)
		if err != nil {
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

func (cb *ContextBuilder) writeContextFile(outputFilename string, files []string) error {
	outputFile, err := cb.fsys.Create(outputFilename)
	if err != nil {
		cb.log.Error("writeContextFile.Create", err)
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	sort.Strings(files)

	for _, path := range files {
		if err := cb.appendFileToContext(outputFile, path); err != nil {
			// Log the warning but continue processing other files.
			cb.log.Warn("writeContextFile.append", map[string]any{
				"message": "Failed to append file to context, skipping",
				"path":    path,
				"error":   err.Error(),
			})
		}
	}

	return nil
}

func (cb *ContextBuilder) appendFileToContext(writer io.Writer, path string) error {
	content, err := cb.fsys.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}

	// Write header
	if _, err := fmt.Fprintf(writer, fileHeaderFormat, path); err != nil {
		return fmt.Errorf("error writing header for file %s: %w", path, err)
	}

	// Write content
	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("error writing content for file %s: %w", path, err)
	}

	// Write footer
	if _, err := fmt.Fprintf(writer, fileFooterFormat, path); err != nil {
		return fmt.Errorf("error writing footer for file %s: %w", path, err)
	}

	return nil
}
