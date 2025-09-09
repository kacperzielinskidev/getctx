package app

type Config struct {
	ExcludedNames map[string]struct{}
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		ExcludedNames: map[string]struct{}{
			// Version control
			".git": {},
			".svn": {},
			".hg":  {},

			// Dependencies
			"node_modules":      {},
			"vendor":            {},
			"package-lock.json": {},
			"yarn.lock":         {},
			"pnpm-lock.yaml":    {},
			"bun.lockb":         {},

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

			".env": {},

			"context.txt": {},
		},
	}
}

func (c *Config) IsExcluded(name string) bool {
	_, found := c.ExcludedNames[name]
	return found
}
