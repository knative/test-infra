package files

import "os"

func WorkingDirectoryOrDie() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
