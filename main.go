// File: main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// readDir is a helper function to get the contents of a directory.
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
		selected: make(map[string]struct{}), // Initialize the map
	}
}

// Init is the first command that can be run when the program starts.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update handles all user input and events.
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

		// Enter key: Enter a directory
		case "enter":
			if len(m.items) == 0 {
				break
			}
			selectedItem := m.items[m.cursor]
			if selectedItem.isDir {
				newPath := filepath.Join(m.path, selectedItem.name)
				newItems, err := readDir(newPath)
				if err != nil {
					// In a real app, you might want to show an error message to the user
					log.Printf("Error reading directory %s: %v", newPath, err)
					break
				}
				m.path = newPath
				m.items = newItems
				m.cursor = 0 // Reset cursor to the top of the new directory
			}

		// Backspace key: Go up to the parent directory
		case "backspace":
			parentPath := filepath.Dir(m.path)
			// Avoid going "up" from the root directory
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

		// Space key: Toggle selection for a file
		case " ":
			if len(m.items) == 0 {
				break
			}
			selectedItem := m.items[m.cursor]
			// We can only select files, not directories
			if !selectedItem.isDir {
				fullPath := filepath.Join(m.path, selectedItem.name)
				// Check if the item is already selected
				if _, ok := m.selected[fullPath]; ok {
					// If it is, unselect it
					delete(m.selected, fullPath)
				} else {
					// If it's not, select it
					m.selected[fullPath] = struct{}{} // Use empty struct for set-like behavior
				}
			}
		}
	}

	return m, nil
}

// View is responsible for rendering the UI.
func (m *model) View() string {
	var s strings.Builder

	s.WriteString("Wybierz pliki do kontekstu (spacja - zaznacz, enter - wejdź, backspace - cofnij, q - wyjdź)\n")
	s.WriteString("Aktualna ścieżka: " + m.path + "\n\n")

	for i, item := range m.items {
		cursor := " " // Default cursor
		if m.cursor == i {
			cursor = ">"
		}

		// Determine the prefix for the item (checkbox or empty space)
		prefix := "   " // Default for directories
		if !item.isDir {
			fullPath := filepath.Join(m.path, item.name)
			if _, ok := m.selected[fullPath]; ok {
				prefix = "[x]" // Selected file
			} else {
				prefix = "[ ]" // Unselected file
			}
		}

		// Get the item name.
		itemName := item.name
		if item.isDir {
			itemName += "/" // Add a slash to directories
		}

		s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, prefix, itemName))
	}

	// Add a footer with the count of selected files
	s.WriteString(fmt.Sprintf("\nZaznaczono %d plików. Naciśnij 'q' aby zakończyć i wypisać ścieżki.", len(m.selected)))

	return s.String()
}

// main is the entry point for the program.
func main() {
	startPath := "."
	if len(os.Args) > 1 {
		startPath = os.Args[1]
	}

	p := tea.NewProgram(initialModel(startPath))

	// Run returns the final model. We can use it to get the selected files.
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Type assertion to get our specific model type
	if m, ok := finalModel.(*model); ok {
		// If any files were selected, print their paths to stdout
		if len(m.selected) > 0 {
			for path := range m.selected {
				fmt.Println(path)
			}
		}
	}
}
