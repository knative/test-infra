package xplat

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetMSystem returns current the MSys2 subsystem.
//
// Allowed values are: (empty), mingw32, mingw64, msys2.
func GetMSystem() string {
	return strings.ToLower(os.Getenv("MSYSTEM"))
}

// DetectShell uses environment variables to determine
// the containing shell environment.
//
// Values: unknown, posix, mingw, powershell, cmd, $SHELL (e.g. sh, bash)
//
// This is not guaranteed to always work depending on customizations to the
// default Windows environment variables. It uses the following heuristics:
// * MSYSTEM is set
// * SHELL is set
// * PSModulePath containing user's home directory
func DetectShell() string {
	msystem := GetMSystem()
	if msystem != "" {
		return msystem
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		return strings.ToLower(filepath.Base(shell))
	}

	if psModules := strings.ToLower(os.Getenv("PSModulePath")); psModules != "" {
		home, err := os.UserHomeDir()
		if err == nil {
			home := strings.ToLower(home)
			if strings.Contains(psModules, home) {
				return "powershell"
			}
		}
	}

	// On Windows assume that if it's not MSys2 or PowerShell, then it is
	// cmd.
	if runtime.GOOS == "windows" {
		return "cmd"
	}

	// We are on some posix shell most likely because it's not windows.
	// This is the best we can narrow it down to.
	return "posix"
}
