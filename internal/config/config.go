package config

import (
	"path/filepath"
	"strings"
)

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
			".cache":            {},
			"target":            {},
			".env":              {},
			"context.txt":       {},
		},
		ExcludedExtensions: map[string]struct{}{
			".jpg": {}, ".jpeg": {}, ".jpe": {}, ".png": {}, ".gif": {},
			".bmp": {}, ".tiff": {}, ".tif": {}, ".webp": {}, ".ico": {},
			".heic": {}, ".heif": {}, ".avif": {}, ".jp2": {}, ".j2k": {},
			".jpf": {}, ".jpx": {}, ".jpm": {}, ".mj2": {}, ".svg": {},
			".ai": {}, ".eps": {}, ".pdf": {}, ".raw": {}, ".cr2": {},
			".nef": {}, ".nrw": {}, ".arw": {}, ".srf": {}, ".sr2": {},
			".orf": {}, ".dng": {}, ".raf": {}, ".pef": {}, ".rw2": {},
			".psd": {}, ".xcf": {}, ".indd": {}, ".tga": {}, ".pcx": {},
			".ppm": {}, ".pgm": {}, ".pbm": {}, ".pnm": {}, ".zip": {},
			".tar": {}, ".gz": {}, ".bz2": {}, ".xz": {}, ".rar": {},
			".7z": {}, ".tgz": {}, ".iso": {}, ".dmg": {}, ".jar": {},
			".war": {}, ".ear": {},
		},
	}
}

func (c *Config) IsExcluded(name string) bool {
	if _, found := c.ExcludedNames[name]; found {
		return true
	}

	ext := filepath.Ext(name)
	if ext == "" {
		return false
	}

	if _, found := c.ExcludedExtensions[strings.ToLower(ext)]; found {
		return true
	}

	return false
}
