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

package config

import (
	"os"
	"testing"
)

func TestProwConfigPathsExist(t *testing.T) {
	pathsToCheck := [][]string{ProdProwConfigPaths, StagingProwKeyConfigPaths, {ProdTestgridConfigPath}}
	checkPaths(pathsToCheck, t)
}

func TestProwKeyConfigPathsExist(t *testing.T) {
	pathsToCheck := [][]string{ProdProwKeyConfigPaths, StagingProwKeyConfigPaths}
	checkPaths(pathsToCheck, t)
}

func checkPaths(pathsArr [][]string, t *testing.T) {
	t.Helper()
	for _, paths := range pathsArr {
		for _, p := range paths {
			info, err := os.Stat(p)
			if os.IsNotExist(err) || !info.IsDir() {
				t.Fatalf("Expected %q to be a dir, but it's not: %v", p, err)
			}
		}
	}
}
