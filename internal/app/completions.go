// internal/app/completion.go

package app

import (
	"path/filepath"
	"strings"
)

// getCompletionParts analizuje ścieżkę wejściową i zwraca katalog do przeszukania oraz prefiks pliku/katalogu.
func (m *Model) getCompletionParts(input string) (dir, prefix string) {
	// Ścieżka, którą będziemy analizować. Zaczynamy od tego, co wpisał użytkownik.
	path := input

	// --- NOWA, POPRAWIONA LOGIKA ---

	// 1. Obsługa specjalnego przypadku `~` (katalog domowy).
	if strings.HasPrefix(path, "~") {
		home, err := m.fsys.UserHomeDir()
		if err == nil {
			// Zamień `~` na pełną ścieżkę do katalogu domowego.
			path = filepath.Join(home, path[1:])
		}
	}

	// 2. Sprawdzenie, czy ścieżka jest względna.
	if !filepath.IsAbs(path) {
		// Jeśli tak, połącz ją z aktualną ścieżką TUI, aby uzyskać pełną, absolutną ścieżkę.
		// To jest kluczowa poprawka, która rozwiązuje zgłoszony problem.
		path = filepath.Join(m.path, path)
	}
	// W tym momencie `path` jest już zawsze poprawną, absolutną ścieżką, którą próbuje wpisać użytkownik.

	// --- KONIEC NOWEJ LOGIKI ---

	// 3. Rozdzielenie ścieżki na katalog do przeszukania i prefiks.
	// Jeśli ścieżka kończy się separatorem (np. "/home/user/"),
	// oznacza to, że chcemy listę zawartości tego katalogu.
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return path, "" // Przeszukaj `path`, prefiks jest pusty.
	}

	// Jeśli ścieżka nie kończy się separatorem (np. "/home/user/Doc"),
	// oznacza to, że chcemy znaleźć dopasowania dla "Doc" wewnątrz "/home/user".
	dir = filepath.Dir(path)
	prefix = filepath.Base(path)
	return dir, prefix
}

// getCompletions przeszukuje dany katalog w poszukiwaniu wpisów pasujących do prefiksu.
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

// findLongestCommonPrefix znajduje najdłuższy wspólny prefiks w tablicy stringów.
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
