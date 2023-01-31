package testdata

import (
	"io/fs"
	"os"
	"testing/fstest"
	"time"
)

// ExampleFS returns a file system with example data.
func ExampleFS() fs.FS {
	return fstest.MapFS{
		"srv":      dir(),
		"code":     dir(),
		"code/foo": dir(),
		"code/foo/go.work": file(`go 1.18

use (
	.
	a
)
`),
		"code/foo/go.mod": file(`module foo

go 1.18
`),
		"code/foo/a": dir(),
		"code/foo/a/go.mod": file(`module foo/a

go 1.18
`),
		"code/foo/b": dir(),
		"code/foo/b/go.mod": file(`module foo/b

go 1.18
`),
		"code/bar": dir(),
		"code/bar/go.mod": file(`module bar

go 1.18
`),
	}
}

func file(data string) *fstest.MapFile {
	return &fstest.MapFile{
		Data:    []byte(data),
		Mode:    0o644,
		ModTime: time.Now(),
	}
}

func dir() *fstest.MapFile {
	return &fstest.MapFile{
		Mode:    os.ModeDir,
		ModTime: time.Now(),
	}
}
