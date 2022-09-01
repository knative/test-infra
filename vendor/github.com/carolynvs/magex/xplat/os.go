package xplat

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// InPath determines if the path is in the PATH environment variable.
func InPath(value string) bool {
	pathSep := string(os.PathSeparator)
	pathListSep := string(os.PathListSeparator)
	value = strings.TrimRight(value, pathSep)

	path := os.Getenv("PATH")
	paths := strings.Split(path, pathListSep)
	for _, p := range paths {
		p = strings.TrimRight(p, pathSep)
		if p == value {
			return true
		}
	}

	return false
}

// EnsureInPath adds the specified path to the beginning of the PATH environment
// variable when it is not already in PATH. Detects if this is an Azure CI build
// and exports the updated PATH.
func EnsureInPath(value string) {
	if !InPath(value) {
		PrependPath(value)
	}
}

// PrependPath adds the specified path to the beginning of the PATH environment
// variable. Detects if this is an Azure CI build and exports the updated PATH.
func PrependPath(value string) {
	path := os.Getenv("PATH")
	sep := string(os.PathListSeparator)

	path = fmt.Sprintf("%s%s%s", value, sep, path)
	os.Setenv("PATH", path)
	log.Printf("Added %s to $PATH\n", value)

	isAzureCI := os.Getenv("TF_BUILD")
	if ok, _ := strconv.ParseBool(isAzureCI); ok {
		fmt.Printf("##vso[task.prependpath]%s\n", value)
	}
}
