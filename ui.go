package main

import "github.com/charmbracelet/lipgloss"

type UIIcons struct {
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

type UIColors struct {
	Red   string
	Reset string
}

type UIStyles struct {
	Selected lipgloss.Style
	Cursor   lipgloss.Style
}

var Icons UIIcons
var Colors UIColors
var Styles UIStyles

func init() {
	Icons = UIIcons{
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

	Colors = UIColors{
		Red:   "\033[31m",
		Reset: "\033[0m",
	}

	Styles = UIStyles{
		Selected: lipgloss.NewStyle().Foreground(lipgloss.Color("34")), // "34" to kod koloru zielonego
		Cursor:   lipgloss.NewStyle().Bold(true),
	}
}
