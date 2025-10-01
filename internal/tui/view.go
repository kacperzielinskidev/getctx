package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	var mainContent string
	switch m.mode {
	case modePathInput:
		mainContent = m.renderCompletionView()
	default:
		m.viewport.SetContent(m.renderFileListView())
		mainContent = m.viewport.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		mainContent,
		footer,
	)
}

func (m *Model) getVisibleItems() []listItem {
	if m.filterQuery == "" {
		return m.items
	}

	var filteredItems []listItem
	lowerQuery := strings.ToLower(m.filterQuery)
	for _, i := range m.items {
		if strings.Contains(strings.ToLower(i.name), lowerQuery) {
			filteredItems = append(filteredItems, i)
		}
	}
	return filteredItems
}

func (m *Model) renderHeader() string {
	if m.mode == modePathInput || m.mode == modeFilter {
		return m.renderTextInput()
	}

	filterIndicator := formatFilterIndicator(m.filterQuery)
	pathStyle := lipgloss.NewStyle().Width(m.width)
	fullPathString := PathPrefix + m.path + filterIndicator
	wrappedPath := pathStyle.Render(fullPathString)

	return lipgloss.JoinVertical(lipgloss.Left,
		HelpHeader,
		wrappedPath,
	)
}

func (m *Model) renderTextInput() string {
	var s strings.Builder

	prompt := InputHeader
	if m.mode == modeFilter {
		prompt = FilterHeader
	}

	s.WriteString(prompt)
	s.WriteString(m.textInput.View())

	if m.inputErrorMsg != "" {
		s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
	}
	s.WriteString("\n")
	return s.String()
}

func (m *Model) renderFooter() string {
	return fmt.Sprintf(StatusFooterFormat, len(m.selected))
}

func (m *Model) renderFileListView() string {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		message := EmptyMessage
		if m.filterQuery != "" {
			message = NoMatchesMessage
		}
		style := Styles.List.Empty.Width(m.viewport.Width).Height(m.viewport.Height).Align(lipgloss.Center, lipgloss.Center)
		return style.Render(message)
	}
	var s strings.Builder
	for i, item := range visibleItems {
		s.WriteString(m.renderListItem(i, item))
	}
	return s.String()
}

func (m *Model) renderListItem(index int, item listItem) string {
	fullPath := filepath.Join(m.path, item.name)
	_, isSelected := m.selected[fullPath]
	isCursorOnItem := m.cursor == index

	var style lipgloss.Style

	if item.isExcluded {
		style = Styles.List.Excluded
	} else if isSelected {
		style = Styles.List.Selected
	} else {
		style = Styles.List.Normal
	}

	cursorStr := Elements.List.CursorEmpty
	if isCursorOnItem {
		cursorStr = Icons.Cursor
	}

	prefix := Elements.List.UnselectedPrefix
	if isSelected && !item.isExcluded {
		prefix = Elements.List.SelectedPrefix
	}

	icon := Icons.File
	if item.isDir {
		icon = Icons.Directory
	}
	if item.isExcluded {
		icon = Icons.Excluded
	}

	itemName := item.name
	if item.isDir {
		itemName += Elements.List.DirectorySuffix
	}

	line := fmt.Sprintf("%s %s%s %s", cursorStr, prefix, icon, itemName)
	return style.Render(line) + "\n"
}

func (m *Model) renderCompletionView() string {
	suggestions := m.completionSuggestions
	if len(suggestions) == 0 {
		return lipgloss.Place(m.width, m.completionViewport.Height,
			lipgloss.Center, lipgloss.Center,
			Styles.List.Empty.Render(NoMatchesMessage),
		)
	}

	sort.Strings(suggestions)

	numCols, colWidths := calculateGridDimensions(suggestions, m.width)
	gridContent := buildGrid(suggestions, numCols, colWidths)

	m.completionViewport.SetContent(gridContent)
	return m.completionViewport.View()
}

func calculateGridDimensions(suggestions []string, maxWidth int) (numCols int, colWidths []int) {
	const maxCols = 7
	const padding = 2
	bestNumCols := 1

	for cols := maxCols; cols >= 1; cols-- {
		numRows := (len(suggestions) + cols - 1) / cols
		if numRows == 0 {
			continue
		}

		currentWidths := make([]int, cols)
		totalWidth := (cols - 1) * padding

		for c := 0; c < cols; c++ {
			maxWidthInCol := 0
			for r := 0; r < numRows; r++ {
				i := c*numRows + r
				if i < len(suggestions) {
					maxWidthInCol = max(maxWidthInCol, len(suggestions[i]))
				}
			}
			currentWidths[c] = maxWidthInCol
			totalWidth += maxWidthInCol
		}

		if totalWidth <= maxWidth {
			bestNumCols = cols
			colWidths = currentWidths
			return bestNumCols, colWidths
		}
	}
	colWidths = []int{maxWidth}
	return 1, colWidths
}

func buildGrid(suggestions []string, numCols int, colWidths []int) string {
	var grid strings.Builder
	numRows := (len(suggestions) + numCols - 1) / numCols

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			i := c*numRows + r
			if i < len(suggestions) {
				item := suggestions[i]
				grid.WriteString(item)
				if c < numCols-1 {
					padCount := colWidths[c] - len(item) + 2 // 2 to padding
					grid.WriteString(strings.Repeat(" ", padCount))
				}
			}
		}
		grid.WriteString("\n")
	}

	return grid.String()
}
