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

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Could not read directory '%s': %v", path, err)
	}

	var items []item
	for _, entry := range dirEntries {
		items = append(items, item{name: entry.Name(), isDir: entry.IsDir()})
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
		case "ctrl+c":
			m.selected = make(map[string]struct{})
			return m, tea.Quit

		case "q":
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
				dirEntries, err := os.ReadDir(newPath)
				if err != nil {
					log.Printf("Error reading directory %s: %v", newPath, err)
					break
				}
				var newItems []item
				for _, entry := range dirEntries {
					newItems = append(newItems, item{name: entry.Name(), isDir: entry.IsDir()})
				}
				m.path = newPath
				m.items = newItems
				m.cursor = 0
			}
		case "backspace":
			parentPath := filepath.Dir(m.path)
			if parentPath != m.path {
				dirEntries, err := os.ReadDir(parentPath)
				if err != nil {
					log.Printf("Error reading directory %s: %v", parentPath, err)
					break
				}
				var newItems []item
				for _, entry := range dirEntries {
					newItems = append(newItems, item{name: entry.Name(), isDir: entry.IsDir()})
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
			cursor = Icons.Cursor
		}

		fullPath := filepath.Join(m.path, item.name)
		_, ok := m.selected[fullPath]

		prefix := "  "
		if ok {
			prefix = Icons.Checkmark + " "
		}

		itemIcon := Icons.File
		if item.isDir {
			itemIcon = Icons.Directory
		}

		itemName := item.name
		if item.isDir {
			itemName += "/"
		}
		line := fmt.Sprintf("%s %s%s %s", cursor, prefix, itemIcon, itemName)

		if ok {
			s.WriteString(Styles.Selected.Render(line))
		} else if m.cursor == i {
			s.WriteString(Styles.Cursor.Render(line))
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}
	s.WriteString(fmt.Sprintf("\nSelected %d files. Press 'q' to save and exit.", len(m.selected)))
	return s.String()
}
