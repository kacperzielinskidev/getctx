package app

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
	m.filterQuery = m.pathInput.Value()
	m.clampCursor()
	return nil
}

func (m *Model) updateNormalMode(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEscape:
			m.handleClearFilter()
			return nil
		case KeyCtrlC:
			if m.filterQuery != "" {
				m.handleClearFilter()
				return nil
			}
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
	m.selected = make(map[string]struct{})
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

	if _, err := loadItems(m.fsys, absPath, m.config); err != nil {
		m.inputErrorMsg = "Error reading directory: " + err.Error()
	} else {
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
	if len(visibleItems) == 0 {
		return
	}
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

func (m *Model) handleSelectAllFiles() {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		return
	}

	allSelectableAreSelected := true
	for _, item := range visibleItems {
		if !item.isExcluded {
			fullPath := filepath.Join(m.path, item.name)
			if _, ok := m.selected[fullPath]; !ok {
				allSelectableAreSelected = false
				break
			}
		}
	}

	for _, item := range visibleItems {
		if !item.isExcluded {
			fullPath := filepath.Join(m.path, item.name)
			if allSelectableAreSelected {
				delete(m.selected, fullPath)
			} else {
				m.selected[fullPath] = struct{}{}
			}
		}
	}
}

func (m *Model) handleGoToTop() {
	if len(m.getVisibleItems()) > 0 {
		m.cursor = 0
	}
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
	pathValue := m.path
	if !strings.HasSuffix(pathValue, string(filepath.Separator)) {
		pathValue += string(filepath.Separator)
	}
	m.pathInput.SetValue(pathValue)

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
