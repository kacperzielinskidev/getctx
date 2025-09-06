package app

import "github.com/charmbracelet/lipgloss"

const (
	colorGreen lipgloss.Color = "34"
	colorRed   lipgloss.Color = "9"
)

type TUIIcons struct {
	Building  string
	Done      string
	Warning   string
	Info      string
	Error     string
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

type TUITextElements struct {
	HelpHeader   string
	PathPrefix   string
	StatusFooter string
}

type TUIElements struct {
	List TUIListElements
	Text TUITextElements
}

type TUIColors struct {
	Green lipgloss.Color
	Red   lipgloss.Color
}

type TUIListStyles struct {
	Selected lipgloss.Style
	Cursor   lipgloss.Style
	Excluded lipgloss.Style
}

type TUILogStyles struct {
	Skipped lipgloss.Style
}

type TUIStyles struct {
	List TUIListStyles
	Log  TUILogStyles
}

var Icons TUIIcons
var Elements TUIElements
var Colors TUIColors
var Styles TUIStyles

func init() {
	Icons = TUIIcons{
		Building:  "🚀",
		Done:      "✅",
		Warning:   "⚠️",
		Info:      "ℹ️",
		Error:     "❌",
		Directory: "📁",
		File:      "📄",
		Checkmark: "✔",
		Cursor:    "❯",
		Excluded:  "🚫",
	}

	Elements = TUIElements{
		List: TUIListElements{
			CursorEmpty:      " ",
			SelectedPrefix:   Icons.Checkmark + " ",
			UnselectedPrefix: "  ",
			DirectorySuffix:  "/",
		},
		Text: TUITextElements{
			HelpHeader:   "Select files for context (space: toggle, enter: open, backspace: up, q: save & quit)\n",
			PathPrefix:   "Current path: ",
			StatusFooter: "\nSelected %d items. Press 'q' to save and exit.",
		},
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
		},
		Log: TUILogStyles{
			Skipped: lipgloss.NewStyle().Foreground(Colors.Red),
		},
	}
}
