package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name       string
	isDir      bool
	isExcluded bool
}

type Model struct {
	path     string
	items    []item
	cursor   int
	selected map[string]struct{}
}

func NewModel(startPath string) *Model {
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
		items = append(items, item{name: entry.Name(), isDir: entry.IsDir(), isExcluded: isExcluded(entry.Name())})
	}

	return &Model{
		path:     path,
		items:    items,
		selected: make(map[string]struct{}),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.items) == 0 && msg.String() != KeyQ && msg.String() != KeyCtrlC {
			return m, nil
		}

		currentItem := m.items[m.cursor]

		switch msg.String() {
		case KeyCtrlC:
			m.selected = make(map[string]struct{})
			return m, tea.Quit

		case KeyQ:
			return m, tea.Quit

		case KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case KeyDown:
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case KeyEnter:
			if currentItem.isDir && !currentItem.isExcluded {

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
			}

		case KeyBackspace:
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

		case KeySpace:
			if !currentItem.isExcluded {
				fullPath := filepath.Join(m.path, currentItem.name)
				if _, ok := m.selected[fullPath]; ok {
					delete(m.selected, fullPath)
				} else {
					m.selected[fullPath] = struct{}{}
				}
			}

		case KeyCtrlA:
			allSelectableAreSelected := true
			hasSelectableItems := false
			for _, item := range m.items {
				if !item.isExcluded {
					hasSelectableItems = true
					fullPath := filepath.Join(m.path, item.name)
					if _, ok := m.selected[fullPath]; !ok {
						allSelectableAreSelected = false
						break
					}
				}
			}

			if !hasSelectableItems {
				break
			}

			if allSelectableAreSelected {
				for _, item := range m.items {
					if !item.isExcluded {
						fullPath := filepath.Join(m.path, item.name)
						delete(m.selected, fullPath)
					}
				}
			} else {
				for _, item := range m.items {
					if !item.isExcluded {
						fullPath := filepath.Join(m.path, item.name)
						m.selected[fullPath] = struct{}{}
					}
				}
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var s strings.Builder
	s.WriteString("Select files for context (space: toggle, enter: open, backspace: up, q: save & quit)\n")
	s.WriteString("Current path: " + m.path + "\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = Icons.Cursor
		}

		var prefix, itemIcon string

		if item.isExcluded {
			prefix = "  "
			itemIcon = Icons.Excluded
		} else {
			fullPath := filepath.Join(m.path, item.name)
			_, isSelected := m.selected[fullPath]

			if isSelected {
				prefix = Icons.Checkmark + " "
			} else {
				prefix = "  "
			}

			if item.isDir {
				itemIcon = Icons.Directory
			} else {
				itemIcon = Icons.File
			}
		}

		itemName := item.name
		if item.isDir {
			itemName += "/"
		}
		line := fmt.Sprintf("%s %s%s %s", cursor, prefix, itemIcon, itemName)

		fullPath := filepath.Join(m.path, item.name)
		_, isSelected := m.selected[fullPath]

		if item.isExcluded {
			s.WriteString(Styles.List.Excluded.Render(line))
		} else if isSelected {
			s.WriteString(Styles.List.Selected.Render(line))
		} else if m.cursor == i {
			s.WriteString(Styles.List.Cursor.Render(line))
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}
	s.WriteString(fmt.Sprintf("\nSelected %d items. Press 'q' to save and exit.", len(m.selected)))
	return s.String()
}
