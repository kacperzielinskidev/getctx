// internal/app/completion.go

package tui

import (
	"path/filepath"
	"strings"
)

func (m *Model) getCompletionParts(input string) (dir, prefix string) {
	path := input

	if strings.HasPrefix(path, "~") {
		home, err := m.fsys.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(m.path, path)
	}
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return path, ""
	}

	dir = filepath.Dir(path)
	prefix = filepath.Base(path)
	return dir, prefix
}

func (m *Model) getCompletions(dir, prefix string) ([]string, error) {
	var matches []string

	entries, err := m.fsys.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prefix) {
			name := entry.Name()
			if entry.IsDir() {
				name += string(filepath.Separator) // Dodaj '/' do katalogów
			}
			matches = append(matches, name)
		}
	}
	return matches, nil
}

func findLongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for !strings.HasPrefix(strs[i], prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}
