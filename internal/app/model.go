// Plik: internal/app/model.go
package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	viewport      viewport.Model
	width         int
	height        int
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

	vp := viewport.New(0, 0)

	m := &Model{
		path:        path,
		items:       items,
		selected:    make(map[string]struct{}),
		config:      config,
		fsys:        fsys,
		pathInput:   ti,
		isInputMode: false,
		viewport:    vp,
	}

	m.viewport.SetContent(m.renderFileList())
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.pathInput, cmd = m.pathInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		var header strings.Builder
		if m.isInputMode {
			header.WriteString(m.renderPathInput())
		} else {
			header.WriteString(Elements.Text.HelpHeader)
		}
		header.WriteString(Elements.Text.PathPrefix + m.path + "\n\n")

		footer := fmt.Sprintf(Elements.Text.StatusFooter, len(m.selected))

		headerHeight := lipgloss.Height(header.String())
		footerHeight := lipgloss.Height(footer)

		m.viewport.Width = m.width
		m.viewport.Height = m.height - headerHeight - footerHeight
	}

	if m.isInputMode {
		cmd = m.updateInputMode(msg)
	} else {
		cmd = m.updateNormalMode(msg)
	}
	cmds = append(cmds, cmd)

	m.viewport.SetContent(m.renderFileList())
	m.ensureCursorVisible()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.handlePathInputConfirm()
			return nil
		case KeyEscape, KeyCtrlC:
			m.handlePathInputCancel()
			return nil
		}
	}
	return nil
}

func (m *Model) updateNormalMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.items) == 0 && !isQuitKey(msg.String()) {
			return nil
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
			return m.handleGoToPath()
		}
	}
	return nil
}

func (m *Model) View() string {
	var header strings.Builder
	if m.isInputMode {
		header.WriteString(m.renderPathInput())
	} else {
		header.WriteString(Elements.Text.HelpHeader)
	}
	header.WriteString(Elements.Text.PathPrefix + m.path + "\n\n")

	footer := fmt.Sprintf(Elements.Text.StatusFooter, len(m.selected))

	return header.String() + m.viewport.View() + footer
}

func (m *Model) renderPathInput() string {
	var s strings.Builder

	s.WriteString(Elements.Text.InputHeader)

	s.WriteString(m.pathInput.View())
	if m.inputErrorMsg != "" {
		s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
	}
	s.WriteString("\n\n")
	return s.String()
}

func (m *Model) renderFileList() string {
	var s strings.Builder
	for i, item := range m.items {
		s.WriteString(m.renderListItem(i, item))
	}
	return s.String()
}

func (m *Model) renderListItem(index int, item item) string {
	var cursorStyle, nameStyle lipgloss.Style

	if m.cursor == index {
		cursorStyle = Styles.List.Cursor
		nameStyle = Styles.List.Cursor
	} else {
		cursorStyle = Styles.List.Normal
		nameStyle = Styles.List.Normal
	}

	fullPath := filepath.Join(m.path, item.name)
	_, isSelected := m.selected[fullPath]

	if item.isExcluded {
		nameStyle = Styles.List.Excluded
	} else if isSelected {
		nameStyle = Styles.List.Selected
	}

	cursorStr := Elements.List.CursorEmpty
	if m.cursor == index {
		cursorStr = Icons.Cursor
	}

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

	// Tworzymy stringi dla każdej części
	finalCursor := cursorStyle.Render(cursorStr)
	finalPrefix := nameStyle.Render(Elements.List.CursorEmpty + prefix + itemIcon + Elements.List.CursorEmpty)
	finalName := nameStyle.Render(itemName)

	// Używamy JoinHorizontal do połączenia komponentów.
	// Lipgloss sam zadba o prawidłowe obliczenie szerokości.
	line := lipgloss.JoinHorizontal(lipgloss.Top,
		finalCursor,
		finalPrefix,
		finalName,
	)

	return line + "\n"
}

// ... reszta kodu jest poprawna i pozostaje bez zmian ...
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
	m.viewport.GotoTop()
}

func (m *Model) handleQuit(msg tea.KeyMsg) tea.Cmd {
	if msg.String() == KeyCtrlC {
		m.selected = make(map[string]struct{})
	}
	return tea.Quit
}

func (m *Model) handleGoToPath() tea.Cmd {
	log.Println("HANDLER: handleGoToPath called!")
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

	// POPRAWKA: Używamy m.fsys zamiast fsys
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

func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	}
	if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
	}
}
