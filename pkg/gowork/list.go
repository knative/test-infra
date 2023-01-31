package gowork

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"runtime"

	"golang.org/x/mod/modfile"
)

// List returns the list of modules in the current Go workspace.
func List(filesystem FileSystem, env Environment) ([]Module, error) {
	gowork, err := findWorkfile(filesystem, env)
	if err != nil {
		return nil, err
	}
	data, err := fs.ReadFile(filesystem, gowork)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGowork, err)
	}
	work, err := modfile.ParseWork(gowork, data, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGowork, err)
	}
	workdir := filepath.Dir(gowork)
	mods := make([]Module, len(work.Use))
	for i, use := range work.Use {
		dir := filepath.Join(workdir, use.Path)
		modfilePath := filepath.Join(dir, "go.mod")
		data, err = fs.ReadFile(filesystem, modfilePath)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidGowork, err)
		}
		mod, err := modfile.ParseLax(modfilePath, data, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidGowork, err)
		}
		mods[i] = Module{
			Name: mod.Module.Mod.Path,
			Path: dir,
		}
	}
	return mods, nil
}

func findWorkfile(fs FileSystem, env Environment) (string, error) {
	switch gowork := env.Get("GOWORK"); gowork {
	case "off":
		return "", nil
	case "", "auto":
		f, err := findEnclosingFile(fs, "go.work", func(d string) bool {
			// As a special case, don't cross GOROOT to find a go.work file.
			// The standard library and commands built in go always use the vendored
			// dependencies, so avoid using a most likely irrelevant go.work file.
			return d == runtime.GOROOT()
		})
		if err != nil {
			return "", err
		}
		if f == "" {
			return "", fmt.Errorf("%w: file not found", ErrInvalidGowork)
		}
		return f, nil
	default:
		if !path.IsAbs(gowork) {
			return "", fmt.Errorf("%w: GOWORK must be an absolute path", ErrInvalidGowork)
		}
		p, err := filepath.Rel("/", gowork)
		if err != nil {
			return "", fmt.Errorf("%w: Cant compute relative path "+
				"to go.work to '/': %v", ErrBug, gowork)
		}
		return p, nil
	}
}
