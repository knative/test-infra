package test

import "io/fs"

// FS is a testable implementation of modules.FileSystem.
type FS struct {
	Dir   string
	Files fs.FS
}

// Open implements modules.FileSystem.
func (e FS) Open(name string) (fs.File, error) {
	return e.Files.Open(name)
}

// Cwd implements modules.FileSystem.
func (e FS) Cwd() (string, error) {
	return e.Dir, nil
}
