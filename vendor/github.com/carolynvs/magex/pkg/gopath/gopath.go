package gopath

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/carolynvs/magex/xplat"
)

// EnsureGopathBin ensures that GOPATH/bin exists and is in PATH.
// Detects if this is an Azure CI build and exports the updated PATH.
func EnsureGopathBin() error {
	gopathBin := GetGopathBin()
	err := os.MkdirAll(gopathBin, 0755)
	if err != nil {
		return fmt.Errorf("could not create GOPATH/bin at %s: %w", gopathBin, err)
	}
	xplat.EnsureInPath(GetGopathBin())
	return nil
}

// GetGopathBin returns GOPATH/bin.
func GetGopathBin() string {
	return filepath.Join(GOPATH(), "bin")
}

// GOPATH returns the current GOPATH.
func GOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
}

// UseTempGopath sets the current GOPATH environment variable to a temporary
// directory, returning a cleanup function which reverts the change.
func UseTempGopath() (error, func()) {
	oldpath := os.Getenv("PATH")
	tmp, err := ioutil.TempDir("", "magex")
	if err != nil {
		return fmt.Errorf("failed to create a temp directory: %w", err), func() {}
	}

	cleanup := func() {
		os.RemoveAll(tmp)
		defer os.Setenv("PATH", oldpath)
		defer os.Setenv("GOPATH", build.Default.GOPATH)
	}

	// Remove actual GOPATH/bin from PATH so the test doesn't accidentally pass because the package was installed before the test was run
	gopathBin := filepath.Join(build.Default.GOPATH, "bin")
	os.Setenv("PATH", strings.ReplaceAll(oldpath, gopathBin, ""))

	// Use temp dir for GOPATH
	os.Setenv("GOPATH", tmp)

	err = EnsureGopathBin()
	if err != nil {
		cleanup()
		return err, nil
	}

	return nil, cleanup
}
