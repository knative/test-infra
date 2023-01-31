package cli

import (
	"knative.dev/test-infra/pkg/gowork"
)

// OS represents a virtual operating system.
type OS interface {
	gowork.FileSystem
	gowork.Environment
	Abs(filepath string) string
}

var _ OS = gowork.RealSystem{}
