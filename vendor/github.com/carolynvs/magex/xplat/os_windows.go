// +build windows

package xplat

// FileExt returns the default file extension based on the operating system.
func FileExt() string {
	return ".exe"
}
