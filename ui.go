package main

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
}

type TUIColors struct {
	Green lipgloss.Color
	Red   lipgloss.Color
}

type TUIListStyles struct {
	Selected lipgloss.Style
	Cursor   lipgloss.Style
}

type TUILogStyles struct {
	Skipped lipgloss.Style
}

type TUIStyles struct {
	List TUIListStyles
	Log  TUILogStyles
}

var Icons TUIIcons
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
	}

	Colors = TUIColors{
		Green: colorGreen,
		Red:   colorRed,
	}

	Styles = TUIStyles{
		List: TUIListStyles{
			Selected: lipgloss.NewStyle().Foreground(Colors.Green),
			Cursor:   lipgloss.NewStyle().Bold(true),
		},
		Log: TUILogStyles{
			Skipped: lipgloss.NewStyle().Foreground(Colors.Red),
		},
	}
}
