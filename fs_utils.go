// File: fs_utils.go
package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Funkcja readDir została usunięta, ponieważ była zbędną abstrakcją.
// Należy używać os.ReadDir bezpośrednio.

// discoverFiles rekurencyjnie przeszukuje podane ścieżki (pliki i katalogi),
// znajduje wszystkie pliki i klasyfikuje je jako tekstowe lub binarne.
func discoverFiles(paths []string) ([]processedFile, error) {
	var allFiles []processedFile

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("%s Warning: Could not stat path %s: %v\n", Icons.Warning, path, err)
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
						fmt.Printf("%s Warning: Could not determine file type for %s: %v\n", Icons.Warning, subPath, err)
						allFiles = append(allFiles, processedFile{Path: subPath, IsText: false})
						return nil
					}
					allFiles = append(allFiles, processedFile{Path: subPath, IsText: isText})
				}
				return nil
			})
			if err != nil {
				fmt.Printf("%s Warning: Error walking directory %s: %v\n", Icons.Warning, path, err)
			}
		} else {
			isText, err := isTextFile(path)
			if err != nil {
				fmt.Printf("%s Warning: Could not determine file type for %s: %v\n", Icons.Warning, path, err)
				allFiles = append(allFiles, processedFile{Path: path, IsText: false})
			} else {
				allFiles = append(allFiles, processedFile{Path: path, IsText: isText})
			}
		}
	}

	return allFiles, nil
}

// isTextFile sprawdza, czy plik jest prawdopodobnie plikiem tekstowym.
// Nazwa jest już zgodna z konwencjami Go i pozostaje bez zmian.
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
