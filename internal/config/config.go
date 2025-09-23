package config

import "strings"

type Config struct {
	ExcludedNames      map[string]struct{}
	ExcludedExtensions map[string]struct{}
}

func NewConfig() *Config {
	return &Config{
		ExcludedNames: map[string]struct{}{
			".git":              {},
			".svn":              {},
			".hg":               {},
			"node_modules":      {},
			"vendor":            {},
			"package-lock.json": {},
			"yarn.lock":         {},
			"pnpm-lock.yaml":    {},
			"bun.lockb":         {},
			".vscode":           {},
			".idea":             {},
			".DS_Store":         {},
			"Thumbs.db":         {},
			"bin":               {},
			"dist":              {},
			"build":             {},
			".cache":            {},
			"target":            {},
			".env":              {},
			"context.txt":       {},
		},
		ExcludedExtensions: map[string]struct{}{
			".jpg":  {},
			".jpeg": {},
			".jpe":  {},
			".png":  {},
			".gif":  {},
			".bmp":  {},
			".tiff": {},
			".tif":  {},
			".webp": {},
			".ico":  {},
			".heic": {},
			".heif": {},
			".avif": {},
			".jp2":  {},
			".j2k":  {},
			".jpf":  {},
			".jpx":  {},
			".jpm":  {},
			".mj2":  {},
			".svg":  {},
			".ai":   {},
			".eps":  {},
			".pdf":  {},
			".raw":  {},
			".cr2":  {},
			".nef":  {},
			".nrw":  {},
			".arw":  {},
			".srf":  {},
			".sr2":  {},
			".orf":  {},
			".dng":  {},
			".raf":  {},
			".pef":  {},
			".rw2":  {},
			".psd":  {},
			".xcf":  {},
			".indd": {},
			".tga":  {},
			".pcx":  {},
			".ppm":  {},
			".pgm":  {},
			".pbm":  {},
			".pnm":  {},
			".zip":  {},
			".tar":  {},
			".gz":   {},
			".bz2":  {},
			".xz":   {},
			".rar":  {},
			".7z":   {},
			".tgz":  {},
			".iso":  {},
			".dmg":  {},
			".jar":  {},
			".war":  {},
			".ear":  {},
		},
	}
}

func (c *Config) IsExcluded(name string) bool {
	if _, found := c.ExcludedNames[name]; found {
		return true
	}

	for ext := range c.ExcludedExtensions {
		if strings.HasSuffix(strings.ToLower(name), ext) {
			return true
		}
	}

	return false
}
