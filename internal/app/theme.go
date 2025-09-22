package app

import "github.com/charmbracelet/lipgloss"

const (
	colorGreen lipgloss.Color = "34"
	colorRed   lipgloss.Color = "9"
	colorCyan  lipgloss.Color = "86"
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
	InputHeader  string
	PathPrefix   string
	StatusFooter string
	EmptyMessage string
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
	Normal   lipgloss.Style
	Hint     lipgloss.Style
	Empty    lipgloss.Style
}

type TUILogStyles struct {
	Skipped lipgloss.Style
	Error   lipgloss.Style
}

type TUIStyles struct {
	List TUIListStyles
	Log  TUILogStyles
}

type TUITexts struct {
	HelpHeaderBase  string
	HelpHeaderHint  string
	InputHeaderBase string
	InputHeaderHint string
	PathPrefix      string
	StatusFooter    string
	EmptyMessage    string
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
			Skipped: lipgloss.NewStyle().Foreground(Colors.Red),
			Error:   lipgloss.NewStyle().Foreground(Colors.Red).Bold(true),
		},
	}

	texts := TUITexts{
		HelpHeaderBase:  "Select files ",
		HelpHeaderHint:  "(SPACE: Select File, CTRL+HOME: Go to Top, CTRL+END: Go to Bottom, CTRL+P: Find Path, Q: Save)",
		InputHeaderBase: "Enter path ",
		InputHeaderHint: "(ENTER: Confirm, ESCAPE: Cancel, CTRL+W: Remove whole line)",
		PathPrefix:      "Current path: ",
		StatusFooter:    "\nSelected %d items. Press 'q' to save and exit.",
		EmptyMessage:    "[ This directory is empty ]",
	}

	Elements = TUIElements{
		List: TUIListElements{
			CursorEmpty:      " ",
			SelectedPrefix:   Icons.Checkmark + " ",
			UnselectedPrefix: "  ",
			DirectorySuffix:  "/",
		},
		Text: TUITextElements{
			HelpHeader: lipgloss.JoinHorizontal(lipgloss.Left,
				texts.HelpHeaderBase,
				Styles.List.Hint.Render(texts.HelpHeaderHint),
			) + "\n",
			InputHeader: lipgloss.JoinHorizontal(lipgloss.Left,
				texts.InputHeaderBase,
				Styles.List.Hint.Render(texts.InputHeaderHint),
			) + "\n",
			PathPrefix:   texts.PathPrefix,
			StatusFooter: texts.StatusFooter,
			EmptyMessage: texts.EmptyMessage,
		},
	}

}
