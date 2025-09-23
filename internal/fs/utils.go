package fs

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

func DiscoverFiles(fsys FileSystem, paths []string, excludedNames map[string]struct{}) ([]string, []string, error) {
	var discoveredPaths []string
	var warnings []string

	for _, path := range paths {
		if _, ok := excludedNames[filepath.Base(path)]; ok {
			continue
		}

		info, err := fsys.Stat(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("Could not stat path %s: %v", path, err))
			continue
		}

		if info.IsDir() {
			err := fsys.WalkDir(path, func(subPath string, d fs.DirEntry, err error) error {
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
				warnings = append(warnings, fmt.Sprintf("Error walking directory %s: %v", path, err))
			}
		} else {
			discoveredPaths = append(discoveredPaths, path)
		}
	}

	return discoveredPaths, warnings, nil
}

func IsTextFile(fsys FileSystem, path string) (bool, error) {
	file, err := fsys.Open(path)
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
