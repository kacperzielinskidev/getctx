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

// A single file or directory item
type item struct {
	name  string
	isDir bool
}

// The main model for our TUI application.
// This is the absolute minimum state we need to display a list.
type model struct {
	path   string // The path we are viewing
	items  []item // The files and directories in the path
	cursor int    // Which item our cursor is pointing at
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
	// We want to work with a clean, absolute path.
	path, err := filepath.Abs(startPath)
	if err != nil {
		log.Fatalf("Could not get absolute path for '%s': %v", startPath, err)
	}

	// Read the items from the starting path.
	items, err := readDir(path)
	if err != nil {
		log.Fatalf("Could not read directory '%s': %v", path, err)
	}

	// Return a pointer to our new model.
	return &model{
		path:  path,
		items: items,
	}
}

// Init is the first command that can be run when the program starts.
// We don't need it to do anything right now.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update handles all user input and events.
// For now, it only handles moving the cursor and quitting.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// tea.KeyMsg is a message sent when the user presses a key.
	case tea.KeyMsg:
		switch msg.String() {
		// Keys to quit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Keys to move the cursor up.
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Keys to move the cursor down.
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime.
	return m, nil
}

// View is responsible for rendering the UI.
func (m *model) View() string {
	var s strings.Builder

	s.WriteString("File Lister (press 'q' to quit)\n")
	s.WriteString("Viewing path: " + m.path + "\n\n")

	// Iterate over our list of items.
	for i, item := range m.items {
		// The cursor shows which line we're on.
		cursor := " " // Default is a space
		if m.cursor == i {
			cursor = ">" // Set to ">" if this is the current line
		}

		// Get the item name.
		itemName := item.name
		// Add a slash to directories to make them identifiable.
		if item.isDir {
			itemName += "/"
		}

		// Render the final line.
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, itemName))
	}

	return s.String()
}

// main is the entry point for the program.
func main() {
	// Determine the starting path. Default to "." (current directory).
	startPath := "."
	// If the user provides an argument, use that as the starting path.
	if len(os.Args) > 1 {
		startPath = os.Args[1]
	}

	// Create a new Bubble Tea program.
	p := tea.NewProgram(initialModel(startPath))

	// Run the program.
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exited getctx.")
}
