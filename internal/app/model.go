package app

import (
	"fmt"
	"log"
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
	fsys          FileSystem
	pathInput     textinput.Model
	isInputMode   bool
	inputErrorMsg string
}

func NewModel(startPath string, config *Config, fsys FileSystem) (*Model, error) {
	path, err := fsys.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", startPath, err)
	}

	items, err := loadItems(fsys, path, config)
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
		fsys:        fsys,
		pathInput:   ti,
		isInputMode: false,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isInputMode {
		return m.updateInputMode(msg)
	}
	return m.updateNormalMode(msg)
}

func (m *Model) View() string {
	if m.isInputMode {
		return m.viewInputMode()
	}
	return m.viewNormalMode()
}

func (m *Model) updateInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.handlePathInputConfirm()
			return m, nil
		case KeyEscape, KeyCtrlC:
			m.handlePathInputCancel()
			return m, nil
		}
	}

	m.pathInput, cmd = m.pathInput.Update(msg)
	return m, cmd
}

func (m *Model) viewInputMode() string {
	var s strings.Builder
	s.WriteString("Enter path (Enter to confirm, Esc to cancel):\n")
	s.WriteString(m.pathInput.View())
	if m.inputErrorMsg != "" {
		s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
	}
	s.WriteString("\n\n")
	return s.String()
}

func (m *Model) updateNormalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.items) == 0 && !isQuitKey(msg.String()) {
			return m, nil
		}

		switch msg.String() {
		case KeyCtrlC, KeyQ:
			return m.handleQuit(msg)
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
			cmd := m.handleGoToPath()
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) viewNormalMode() string {
	var s strings.Builder
	s.WriteString(Elements.Text.HelpHeader)
	s.WriteString(Elements.Text.PathPrefix + m.path + "\n\n")
	s.WriteString(m.renderFileList())
	s.WriteString(fmt.Sprintf(Elements.Text.StatusFooter, len(m.selected)))
	return s.String()
}

func (m *Model) renderFileList() string {
	var s strings.Builder
	for i, item := range m.items {
		s.WriteString(m.renderListItem(i, item))
		s.WriteString("\n")
	}
	return s.String()
}
func (m *Model) renderListItem(index int, item item) string {
	cursor := Elements.List.CursorEmpty
	if m.cursor == index {
		cursor = Icons.Cursor
	}

	fullPath := filepath.Join(m.path, item.name)
	_, isSelected := m.selected[fullPath]

	var prefix, itemIcon string
	if item.isExcluded {
		prefix = Elements.List.UnselectedPrefix
		itemIcon = Icons.Excluded
	} else {
		if isSelected {
			prefix = Elements.List.SelectedPrefix
		} else {
			prefix = Elements.List.UnselectedPrefix
		}
		itemIcon = Icons.File
		if item.isDir {
			itemIcon = Icons.Directory
		}
	}

	itemName := item.name
	if item.isDir {
		itemName += Elements.List.DirectorySuffix
	}
	line := fmt.Sprintf("%s %s%s %s", cursor, prefix, itemIcon, itemName)

	if item.isExcluded {
		return Styles.List.Excluded.Render(line)
	}
	if isSelected {
		return Styles.List.Selected.Render(line)
	}
	if m.cursor == index {
		return Styles.List.Cursor.Render(line)
	}
	return line
}

func isQuitKey(key string) bool {
	return key == KeyQ || key == KeyCtrlC
}

func loadItems(fsys FileSystem, path string, config *Config) ([]item, error) {
	dirEntries, err := fsys.ReadDir(path)
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
func (m *Model) changeDirectory(newPath string) {
	newItems, err := loadItems(m.fsys, newPath, m.config)
	if err != nil {
		log.Printf("Error reading directory %s: %v", newPath, err)
		return
	}
	m.path = newPath
	m.items = newItems
	m.cursor = 0
}

func (m *Model) handleQuit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == KeyCtrlC {
		m.selected = make(map[string]struct{})
	}
	return m, tea.Quit
}

func (m *Model) handleGoToPath() tea.Cmd {
	m.isInputMode = true
	m.inputErrorMsg = ""
	m.pathInput.SetValue(m.path + string(filepath.Separator))
	return m.pathInput.Focus()
}

func (m *Model) handlePathInputConfirm() {
	newPath := m.pathInput.Value()
	if strings.HasPrefix(newPath, "~") {
		home, err := m.fsys.UserHomeDir()
		if err == nil {
			newPath = filepath.Join(home, newPath[1:])
		}
	}

	absPath, err := m.fsys.Abs(newPath)
	if err != nil {
		m.inputErrorMsg = fmt.Sprintf("Invalid path: %v", err)
		return
	}

	if _, err := loadItems(m.fsys, absPath, m.config); err != nil {
		m.inputErrorMsg = fmt.Sprintf("Error reading directory: %v", err)
	} else {
		m.changeDirectory(absPath)
		m.isInputMode = false
		m.inputErrorMsg = ""
		m.pathInput.Reset()
	}
}

func (m *Model) handlePathInputCancel() {
	m.isInputMode = false
	m.inputErrorMsg = ""
	m.pathInput.Reset()
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
		m.changeDirectory(filepath.Join(m.path, currentItem.name))
	}
}

func (m *Model) handleBackspace() {
	parentPath := filepath.Dir(m.path)
	if parentPath != m.path {
		m.changeDirectory(parentPath)
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
				delete(m.selected, filepath.Join(m.path, item.name))
			}
		}
	} else {
		for _, item := range m.items {
			if !item.isExcluded {
				m.selected[filepath.Join(m.path, item.name)] = struct{}{}
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
