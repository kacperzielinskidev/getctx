package app

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
	prompt := Elements.Text.InputHeader
	if m.isFilterMode {
		prompt = Elements.Text.FilterHeader
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

	helpHeader := Elements.Text.HelpHeader
	pathStyle := lipgloss.NewStyle().Width(m.width)
	fullPathString := Elements.Text.PathPrefix + m.path + filterIndicator
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
	visibleItems := m.getVisibleItems()
	if len(visibleItems) == 0 {
		return m.renderEmptyFileList()
	}

	var s strings.Builder
	for i, item := range visibleItems {
		s.WriteString(m.renderListItem(i, item))
	}
	return s.String()
}

func (m *Model) renderEmptyFileList() string {
	emptyMessage := Elements.Text.EmptyMessage
	padding := strings.Repeat("\n", m.viewport.Height/2)

	style := Styles.List.Empty.Width(m.viewport.Width).Align(lipgloss.Center)

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

	line := lipgloss.JoinHorizontal(lipgloss.Top, finalCursor, finalPrefix, finalName)

	return line + "\n"
}

func (m *Model) renderCompletionGrid() string {
	suggestions := m.completionSuggestions
	if len(suggestions) == 0 {
		availableHeight := max(0, m.height-lipgloss.Height(m.renderHeader())-lipgloss.Height(m.renderFooter()))
		message := Styles.List.Empty.Render(Elements.Text.NoMatchesMessage)

		return lipgloss.Place(
			m.width,
			availableHeight,
			lipgloss.Center,
			lipgloss.Center,
			message,
		)
	}

	sort.Strings(suggestions)

	const maxCols = 7
	const padding = 2

	var (
		bestNumCols   = 1
		bestColWidths []int
	)

	for numCols := maxCols; numCols >= 1; numCols-- {
		numRows := (len(suggestions) + numCols - 1) / numCols
		if numRows == 0 {
			continue
		}

		colWidths := make([]int, numCols)
		totalWidth := 0

		for c := range numCols {
			maxWidthInCol := 0
			for r := range numRows {
				index := c*numRows + r
				if index < len(suggestions) {
					maxWidthInCol = max(maxWidthInCol, len(suggestions[index]))
				}
			}
			colWidths[c] = maxWidthInCol
		}

		totalWidth = (numCols - 1) * padding
		for _, w := range colWidths {
			totalWidth += w
		}

		if totalWidth <= m.width {
			bestNumCols = numCols
			bestColWidths = colWidths
			break
		}
	}

	var grid strings.Builder
	numRows := (len(suggestions) + bestNumCols - 1) / bestNumCols

	for r := range numRows {
		for c := range bestNumCols {
			index := c*numRows + r
			if index < len(suggestions) {
				item := suggestions[index]
				grid.WriteString(item)

				if c < bestNumCols-1 {
					padCount := bestColWidths[c] - len(item) + padding
					grid.WriteString(strings.Repeat(" ", padCount))
				}
			}
		}
		grid.WriteString("\n")
	}

	return lipgloss.NewStyle().
		Height(m.height - lipgloss.Height(m.renderHeader()) - lipgloss.Height(m.renderFooter())).
		Render(grid.String())
}
