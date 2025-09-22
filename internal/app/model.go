package app

import (
	"fmt"
	"getctx/internal/logger"
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		inputWidth := m.width - len(m.pathInput.Prompt) - 1
		inputWidth = max(inputWidth, 1)

		m.pathInput.Width = inputWidth
	}

	if m.isInputMode {
		m.pathInput, cmd = m.pathInput.Update(msg)
		cmds = append(cmds, cmd)
		cmd = m.updateInputMode(msg)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		cmd = m.updateNormalMode(msg)
	}
	cmds = append(cmds, cmd)

	headerContent := m.renderHeader()
	footerContent := m.renderFooter()

	headerHeight := lipgloss.Height(headerContent)
	footerHeight := lipgloss.Height(footerContent)

	m.viewport.Width = m.width
	m.viewport.Height = m.height - headerHeight - footerHeight

	m.viewport.SetContent(m.renderFileList())
	m.ensureCursorVisible()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.handleConfirmPathChange()
			return nil
		case KeyEscape, KeyCtrlC:
			m.handleCancelPathChange()
			return nil
		}
	}
	return nil
}

func (m *Model) updateNormalMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyQ:
			return m.handleConfirmAndExit()
		case KeyCtrlC:
			return m.handleCancelAndExit()
		case KeyUp:
			m.handleMoveCursorUp()
		case KeyDown:
			m.handleMoveCursorDown()
		case KeyEnter:
			m.handleEnterDirectory()
		case KeyBackspace:
			m.handleNavigateToParent()
		case KeySpace:
			m.handleSelectFile()
		case KeyCtrlA:
			m.handleSelectAllFiles()
		case KeyCtrlHome:
			m.handleGoToTop()
		case KeyCtrlEnd:
			m.handleGoToBottom()
		case KeyP:
			return m.handleEnterPathInputMode()
		}
	}
	return nil
}

func (m *Model) View() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}

func (m *Model) renderPathInput() string {
	var s strings.Builder
	s.WriteString(Elements.Text.InputHeader)
	s.WriteString(m.pathInput.View())

	if m.inputErrorMsg != "" {
		s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
	}
	s.WriteString("\n")
	return s.String()
}

func (m *Model) renderHeader() string {
	if m.isInputMode {
		return m.renderPathInput()
	}

	helpHeader := Elements.Text.HelpHeader
	pathStyle := lipgloss.NewStyle().Width(m.width)
	fullPathString := Elements.Text.PathPrefix + m.path
	wrappedPath := pathStyle.Render(fullPathString)

	return lipgloss.JoinVertical(lipgloss.Left,
		helpHeader,
		wrappedPath,
	)
}

func (m *Model) renderFooter() string {
	return fmt.Sprintf(Elements.Text.StatusFooter, len(m.selected))
}

func (m *Model) renderFileList() string {

	if len(m.items) == 0 {
		return m.renderEmptyFileList()
	}

	var s strings.Builder
	for i, item := range m.items {
		s.WriteString(m.renderListItem(i, item))
	}
	return s.String()
}

func (m *Model) renderEmptyFileList() string {
	emptyMessage := Elements.Text.EmptyMessage
	padding := strings.Repeat("\n", m.viewport.Height/2)

	style := Styles.List.Empty.
		Width(m.viewport.Width).
		Align(lipgloss.Center)

	return style.Render(padding + emptyMessage)
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

	finalCursor := cursorStyle.Render(cursorStr)
	finalPrefix := nameStyle.Render(Elements.List.CursorEmpty + prefix + itemIcon + Elements.List.CursorEmpty)
	finalName := nameStyle.Render(itemName)

	line := lipgloss.JoinHorizontal(lipgloss.Top,
		finalCursor,
		finalPrefix,
		finalName,
	)

	return line + "\n"
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
		logger.Error("changeDirectory.loadItems", map[string]any{
			"message": "Failed to load directory items",
			"path":    newPath,
			"error":   err.Error(),
		})
		return
	}

	logger.Debug("changeDirectory", map[string]any{
		"path":       newPath,
		"item_count": len(newItems),
	})

	m.path = newPath
	m.items = newItems
	m.cursor = 0
	m.viewport.GotoTop()
}

func (m *Model) handleConfirmAndExit() tea.Cmd {
	return tea.Quit
}

func (m *Model) handleCancelAndExit() tea.Cmd {
	m.selected = make(map[string]struct{})
	return tea.Quit
}

func (m *Model) handleConfirmPathChange() {
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
		logger.Warn("handleConfirmPathChange.Abs", map[string]any{
			"message": "Could not get absolute path",
			"path":    newPath,
			"error":   err.Error(),
		})
		return
	}

	if _, err := loadItems(m.fsys, absPath, m.config); err != nil {
		m.inputErrorMsg = fmt.Sprintf("Error reading directory: %v", err)
		logger.Error("handleConfirmPathChange.loadItems", map[string]any{
			"message": "Failed to load directory items on path confirm",
			"path":    absPath,
			"error":   err.Error(),
		})
	} else {
		logger.Info("handleConfirmPathChange", map[string]any{
			"message": "Directory changed successfully via path input",
			"path":    absPath,
		})
		m.changeDirectory(absPath)
		m.isInputMode = false
		m.inputErrorMsg = ""
		m.pathInput.Reset()
	}
}

func (m *Model) handleCancelPathChange() {
	m.isInputMode = false
	m.inputErrorMsg = ""
	m.pathInput.Reset()
}

func (m *Model) handleMoveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) handleMoveCursorDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
}

func (m *Model) handleEnterDirectory() {
	if len(m.items) == 0 {
		return
	}
	currentItem := m.items[m.cursor]
	if currentItem.isDir && !currentItem.isExcluded {
		m.changeDirectory(filepath.Join(m.path, currentItem.name))
	}
}

func (m *Model) handleNavigateToParent() {
	parentPath := filepath.Dir(m.path)
	if parentPath != m.path {
		m.changeDirectory(parentPath)
	}
}

func (m *Model) handleSelectFile() {
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

func (m *Model) handleSelectAllFiles() {
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

func (m *Model) handleGoToTop() {
	if len(m.items) > 0 {
		m.cursor = 0
	}
}

func (m *Model) handleGoToBottom() {
	if len(m.items) > 0 {
		m.cursor = len(m.items) - 1
	}
}

func (m *Model) handleEnterPathInputMode() tea.Cmd {
	m.isInputMode = true
	m.inputErrorMsg = ""
	m.pathInput.SetValue(m.path + string(filepath.Separator))
	return m.pathInput.Focus()
}

func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	}
	if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
	}
}
