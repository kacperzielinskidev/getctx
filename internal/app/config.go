package app

var ExcludedNames = map[string]struct{}{
	// Version control
	".git": {},
	".svn": {},
	".hg":  {},

	// Dependencies
	"node_modules": {},
	"vendor":       {},

	// IDE configuration
	".vscode": {},
	".idea":   {},

	// System files
	".DS_Store": {},
	"Thumbs.db": {},

	// Build artifacts & cache
	"bin":    {},
	"dist":   {},
	"build":  {},
	".cache": {},
	"target": {},
}
