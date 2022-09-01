// Helper methods for installing Go packages on any platform.
//
// This takes into account both the operating system and the shell environment
// that the process is running within. This enables cross-platform installations
// on MacOS, Windows with WSL, Windows with PowerShell/CMD and Windows with Git
// Bash (mingw).
package pkg
