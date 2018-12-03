package test

import (
	"os"
	"path"
)

const projectPathLessGoPath = "src/github.com/knative/test-infra/tools/coverage"

func ProjDir() string {
	gopath := os.Getenv("GOPATH")
	return path.Join(gopath, projectPathLessGoPath)
}

func AbsPath(pathRelToProj string) string {
	return path.Join(ProjDir(), pathRelToProj)
}
