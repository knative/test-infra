package gowork

import (
	"io/fs"
	"path/filepath"
)

// FileSystem is used to read files from the file system in a testable way.
type FileSystem interface {
	fs.FS
	// Cwd works like os.Getwd
	Cwd() (string, error)
}

// Environment is used to get environment variables.
type Environment interface {
	Get(name string) string
}

func findEnclosingFile(fs FileSystem, file string, drops ...func(string) bool) (string, error) {
	dir, err := fs.Cwd()
	if err != nil {
		return "", errWrap(err, ErrBug)
	}
	dir = filepath.Clean(dir)

	// Look for enclosing file.
	for {
		fp := filepath.Join(dir, file)
		f, err2 := fs.Open(fp)
		if err2 == nil {
			defer func() {
				_ = f.Close()
			}()
			fi, err := f.Stat()
			if err != nil {
				return "", errWrap(err, ErrBug)
			}

			if !fi.IsDir() {
				return fp, nil
			}
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		for _, drop := range drops {
			if drop(d) {
				return "", nil
			}
		}
		dir = d
	}
	return "", nil
}
