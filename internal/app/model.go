package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name       string
	isDir      bool
	isExcluded bool
}

type Model struct {
	path          string
	items         []item
	cursor        int
	selected      map[string]struct{}
	config        *Config
	pathInput     textinput.Model
	isInputMode   bool
	inputErrorMsg string
}

func loadItems(path string, config *Config) ([]item, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []item
	for _, entry := range dirEntries {
		items = append(items, item{
			name:       entry.Name(),
			isDir:      entry.IsDir(),
			isExcluded: config.IsExcluded(entry.Name()),
		})
	}
	return items, nil
}

func NewModel(startPath string, config *Config) (*Model, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", startPath, err)
	}

	items, err := loadItems(path, config)
	if err != nil {
		return nil, fmt.Errorf("could not read directory '%s': %w", path, err)
	}

	ti := textinput.New()
	ti.Placeholder = "/home/user/project..."
	ti.Prompt = Icons.Cursor + Elements.List.CursorEmpty
	ti.Focus()

	return &Model{
		path:        path,
		items:       items,
		selected:    make(map[string]struct{}),
		config:      config,
		pathInput:   ti,
		isInputMode: false,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.isInputMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case KeyEnter:
				newPath := m.pathInput.Value()
				if strings.HasPrefix(newPath, "~") {
					home, err := os.UserHomeDir()
					if err == nil {
						newPath = filepath.Join(home, newPath[1:])
					}
				}
				absPath, err := filepath.Abs(newPath)
				if err != nil {
					m.inputErrorMsg = fmt.Sprintf("Invalid path: %v", err)
					return m, nil
				}
				newItems, err := loadItems(absPath, m.config)
				if err != nil {
					m.inputErrorMsg = fmt.Sprintf("Error reading directory: %v", err)
				} else {
					m.path = absPath
					m.items = newItems
					m.cursor = 0
					m.isInputMode = false
					m.inputErrorMsg = ""
					m.pathInput.Reset()
				}
				return m, nil

			case KeyEscape, KeyCtrlC:
				m.isInputMode = false
				m.inputErrorMsg = ""
				m.pathInput.Reset()
				return m, nil
			}
		}
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.items) == 0 && msg.String() != KeyQ && msg.String() != KeyCtrlC {
			return m, nil
		}

		switch msg.String() {
		case KeyCtrlC, KeyQ:
			if msg.String() == KeyCtrlC {
				m.selected = make(map[string]struct{})
			}
			return m, tea.Quit
		case KeyUp:
			m.handleKeyUp()
		case KeyDown:
			m.handleKeyDown()
		case KeyEnter:
			m.handleEnter()
		case KeyBackspace:
			m.handleBackspace()
		case KeySpace:
			m.handleSpace()
		case KeyCtrlA:
			m.handleCtrlA()
		case KeyCtrlHome:
			m.handleCtrlHome()
		case KeyCtrlEnd:
			m.handleCtrlEnd()
		case KeyP:
			m.isInputMode = true
			m.inputErrorMsg = ""
			m.pathInput.SetValue(m.path + string(filepath.Separator))
			return m, m.pathInput.Focus()
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var s strings.Builder

	if m.isInputMode {
		s.WriteString("Enter path (Enter to confirm, Esc to cancel):\n")
		s.WriteString(m.pathInput.View())
		if m.inputErrorMsg != "" {
			s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
		}
		s.WriteString("\n\n")
	} else {
		s.WriteString(Elements.Text.HelpHeader)
	}

	s.WriteString(Elements.Text.PathPrefix + m.path + "\n\n")

	for i, item := range m.items {
		cursor := Elements.List.CursorEmpty
		if m.cursor == i {
			cursor = Icons.Cursor
		}

		var prefix, itemIcon string

		if item.isExcluded {
			prefix = Elements.List.UnselectedPrefix
			itemIcon = Icons.Excluded
		} else {
			fullPath := filepath.Join(m.path, item.name)
			_, isSelected := m.selected[fullPath]

			if isSelected {
				prefix = Elements.List.SelectedPrefix
			} else {
				prefix = Elements.List.UnselectedPrefix
			}

			if item.isDir {
				itemIcon = Icons.Directory
			} else {
				itemIcon = Icons.File
			}
		}

		itemName := item.name
		if item.isDir {
			itemName += Elements.List.DirectorySuffix
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
	s.WriteString(fmt.Sprintf(Elements.Text.StatusFooter, len(m.selected)))
	return s.String()
}

func (m *Model) handleKeyUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) handleKeyDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
}

func (m *Model) handleEnter() {
	if len(m.items) == 0 {
		return
	}
	currentItem := m.items[m.cursor]
	if currentItem.isDir && !currentItem.isExcluded {
		newPath := filepath.Join(m.path, currentItem.name)
		newItems, err := loadItems(newPath, m.config)
		if err != nil {
			log.Printf("Error reading directory %s: %v", newPath, err)
			return
		}
		m.path = newPath
		m.items = newItems
		m.cursor = 0
	}
}

func (m *Model) handleBackspace() {
	parentPath := filepath.Dir(m.path)
	if parentPath != m.path {
		newItems, err := loadItems(parentPath, m.config)
		if err != nil {
			log.Printf("Error reading directory %s: %v", parentPath, err)
			return
		}
		m.path = parentPath
		m.items = newItems
		m.cursor = 0
	}
}

func (m *Model) handleSpace() {
	if len(m.items) == 0 {
		return
	}
	currentItem := m.items[m.cursor]
	if !currentItem.isExcluded {
		fullPath := filepath.Join(m.path, currentItem.name)
		if _, ok := m.selected[fullPath]; ok {
			delete(m.selected, fullPath)
		} else {
			m.selected[fullPath] = struct{}{}
		}
	}
}

func (m *Model) handleCtrlA() {
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
		return
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

func (m *Model) handleCtrlHome() {
	if len(m.items) > 0 {
		m.cursor = 0
	}
}

func (m *Model) handleCtrlEnd() {
	if len(m.items) > 0 {
		m.cursor = len(m.items) - 1
	}
}
