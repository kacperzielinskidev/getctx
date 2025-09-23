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
	if m.isInputMode {
		mainContent = m.renderCompletionGrid()
	} else {
		mainContent = m.viewport.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		mainContent,
		footer,
	)
}

func (m *Model) getVisibleItems() []item {
	if m.filterQuery == "" {
		return m.items
	}
	var filteredItems []item
	lowerQuery := strings.ToLower(m.filterQuery)
	for _, i := range m.items {
		if strings.Contains(strings.ToLower(i.name), lowerQuery) {
			filteredItems = append(filteredItems, i)
		}
	}
	return filteredItems
}

func (m *Model) renderPathInput() string {
	var s strings.Builder
	prompt := InputHeader
	if m.isFilterMode {
		prompt = FilterHeader
	}
	s.WriteString(prompt)
	s.WriteString(m.pathInput.View())

	if m.inputErrorMsg != "" {
		s.WriteString("\n" + Styles.Log.Error.Render(m.inputErrorMsg))
	}
	s.WriteString("\n")
	return s.String()
}

func (m *Model) renderHeader() string {
	if m.isInputMode || m.isFilterMode {
		return m.renderPathInput()
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

func (m *Model) renderFooter() string {
	return fmt.Sprintf(StatusFooterFormat, len(m.selected))
}

func (m *Model) renderFileList() string {
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		message := EmptyMessage
		if m.filterQuery != "" {
			message = NoMatchesMessage
		}
		padding := strings.Repeat("\n", m.viewport.Height/2)
		style := Styles.List.Empty.Width(m.viewport.Width).Align(lipgloss.Center)
		return style.Render(padding + message)
	}

	var s strings.Builder
	for i, item := range visibleItems {
		s.WriteString(m.renderListItem(i, item))
	}
	return s.String()
}

func (m *Model) renderListItem(index int, item item) string {
	var nameStyle lipgloss.Style
	fullPath := filepath.Join(m.path, item.name)
	_, isSelected := m.selected[fullPath]

	if item.isExcluded {
		nameStyle = Styles.List.Excluded
	} else if isSelected {
		nameStyle = Styles.List.Selected
	} else {
		nameStyle = Styles.List.Normal
	}

	cursorStr := Elements.List.CursorEmpty
	if m.cursor == index {
		cursorStr = Icons.Cursor
		nameStyle = nameStyle.Copy().Bold(true)
	}

	var prefix, itemIcon string
	if item.isExcluded {
		prefix = Elements.List.UnselectedPrefix
		itemIcon = Icons.Excluded
	} else {
		prefix = Elements.List.UnselectedPrefix
		if isSelected {
			prefix = Elements.List.SelectedPrefix
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

	finalLine := fmt.Sprintf("%s %s%s %s", cursorStr, prefix, itemIcon, itemName)
	return nameStyle.Render(finalLine) + "\n"
}

func (m *Model) renderCompletionGrid() string {
	suggestions := m.completionSuggestions
	if len(suggestions) == 0 {
		return lipgloss.Place(m.width, m.height-lipgloss.Height(m.renderHeader())-lipgloss.Height(m.renderFooter()), lipgloss.Center, lipgloss.Center, Styles.List.Empty.Render(NoMatchesMessage))
	}
	sort.Strings(suggestions)

	const maxCols = 7
	const padding = 2
	var bestNumCols = 1
	var bestColWidths []int

	for numCols := maxCols; numCols >= 1; numCols-- {
		numRows := (len(suggestions) + numCols - 1) / numCols
		if numRows == 0 {
			continue
		}
		colWidths := make([]int, numCols)
		totalWidth := (numCols - 1) * padding
		for c := 0; c < numCols; c++ {
			maxWidthInCol := 0
			for r := 0; r < numRows; r++ {
				i := c*numRows + r
				if i < len(suggestions) {
					maxWidthInCol = max(maxWidthInCol, len(suggestions[i]))
				}
			}
			colWidths[c] = maxWidthInCol
			totalWidth += maxWidthInCol
		}
		if totalWidth <= m.width {
			bestNumCols = numCols
			bestColWidths = colWidths
			break
		}
	}

	var grid strings.Builder
	numRows := (len(suggestions) + bestNumCols - 1) / bestNumCols
	for r := 0; r < numRows; r++ {
		for c := 0; c < bestNumCols; c++ {
			i := c*numRows + r
			if i < len(suggestions) {
				item := suggestions[i]
				grid.WriteString(item)
				if c < bestNumCols-1 {
					padCount := bestColWidths[c] - len(item) + padding
					grid.WriteString(strings.Repeat(" ", padCount))
				}
			}
		}
		grid.WriteString("\n")
	}
	return grid.String()
}
