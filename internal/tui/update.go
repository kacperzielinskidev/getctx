package tui

import (
	"path/filepath"
	"runtime"
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
		m.textInput.Width = max(m.width-len(m.textInput.Prompt)-1, 1)
	}

	switch m.mode {
	case modeNormal:
		cmd = m.updateNormalMode(msg)
	case modePathInput:
		cmd = m.updatePathInputMode(msg)
	case modeFilter:
		cmd = m.updateFilterMode(msg)
	}
	cmds = append(cmds, cmd)

	headerContent := m.renderHeader()
	footerContent := m.renderFooter()
	headerHeight := lipgloss.Height(headerContent)
	footerHeight := lipgloss.Height(footerContent)
	viewportHeight := m.height - headerHeight - footerHeight

	m.viewport.Height = viewportHeight
	m.completionViewport.Height = viewportHeight

	m.ensureCursorVisible()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateNormalMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyUp:
			m.handleMoveCursorUp()
		case KeyDown:
			m.handleMoveCursorDown()
		case KeyCtrlHome:
			m.handleGoToTop()
		case KeyCtrlEnd:
			m.handleGoToBottom()
		case KeyEnter:
			m.enterDirectory()
		case KeyBackspace:
			m.navigateToParent()
		case KeySpace:
			m.toggleSelection()
		case KeyCtrlA:
			m.toggleSelectAll()
		case KeySlash:
			return m.enterFilterMode()
		case KeyP:
			return m.enterPathInputMode()
		case KeyEscape:
			m.clearFilter()
		case KeyQ:
			return tea.Quit
		case KeyCtrlC:
			m.Aborted = true
			return tea.Quit
		}
	}
	return nil
}

func (m *Model) updatePathInputMode(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	oldValue := m.textInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.confirmPathChange()
			return nil
		case KeyEscape, KeyCtrlC:
			m.cancelInputMode()
			return nil
		case KeyTab:
			m.autoCompletePath()
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	default:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.textInput.Value() != oldValue {
		m.updateCompletions()
	}

	m.completionViewport, cmd = m.completionViewport.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) updateFilterMode(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			m.mode = modeNormal
			m.textInput.Blur()
			m.clampCursor()
			return nil
		case KeyEscape, KeyCtrlC:
			m.clearFilter()
			return nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.filterQuery = m.textInput.Value()
	m.clampCursor()
	return cmd
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

func (m *Model) handleGoToTop() {
	m.cursor = 0
}

func (m *Model) handleGoToBottom() {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) > 0 {
		m.cursor = len(visibleItems) - 1
	}
}

func (m *Model) enterDirectory() {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		return
	}
	currentItem := visibleItems[m.cursor]
	if currentItem.isDir && !currentItem.isExcluded {
		m.changeDirectory(filepath.Join(m.path, currentItem.name))
	}
}

func (m *Model) navigateToParent() {
	parentPath := filepath.Dir(m.path)
	if parentPath != m.path {
		m.changeDirectory(parentPath)
	}
}

func (m *Model) toggleSelection() {
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

func (m *Model) toggleSelectAll() {
	visibleItems := m.getVisibleItems()
	allSelected := true
	if len(visibleItems) == 0 {
		allSelected = false
	}
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

func (m *Model) confirmPathChange() {
	inputPath := m.textInput.Value()

	if strings.HasPrefix(inputPath, "~") {
		home, err := m.fsys.UserHomeDir()
		if err == nil {
			inputPath = filepath.Join(home, inputPath[1:])
		}
	}

	var finalPath string
	if filepath.IsAbs(inputPath) {
		finalPath = inputPath
	} else if runtime.GOOS == "windows" && strings.HasPrefix(inputPath, `\`) {
		finalPath = filepath.VolumeName(m.path) + inputPath
	} else {
		finalPath = filepath.Join(m.path, inputPath)
	}

	cleanedPath := filepath.Clean(finalPath)
	info, err := m.fsys.Stat(cleanedPath)
	if err != nil {
		m.inputErrorMsg = "Error: Path not found or is inaccessible."
		return
	}
	if !info.IsDir() {
		m.inputErrorMsg = "Error: Path is a file, not a directory."
		return
	}

	m.changeDirectory(cleanedPath)
	m.cancelInputMode()
}

func (m *Model) enterPathInputMode() tea.Cmd {
	m.mode = modePathInput
	m.inputErrorMsg = ""
	pathValue := m.path + string(filepath.Separator)
	m.textInput.SetValue(pathValue)
	m.textInput.SetCursor(len(pathValue))
	m.updateCompletions()
	return m.textInput.Focus()
}

func (m *Model) enterFilterMode() tea.Cmd {
	m.mode = modeFilter
	m.textInput.SetValue(m.filterQuery)
	m.textInput.SetCursor(len(m.filterQuery))
	return m.textInput.Focus()
}

func (m *Model) cancelInputMode() {
	m.mode = modeNormal
	m.inputErrorMsg = ""
	m.textInput.Blur()
	m.textInput.Reset()
	m.completionSuggestions = nil
}

func (m *Model) clearFilter() {
	if m.filterQuery != "" {
		m.filterQuery = ""
		m.textInput.Reset()
		m.clampCursor()
	}
	if m.mode == modeFilter {
		m.mode = modeNormal
		m.textInput.Blur()
	}
}

func (m *Model) updateCompletions() {
	currentInput := m.textInput.Value()
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

func (m *Model) autoCompletePath() {
	if len(m.completionSuggestions) == 0 {
		return
	}
	suggestions := m.completionSuggestions
	dirToSearch, _ := m.getCompletionParts(m.textInput.Value())

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

	m.textInput.SetValue(newInputValue)
	m.textInput.SetCursor(len(newInputValue))
	m.updateCompletions()
}
