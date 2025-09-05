// File: file_builder.go
package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func fileContextBuilder(selectedPaths []string, outputFilename string) error {

	allFilePaths, skippedBinaryCount, err := expandPaths(selectedPaths)
	if err != nil {
		return fmt.Errorf("error expanding selected paths: %w", err)

	}

	if len(allFilePaths) == 0 {
		if skippedBinaryCount > 0 {
			fmt.Printf("ℹ️ No text files found to include. %d binary file(s) were skipped.\n", skippedBinaryCount)
		} else {
			fmt.Println("ℹ️ No text files found in the selected paths. Output file was not created.")
		}
		return nil
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	fmt.Printf("🚀 Building context file: %s\n", outputFilename)

	sort.Strings(allFilePaths)

	for _, path := range allFilePaths {
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

	fmt.Printf("✅ Done! All content has been combined into %s\n", outputFilename)

	return nil
}

func expandPaths(paths []string) (textFiles []string, skippedCount int, err error) {
	fileSet := make(map[string]struct{})

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("⚠️ Warning: Could not stat path %s: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			err := filepath.WalkDir(path, func(subPath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !d.IsDir() {
					isText, err := isTextFile(subPath)
					if err != nil {
						fmt.Printf("⚠️ Warning: Could not determine file type for %s: %v\n", subPath, err)
						return nil
					}

					if isText {
						fileSet[subPath] = struct{}{}
					} else {
						skippedCount++
					}

				}
				return nil
			})

			if err != nil {
				fmt.Printf("⚠️ Warning: Error walking directory %s: %v\n", path, err)

			}
		} else {
			isText, err := isTextFile(path)
			if err != nil {
				fmt.Printf("⚠️ Warning: Could not determine file type for %s: %v\n", path, err)

			} else if isText {

				fileSet[path] = struct{}{}
			} else {
				skippedCount++
			}

		}
	}

	// Konwertuj mapę (set) z powrotem na plasterek (slice)
	finalPaths := make([]string, 0, len(fileSet))
	for path := range fileSet {
		finalPaths = append(finalPaths, path)
	}

	return finalPaths, skippedCount, nil
}

func isTextFile(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	contentType := http.DetectContentType(buffer[:n])

	return strings.HasPrefix(contentType, "text/"), nil
}
