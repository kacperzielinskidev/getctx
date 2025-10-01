package config

import (
	"path/filepath"
	"strings"
)

type Config struct {
	ExcludedNames      map[string]struct{}
	ExcludedExtensions map[string]struct{}
}

var defaultExcludedNames = []string{
	".git", ".svn", ".hg",
	"node_modules", "vendor",
	"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb",
	".vscode", ".idea",
	".DS_Store", "Thumbs.db",
	"bin", "dist", "build", "target",
	".cache",
	".env",
	"context.txt",
}

var defaultExcludedExtensions = []string{
	".jpg", ".jpeg", ".jpe", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp",
	".ico", ".heic", ".heif", ".avif", ".jp2", ".j2k", ".jpf", ".jpx", ".jpm",
	".mj2", ".svg", ".ai", ".eps",
	".pdf", ".psd", ".xcf", ".indd",
	".raw", ".cr2", ".nef", ".nrw", ".arw", ".srf", ".sr2", ".orf", ".dng",
	".raf", ".pef", ".rw2",
	".tga", ".pcx", ".ppm", ".pgm", ".pbm", ".pnm",
	".zip", ".tar", ".gz", ".bz2", ".xz", ".rar", ".7z", ".tgz", ".iso",
	".dmg", ".jar", ".war", ".ear",
}

func NewConfig() *Config {
	cfg := &Config{
		ExcludedNames:      make(map[string]struct{}),
		ExcludedExtensions: make(map[string]struct{}),
	}

	for _, name := range defaultExcludedNames {
		cfg.ExcludedNames[name] = struct{}{}
	}

	for _, ext := range defaultExcludedExtensions {
		cfg.ExcludedExtensions[ext] = struct{}{}
	}

	return cfg
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
