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
		emptyMessage := Elements.Text.NoMatchesMessage
		availableHeight := m.height - lipgloss.Height(m.renderHeader()) - lipgloss.Height(m.renderFooter())
		padding := strings.Repeat("\n", availableHeight/2)

		style := Styles.List.Empty.Width(m.width).Align(lipgloss.Center)

		return style.Render(padding + emptyMessage)

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

		// Obliczamy maksymalną szerokość dla każdej kolumny indywidualnie.
		for c := 0; c < numCols; c++ {
			maxWidthInCol := 0
			// Układ "ls" wypełnia najpierw kolumny, potem wiersze.
			for r := 0; r < numRows; r++ {
				index := c*numRows + r
				if index < len(suggestions) {
					if len(suggestions[index]) > maxWidthInCol {
						maxWidthInCol = len(suggestions[index])
					}
				}
			}
			colWidths[c] = maxWidthInCol
		}

		// Sumujemy szerokości wszystkich kolumn oraz padding między nimi.
		totalWidth = (numCols - 1) * padding
		for _, w := range colWidths {
			totalWidth += w
		}

		// Jeśli ten układ mieści się na ekranie, znaleźliśmy optymalne rozwiązanie.
		if totalWidth <= m.width {
			bestNumCols = numCols
			bestColWidths = colWidths
			break
		}
	}
	// Jeśli nawet układ jednokolumnowy się nie mieści (bardzo długa nazwa pliku),
	// to i tak pozostanie on jednokolumnowy, co jest poprawnym zachowaniem.

	// Renderujemy siatkę, używając obliczonego optymalnego układu.
	var grid strings.Builder
	numRows := (len(suggestions) + bestNumCols - 1) / bestNumCols

	for r := 0; r < numRows; r++ {
		for c := 0; c < bestNumCols; c++ {
			index := c*numRows + r
			if index < len(suggestions) {
				item := suggestions[index]
				grid.WriteString(item)

				// Dodajemy padding tylko wtedy, gdy to nie jest ostatnia kolumna.
				if c < bestNumCols-1 {
					// Obliczamy padding potrzebny do wyrównania do następnej kolumny.
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
