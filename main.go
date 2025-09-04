// File: main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name  string
	isDir bool
}

// The main model for our TUI application.
type model struct {
	path     string              // The path we are viewing
	items    []item              // The files and directories in the path
	cursor   int                 // Which item our cursor is pointing at
	selected map[string]struct{} // A set of selected file paths
}

func readDir(path string) ([]item, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []item
	for _, file := range files {
		items = append(items, item{name: file.Name(), isDir: file.IsDir()})
	}
	return items, nil
}

// initialModel sets up the very first state of our application.
func initialModel(startPath string) *model {
	path, err := filepath.Abs(startPath)
	if err != nil {
		log.Fatalf("Could not get absolute path for '%s': %v", startPath, err)
	}

	items, err := readDir(path)
	if err != nil {
		log.Fatalf("Could not read directory '%s': %v", path, err)
	}

	return &model{
		path:     path,
		items:    items,
		selected: make(map[string]struct{}),
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.items) == 0 {
				break
			}
			selectedItem := m.items[m.cursor]
			if selectedItem.isDir {
				newPath := filepath.Join(m.path, selectedItem.name)
				newItems, err := readDir(newPath)
				if err != nil {
					log.Printf("Error reading directory %s: %v", newPath, err)
					break
				}
				m.path = newPath
				m.items = newItems
				m.cursor = 0
			}

		case "backspace":
			parentPath := filepath.Dir(m.path)
			if parentPath != m.path {
				newItems, err := readDir(parentPath)
				if err != nil {
					log.Printf("Error reading directory %s: %v", parentPath, err)
					break
				}
				m.path = parentPath
				m.items = newItems
				m.cursor = 0
			}

		case " ":
			if len(m.items) == 0 {
				break
			}
			selectedItem := m.items[m.cursor]
			if !selectedItem.isDir {
				fullPath := filepath.Join(m.path, selectedItem.name)
				if _, ok := m.selected[fullPath]; ok {
					delete(m.selected, fullPath)
				} else {
					m.selected[fullPath] = struct{}{}
				}
			}
		}
	}

	return m, nil
}

func (m *model) View() string {
	var s strings.Builder

	s.WriteString("Wybierz pliki do kontekstu (spacja - zaznacz, enter - wejdź, backspace - cofnij, q - zapisz i wyjdź)\n")
	s.WriteString("Aktualna ścieżka: " + m.path + "\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		prefix := "   "
		if !item.isDir {
			fullPath := filepath.Join(m.path, item.name)
			if _, ok := m.selected[fullPath]; ok {
				prefix = "[x]"
			} else {
				prefix = "[ ]"
			}
		}

		itemName := item.name
		if item.isDir {
			itemName += "/"
		}

		s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, prefix, itemName))
	}

	s.WriteString(fmt.Sprintf("\nZaznaczono %d plików. Naciśnij 'q' aby zapisać i zakończyć.", len(m.selected)))

	return s.String()
}

// createContextFile takes a list of file paths and an output file name,
// and concatenates the content of the source files into the output file.
func createContextFile(selectedPaths []string, outputFilename string) error {
	// Create or truncate the output file.
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("nie udało się stworzyć pliku wyjściowego %s: %w", outputFilename, err)
	}
	defer outputFile.Close()

	fmt.Printf("🚀 Rozpoczynam budowanie pliku kontekstu: %s\n", outputFilename)

	// Sort paths for consistent output
	sort.Strings(selectedPaths)

	for _, path := range selectedPaths {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("⚠️ Ostrzeżenie: Nie udało się odczytać pliku %s: %v\n", path, err)
			continue // Skip this file and continue with others
		}

		fmt.Printf("   -> Dodawanie zawartości z: %s\n", path)

		header := fmt.Sprintf("--- START OF FILE: %s ---\n", path)
		footer := fmt.Sprintf("\n--- END OF FILE: %s ---\n\n", path)

		if _, err := outputFile.WriteString(header); err != nil {
			return fmt.Errorf("błąd zapisu nagłówka dla pliku %s: %w", path, err)
		}
		if _, err := outputFile.Write(content); err != nil {
			return fmt.Errorf("błąd zapisu zawartości dla pliku %s: %w", path, err)
		}
		if _, err := outputFile.WriteString(footer); err != nil {
			return fmt.Errorf("błąd zapisu stopki dla pliku %s: %w", path, err)
		}
	}

	return nil
}

func main() {
	// Define and parse command-line flags
	outputFilename := flag.String("o", "context.txt", "Nazwa pliku wyjściowego")
	flag.Parse()

	// The starting path is the first non-flag argument, or "." if not provided.
	startPath := "."
	if flag.NArg() > 0 {
		startPath = flag.Arg(0)
	}

	p := tea.NewProgram(initialModel(startPath))

	// Run returns the final model. We can use it to get the selected files.
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Wystąpił błąd podczas uruchamiania programu: %v", err)
	}

	// Type assertion to get our specific model type
	if m, ok := finalModel.(*model); ok {
		// If no files were selected, inform the user and exit.
		if len(m.selected) == 0 {
			fmt.Println("❌ Nie wybrano żadnych plików. Program zakończył działanie.")
			return
		}

		// Convert map keys to a slice
		selectedPaths := make([]string, 0, len(m.selected))
		for path := range m.selected {
			selectedPaths = append(selectedPaths, path)
		}

		// Create the context file
		if err := createContextFile(selectedPaths, *outputFilename); err != nil {
			log.Fatalf("Błąd krytyczny podczas tworzenia pliku kontekstu: %v", err)
		}

		fmt.Printf("✅ Gotowe! Cała zawartość została połączona w pliku %s\n", *outputFilename)
	}
}
