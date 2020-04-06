/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	gb "go/build"
	"log"
	"time"

	"knative.dev/pkg/test/cmd"
)

// https://golang.org/pkg/cmd/go/internal/modinfo/#ModulePublic
type modInfo struct {
	Path string
	Dir  string
}

type gobuild struct {
	mod *modInfo
}

// moduleInfo returns the module path and module root directory for a project
// using go modules, otherwise returns nil.
//
// Related: https://github.com/golang/go/issues/26504
func moduleInfo() *modInfo {
	output, err := cmd.RunCommand("go list -mod=readonly -m -json")
	if err != nil {
		return nil
	}
	var info modInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil
	}
	return &info
}

// importPackage wraps go/build.Import to handle go modules.
//
// Note that we will fall back to GOPATH if the project isn't using go modules.
func (g *gobuild)importPackage(s string) (*gb.Package, error) {
	if g.mod == nil {
		return gb.Import(s, gb.Default.GOPATH, gb.ImportComment)
	}

	// If we're inside a go modules project, try to use the module's directory
	// as our source root to import:
	// * paths that match module path prefix (they should be in this project)
	// * relative paths (they should also be in this project)
	// if strings.HasPrefix(s, mod.Path) || gb.IsLocalImport(s) {
	now := time.Now()
	gp, err := gb.Import(s, g.mod.Dir, gb.ImportComment)
	then := time.Now()
	log.Printf("time spent: %s", then.Sub(now))
	return gp, err
	// }

	// return nil, errors.New("unmatched importPackage with Go modules")
}
