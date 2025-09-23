package app

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateInputMode(msg tea.Msg) tea.Cmd {
	// Zapamiętaj wartość pola tekstowego PRZED aktualizacją
	oldValue := m.pathInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyTab:
			// TAB teraz tylko uzupełnia, nie generuje listy.
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

	m.pathInput, _ = m.pathInput.Update(msg)

	if m.pathInput.Value() != oldValue {
		m.updateCompletions()
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
		m.completionSuggestions = nil
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
	m.pathInput.SetCursor(len(pathValue)) // Ustaw kursor na końcu

	// KLUCZOWA ZMIANA: Natychmiast generuj sugestie dla bieżącej ścieżki.
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

func (m *Model) handleTabCompletion() {
	currentInput := m.pathInput.Value()
	dirToSearch, prefix := m.getCompletionParts(currentInput)

	// 1. Pobierz sugestie dla aktualnie wpisanego tekstu.
	initialSuggestions, err := m.getCompletions(dirToSearch, prefix)
	if err != nil || len(initialSuggestions) == 0 {
		m.completionSuggestions = nil // Brak sugestii, wyczyść widok siatki.
		return
	}

	// 2. Ustal, jaka będzie nowa, uzupełniona wartość pola tekstowego.
	var newInputValue string
	if len(initialSuggestions) == 1 {
		// Scenariusz A: Tylko jedna możliwość. Uzupełnij ją w całości.
		newInputValue = filepath.Join(dirToSearch, initialSuggestions[0])
	} else {
		// Scenariusz B: Wiele możliwości. Znajdź najdłuższy wspólny prefiks.
		commonPrefix := findLongestCommonPrefix(initialSuggestions)
		if commonPrefix == "" {
			// Jeśli nie ma wspólnego prefiksu, tylko pokaż listę i nie zmieniaj tekstu.
			m.completionSuggestions = initialSuggestions
			return
		}
		newInputValue = filepath.Join(dirToSearch, commonPrefix)
	}

	// 3. Sprawdź, czy nowo uzupełniona ścieżka jest katalogiem.
	info, err := m.fsys.Stat(newInputValue)
	isDir := err == nil && info.IsDir()

	// 4. Zaktualizuj pole tekstowe i zdecyduj, jaka lista sugestii ma się pojawić.
	m.pathInput.SetValue(newInputValue)
	m.pathInput.SetCursor(len(newInputValue))

	if isDir {
		// To jest katalog! Upewnij się, że ma ukośnik na końcu...
		if !strings.HasSuffix(newInputValue, string(filepath.Separator)) {
			newInputValue += string(filepath.Separator)
			m.pathInput.SetValue(newInputValue) // Zaktualizuj ponownie, by dodać ukośnik
			m.pathInput.SetCursor(len(newInputValue))
		}
		// ...a następnie NATYCHMIAST pobierz i wyświetl jego zawartość.
		finalSuggestions, _ := m.getCompletions(newInputValue, "")
		m.completionSuggestions = finalSuggestions
	} else {
		// To nie jest katalog (np. plik lub częściowy prefiks).
		// Pokaż oryginalną listę sugestii, która doprowadziła do uzupełnienia.
		m.completionSuggestions = initialSuggestions
	}
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

// handleAutoComplete jest wywoływana przez TAB i uzupełnia tekst na podstawie widocznych sugestii.
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
	if err == nil && info.IsDir() {
		if !strings.HasSuffix(newInputValue, string(filepath.Separator)) {
			newInputValue += string(filepath.Separator)
		}
	}

	m.pathInput.SetValue(newInputValue)
	m.pathInput.SetCursor(len(newInputValue))

	// Po uzupełnieniu, natychmiast odśwież listę sugestii.
	m.updateCompletions()
}

func (m *Model) handleInputModeKeys(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}

	switch keyMsg.String() {
	case KeyTab:
		m.handleAutoComplete()
		return nil, true // Klucz obsłużony
	case KeyEnter:
		m.handleConfirmPathChange()
		return nil, true // Klucz obsłużony
	case KeyEscape, KeyCtrlC:
		m.handleCancelPathChange()
		return nil, true // Klucz obsłużony
	}

	return nil, false // To nie był nasz skrót, pozwól polu tekstowemu go obsłużyć.
}
