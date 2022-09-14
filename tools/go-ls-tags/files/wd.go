package files

import "os"

// WorkingDirectoryOrDie will get the working directory, or die trying.
func WorkingDirectoryOrDie() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
