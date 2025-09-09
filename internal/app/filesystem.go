package app

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type FileSystem interface {
	ReadDir(name string) ([]fs.DirEntry, error)
	Stat(name string) (fs.FileInfo, error)
	Abs(path string) (string, error)
	ReadFile(name string) ([]byte, error)
	Create(name string) (io.WriteCloser, error)
	WalkDir(root string, fn fs.WalkDirFunc) error
	UserHomeDir() (string, error)
	Open(name string) (fs.File, error)
}

type OSFileSystem struct{}

func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

func (fsys *OSFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (fsys *OSFileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (fsys *OSFileSystem) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (fsys *OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fsys *OSFileSystem) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (fsys *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

func (fsys *OSFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (fsys *OSFileSystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}
