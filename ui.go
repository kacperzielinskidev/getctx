package main

type UIIcons struct {
	Building string
	Done     string
	Warning  string
	Info     string
	Error    string
}

type UIColors struct {
	Red   string
	Reset string
}

var Icons UIIcons
var Colors UIColors

func init() {
	Icons = UIIcons{
		Building: "🚀",
		Done:     "✅",
		Warning:  "⚠️",
		Info:     "ℹ️",
		Error:    "❌",
	}

	Colors = UIColors{
		Red:   "\033[31m",
		Reset: "\033[0m",
	}
}
