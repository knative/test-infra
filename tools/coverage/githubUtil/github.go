/* Copyright 2019 The Knative Authors.

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

package githubUtil

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// convert a file path from profile format to github format.
// Equivalent to remove all path prefix up to and including the repo name
// e.g.
//   knative.dev/$REPO_NAME/pkg/... -> pkg/...
//   github.com/$REPO_OWNER/$REPO_NAME/pkg/... -> pkg/...
func FilePathProfileToGithub(file string) string {
	repoPath, err := GetRepoPath()
	if err != nil {
		log.Fatalf("Cannot get relative repo path: %v", err)
	}
	result := strings.TrimPrefix(file, repoPath+"/")
	if result == file {
		log.Fatalf("repo path (%s) is not a prefix of filepath (%s):", repoPath, file)
	}
	return result
}

// GetRepoPath gets repository path relative to GOPATH/src
func GetRepoPath() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("failed git rev-parse --show-toplevel: '%v'", err)
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("GOPATH is empty")
	}
	relPath, err := filepath.Rel(gopath+"/src", string(out))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(relPath), nil
}
