package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"knative.dev/test-infra/pkg/gowork"
)

const rootDir = "/"

var systemFS = os.DirFS(rootDir) //nolint:gochecknoglobals

// OS represents a virtual operating system.
type OS interface {
	gowork.FileSystem
	gowork.Environment
	Abs(filepath string) string
}

type system struct{}

var _ OS = system{}

func (s system) Get(name string) string {
	return os.Getenv(name)
}

func (s system) Open(name string) (fs.File, error) {
	f, err := systemFS.Open(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", gowork.ErrBug, err)
	}
	return f, nil
}

func (s system) Cwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w: %v", gowork.ErrBug, err)
	}
	p, err := filepath.Rel(rootDir, dir)
	if err != nil {
		return "", fmt.Errorf("%w: %v", gowork.ErrBug, err)
	}
	return p, nil
}

func (s system) Abs(filepath string) string {
	return path.Clean(path.Join(rootDir, filepath))
}
