package test

import "io/fs"

type FS struct {
	Dir   string
	Files fs.FS
}

func (e FS) Open(name string) (fs.File, error) {
	return e.Files.Open(name)
}

func (e FS) Cwd() (string, error) {
	return e.Dir, nil
}
