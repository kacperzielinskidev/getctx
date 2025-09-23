package tui

import (
	"path/filepath"
	"runtime"
	"strings"
)

func (m *Model) getCompletionParts(input string) (dir, prefix string) {
	path := input

	// 1. Obsługa katalogu domowego (bez zmian)
	if strings.HasPrefix(path, "~") {
		home, err := m.fsys.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}

	var analysisPath string

	// 2. KLUCZOWA POPRAWKA LOGIKI
	if filepath.IsAbs(path) {
		// Ścieżka jest już absolutna (np. C:\Users, /home/user). Użyj jej.
		analysisPath = path
	} else if runtime.GOOS == "windows" && (path == `\` || path == `/`) {
		// --- TO JEST SERCE POPRAWKI DLA WINDOWS ---
		// Jeśli użytkownik wpisał sam separator, jego intencją jest
		// przejście do roota dysku, na którym aktualnie się znajduje w TUI.
		// Wyciągamy nazwę woluminu z `m.path` (np. "C:") i dodajemy separator.
		analysisPath = filepath.VolumeName(m.path) + string(filepath.Separator)
	} else {
		// W każdym innym przypadku jest to ścieżka względna.
		// Łączymy ją z bieżącym katalogiem w TUI.
		analysisPath = filepath.Join(m.path, path)
	}

	// 3. Ustalenie katalogu do przeszukania i prefiksu (bez zmian)
	if strings.HasSuffix(input, string(filepath.Separator)) {
		return filepath.Clean(analysisPath), ""
	}

	dir = filepath.Dir(analysisPath)
	prefix = filepath.Base(analysisPath)

	return filepath.Clean(dir), prefix
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
				name += string(filepath.Separator)
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
