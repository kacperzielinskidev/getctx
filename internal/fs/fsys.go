package fs

import (
	"io"
	"io/fs"
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
