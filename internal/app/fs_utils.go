package app

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func discoverFiles(paths []string, excludedNames map[string]struct{}) ([]string, error) {
	var discoveredPaths []string

	for _, path := range paths {
		if _, ok := excludedNames[filepath.Base(path)]; ok {
			continue
		}

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

				if _, ok := excludedNames[d.Name()]; ok {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				if !d.IsDir() {
					discoveredPaths = append(discoveredPaths, subPath)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("%s Warning: Error walking directory %s: %v\n", Icons.Warning, path, err)
			}
		} else {
			discoveredPaths = append(discoveredPaths, path)
		}
	}

	return discoveredPaths, nil
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
