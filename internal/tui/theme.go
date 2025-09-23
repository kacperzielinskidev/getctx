package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	colorGreen lipgloss.Color = "34"
	colorRed   lipgloss.Color = "9"
	colorCyan  lipgloss.Color = "86"
)

type TUIIcons struct {
	Directory string
	File      string
	Checkmark string
	Cursor    string
	Excluded  string
}

type TUIListElements struct {
	CursorEmpty      string
	SelectedPrefix   string
	UnselectedPrefix string
	DirectorySuffix  string
}

type TUIElements struct {
	List TUIListElements
}

type TUIColors struct {
	Green lipgloss.Color
	Red   lipgloss.Color
}

type TUIListStyles struct {
	Selected lipgloss.Style
	Cursor   lipgloss.Style
	Excluded lipgloss.Style
	Normal   lipgloss.Style
	Hint     lipgloss.Style
	Empty    lipgloss.Style
}

type TUILogStyles struct {
	Error lipgloss.Style
}

type TUIStyles struct {
	List TUIListStyles
	Log  TUILogStyles
}

var Icons TUIIcons
var Elements TUIElements
var Colors TUIColors
var Styles TUIStyles
var HelpHeader string
var InputHeader string
var FilterHeader string
var FilterIndicatorFormat string
var PathPrefix string
var StatusFooterFormat string
var EmptyMessage string
var NoMatchesMessage string

func init() {
	Icons = TUIIcons{
		Directory: "📁",
		File:      "📄",
		Checkmark: "✔",
		Cursor:    "❯",
		Excluded:  "🚫",
	}

	Colors = TUIColors{
		Green: colorGreen,
		Red:   colorRed,
	}

	Styles = TUIStyles{
		List: TUIListStyles{
			Selected: lipgloss.NewStyle().Foreground(Colors.Green),
			Cursor:   lipgloss.NewStyle().Bold(true),
			Excluded: lipgloss.NewStyle().Faint(true),
			Normal:   lipgloss.NewStyle(),
			Hint:     lipgloss.NewStyle().Foreground(colorCyan),
			Empty:    lipgloss.NewStyle().Faint(true),
		},
		Log: TUILogStyles{
			Error: lipgloss.NewStyle().Foreground(Colors.Red).Bold(true),
		},
	}

	Elements = TUIElements{
		List: TUIListElements{
			CursorEmpty:      " ",
			SelectedPrefix:   Icons.Checkmark + " ",
			UnselectedPrefix: "  ",
			DirectorySuffix:  "/",
		},
	}

	HelpHeader = lipgloss.JoinHorizontal(lipgloss.Left,
		"Select files ",
		Styles.List.Hint.Render("(space: select, a: select all, ctrl+p: find path, /: filter, q: save, ctrl+c: quit)"),
	) + "\n"

	InputHeader = lipgloss.JoinHorizontal(lipgloss.Left,
		"Enter path ",
		Styles.List.Hint.Render("(enter: confirm, esc: cancel, tab: autocomplete)"),
	) + "\n"

	FilterHeader = lipgloss.JoinHorizontal(lipgloss.Left,
		"Filter ",
		Styles.List.Hint.Render("(type to filter, enter: confirm, esc: cancel)"),
	) + "\n"

	FilterIndicatorFormat = " [Filtering by: \"%s\"]"
	PathPrefix = "Current path: "
	StatusFooterFormat = "\nSelected %d items. Press 'q' to save and exit."
	EmptyMessage = "[ This directory is empty ]"
	NoMatchesMessage = "[ No matching files or directories found ]"
}

func formatFilterIndicator(query string) string {
	if query == "" {
		return ""
	}
	indicator := fmt.Sprintf(FilterIndicatorFormat, query)
	return Styles.List.Hint.Render(indicator)
}

func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	}
	if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
	}
}
