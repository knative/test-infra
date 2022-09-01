package shx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CopyOption int

const (
	// CopyNoOverwrite does not overwrite existing files in the destination
	CopyDefault CopyOption = iota
	CopyNoOverwrite
	CopyRecursive
)

// Copy a file or directory with the specified set of CopyOption.
// The source may use globbing, which is resolved with filepath.Glob.
// Notes:
//   * Does not copy file owner/group.
func Copy(src string, dest string, opts ...CopyOption) error {
	items, err := filepath.Glob(src)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("no such file or directory '%s'", src)
	}

	var combinedOpts CopyOption
	for _, opt := range opts {
		combinedOpts |= opt
	}

	// Check if the destination exists, e.g. if we are copying to /tmp/foo, /tmp should already exist
	if _, err := os.Stat(filepath.Dir(dest)); err != nil {
		return err
	}

	for _, item := range items {
		err := copyFileOrDirectory(item, dest, combinedOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileOrDirectory(src string, dest string, opts CopyOption) error {
	// If the destination is a directory that exists,
	// copy into the directory.
	destInfo, err := os.Stat(dest)
	if err == nil && destInfo.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}

	return filepath.Walk(src, func(srcPath string, srcInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only copy the first item if CopyRecursive wasn't set
		if opts&CopyRecursive != CopyRecursive && src != srcPath {
			return nil
		}

		relPath, err := filepath.Rel(src, srcPath)
		if err != nil {
			return fmt.Errorf("error determining the relative path between %s and %s: %w", src, srcPath, err)
		}
		destPath := filepath.Join(dest, relPath)

		if srcInfo.IsDir() {
			return os.MkdirAll(destPath, srcInfo.Mode())
		}

		return copyFile(srcPath, destPath, opts)
	})
}

func copyFile(src string, dest string, opts CopyOption) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	// Check if we should skip existing files
	overwrite := opts&CopyNoOverwrite != CopyNoOverwrite
	createOpts := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if !overwrite { // Return an error if the file exists
		createOpts |= os.O_EXCL
	}

	destF, err := os.OpenFile(dest, createOpts, srcInfo.Mode())
	if err != nil {
		if os.IsExist(err) && !overwrite {
			return nil
		}
		return err
	}
	defer destF.Close()

	_, err = io.Copy(destF, srcF)
	if err != nil {
		return fmt.Errorf("error copying %s to %s: %w", src, dest, err)
	}
	return destF.Close()
}
