package gowork

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const rootDir = "/"

var (
	systemFS                 = os.DirFS(rootDir) //nolint:gochecknoglobals
	_        operatingsystem = RealSystem{}
)

// RealSystem is the real, live OS implementation.
type RealSystem struct{}

func (s RealSystem) Get(name string) string {
	return os.Getenv(name)
}

func (s RealSystem) Open(name string) (fs.File, error) {
	f, err := systemFS.Open(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBug, err)
	}
	return f, nil
}

func (s RealSystem) Cwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrBug, err)
	}
	p, err := filepath.Rel(rootDir, dir)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrBug, err)
	}
	return p, nil
}

func (s RealSystem) Abs(filepath string) string {
	return path.Clean(path.Join(rootDir, filepath))
}

// OS represents a virtual operating system.
type operatingsystem interface {
	FileSystem
	Environment
	Abs(filepath string) string
}
