package mgx

import "github.com/magefile/mage/mg"

// Must stops the build when an error occurs.
func Must(err error) {
	if err != nil {
		panic(mg.Fatal(1, err))
	}
}
