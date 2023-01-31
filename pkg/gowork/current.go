package gowork

import (
	"fmt"
	"io/fs"
	"path"

	"golang.org/x/mod/modfile"
)

// Current returns the current Go module.
func Current(filesystem FileSystem, env Environment) (*Module, error) {
	moduleFilepath, err := findModfile(filesystem, env)
	if err != nil {
		return nil, err
	}
	data, err := fs.ReadFile(filesystem, moduleFilepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGomod, err)
	}
	mod, err := modfile.ParseLax(moduleFilepath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGomod, err)
	}
	return &Module{
		Name: mod.Module.Mod.Path,
		Path: path.Dir(moduleFilepath),
	}, nil
}

func findModfile(fs FileSystem, env Environment) (string, error) {
	switch go111module := env.Get("GO111MODULE"); go111module {
	case "", "auto", "on":
		f, err := findEnclosingFile(fs, "go.mod")
		if err != nil {
			return "", err
		}
		if f == "" {
			return "", fmt.Errorf("%w: file not found", ErrInvalidGomod)
		}
		return f, nil
	default:
		return "", fmt.Errorf("%w: unsupported value of GO111MODULE=%v "+
			"(supported are '', 'auto', 'on')", ErrInvalidGomod, go111module)
	}
}
