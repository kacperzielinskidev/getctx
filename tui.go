// File: tui.go
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

type model struct {
	path     string
	items    []item
	cursor   int
	selected map[string]struct{}
}

func newModel(startPath string) *model {
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
			fullPath := filepath.Join(m.path, selectedItem.name)
			if _, ok := m.selected[fullPath]; ok {
				delete(m.selected, fullPath)
			} else {
				m.selected[fullPath] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m *model) View() string {
	var s strings.Builder

	s.WriteString("Select files for context (space: toggle, enter: open, backspace: up, q: save & quit)\n")
	s.WriteString("Current path: " + m.path + "\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		fullPath := filepath.Join(m.path, item.name)
		prefix := "[ ]"
		if _, ok := m.selected[fullPath]; ok {
			prefix = "[x]"
		}

		itemName := item.name
		if item.isDir {
			itemName += "/"
		}

		s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, prefix, itemName))
	}

	s.WriteString(fmt.Sprintf("\nSelected %d files. Press 'q' to save and exit.", len(m.selected)))

	return s.String()
}
