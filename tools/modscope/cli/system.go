package cli

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"knative.dev/test-infra/tools/modscope/modules"
)

const rootDir = "/"

var systemFS = os.DirFS(rootDir)

type OS interface {
	modules.FileSystem
	modules.Environment
	Abs(filepath string) string
}

type system struct{}

var _ OS = system{}

func (s system) Get(name string) string {
	return os.Getenv(name)
}

func (s system) Open(name string) (fs.File, error) {
	return systemFS.Open(name)
}

func (s system) Cwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Rel(rootDir, dir)
}

func (s system) Abs(filepath string) string {
	return path.Clean(path.Join(rootDir, filepath))
}
