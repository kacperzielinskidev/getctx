package tui

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputWidth := max(m.width-len(m.pathInput.Prompt)-1, 1)
		m.pathInput.Width = inputWidth
	}

	if m.isInputMode {
		cmd = m.updateInputMode(msg)
		cmds = append(cmds, cmd)
	} else if m.isFilterMode {
		cmd = m.updateFilterMode(msg)
		cmds = append(cmds, cmd)
	} else {
		cmd = m.updateNormalMode(msg)
		cmds = append(cmds, cmd)
	}

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
	oldValue := m.pathInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyTab:
			m.handleAutoComplete()
			return nil
		case KeyEnter:
			m.handleConfirmPathChange()
			return nil
		case KeyEscape, KeyCtrlC:
			m.handleCancelPathChange()
			return nil
		}
	}

	var cmd tea.Cmd
	m.pathInput, cmd = m.pathInput.Update(msg)

	if m.pathInput.Value() != oldValue {
		m.updateCompletions()
	}

	return cmd
}

func (m *Model) updateFilterMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.isFilterMode = false
			return nil
		case KeyEscape, KeyCtrlC:
			m.handleCancelFilter()
			return nil
		}
	}
	var cmd tea.Cmd
	m.pathInput, cmd = m.pathInput.Update(msg)
	m.filterQuery = m.pathInput.Value()
	m.clampCursor()
	return cmd
}

func (m *Model) updateNormalMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEscape:
			m.handleClearFilter()
		case KeyCtrlC:
			return m.handleCancelAndExit()
		case KeyQ:
			return m.handleConfirmAndExit()
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
		case KeySlash:
			return m.handleEnterFilterMode()
		}
	}
	return nil
}

func (m *Model) handleConfirmAndExit() tea.Cmd {
	return tea.Quit
}

func (m *Model) handleCancelAndExit() tea.Cmd {
	m.Aborted = true
	return tea.Quit
}

func (m *Model) handleClearFilter() {
	if m.filterQuery != "" {
		m.filterQuery = ""
		m.clampCursor()
	}
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
		m.inputErrorMsg = "Invalid path: " + err.Error()
		return
	}

	m.changeDirectory(absPath)
	m.isInputMode = false
	m.inputErrorMsg = ""
	m.pathInput.Reset()
	m.completionSuggestions = nil
}

func (m *Model) handleCancelPathChange() {
	m.isInputMode = false
	m.inputErrorMsg = ""
	m.pathInput.Reset()
	m.completionSuggestions = nil
}

func (m *Model) handleMoveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) handleMoveCursorDown() {
	visibleItems := m.getVisibleItems()
	if m.cursor < len(visibleItems)-1 {
		m.cursor++
	}
}

func (m *Model) handleEnterDirectory() {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		return
	}
	currentItem := visibleItems[m.cursor]
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
	visibleItems := m.getVisibleItems()
	if len(visibleItems) > 0 {
		currentItem := visibleItems[m.cursor]
		if !currentItem.isExcluded {
			fullPath := filepath.Join(m.path, currentItem.name)
			if _, ok := m.selected[fullPath]; ok {
				delete(m.selected, fullPath)
			} else {
				m.selected[fullPath] = struct{}{}
			}
		}
	}
}

func (m *Model) handleSelectAllFiles() {
	visibleItems := m.getVisibleItems()
	allSelected := true
	for _, item := range visibleItems {
		if !item.isExcluded {
			fullPath := filepath.Join(m.path, item.name)
			if _, ok := m.selected[fullPath]; !ok {
				allSelected = false
				break
			}
		}
	}

	for _, item := range visibleItems {
		if !item.isExcluded {
			fullPath := filepath.Join(m.path, item.name)
			if allSelected {
				delete(m.selected, fullPath)
			} else {
				m.selected[fullPath] = struct{}{}
			}
		}
	}
}

func (m *Model) handleGoToTop() {
	m.cursor = 0
}

func (m *Model) handleGoToBottom() {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) > 0 {
		m.cursor = len(visibleItems) - 1
	}
}

func (m *Model) handleEnterPathInputMode() tea.Cmd {
	m.isInputMode = true
	m.inputErrorMsg = ""
	pathValue := m.path + string(filepath.Separator)
	m.pathInput.SetValue(pathValue)
	m.pathInput.SetCursor(len(pathValue))
	m.updateCompletions()
	return m.pathInput.Focus()
}

func (m *Model) handleEnterFilterMode() tea.Cmd {
	m.isFilterMode = true
	m.pathInput.SetValue(m.filterQuery)
	return m.pathInput.Focus()
}

func (m *Model) handleCancelFilter() {
	m.isFilterMode = false
	m.filterQuery = ""
	m.pathInput.Reset()
}

func (m *Model) updateCompletions() {
	currentInput := m.pathInput.Value()
	if currentInput == "" {
		m.completionSuggestions = nil
		return
	}
	dirToSearch, prefix := m.getCompletionParts(currentInput)
	suggestions, err := m.getCompletions(dirToSearch, prefix)
	if err != nil {
		m.completionSuggestions = nil
		return
	}
	m.completionSuggestions = suggestions
}

func (m *Model) handleAutoComplete() {
	if len(m.completionSuggestions) == 0 {
		return
	}
	suggestions := m.completionSuggestions
	dirToSearch, _ := m.getCompletionParts(m.pathInput.Value())
	var newInputValue string
	if len(suggestions) == 1 {
		newInputValue = filepath.Join(dirToSearch, suggestions[0])
	} else {
		commonPrefix := findLongestCommonPrefix(suggestions)
		if commonPrefix == "" {
			return
		}
		newInputValue = filepath.Join(dirToSearch, commonPrefix)
	}

	info, err := m.fsys.Stat(newInputValue)
	if err == nil && info.IsDir() && !strings.HasSuffix(newInputValue, string(filepath.Separator)) {
		newInputValue += string(filepath.Separator)
	}

	m.pathInput.SetValue(newInputValue)
	m.pathInput.SetCursor(len(newInputValue))
	m.updateCompletions()
}
