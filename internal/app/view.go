package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
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

	filterIndicator := FormatFilterIndicator(m.filterQuery)

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
