/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package git

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

const (
	gitAttrLinguistGenerated = "linguist-generated"
	gitAttrCoverageExcluded  = "coverage-excluded"
)

// hasGitAttr checks git attribute value exist for the file
func hasGitAttr(attr string, fileName string) bool {
	var val bytes.Buffer
	attrCmd := exec.Command("git", "check-attr", attr, "--", fileName)

	attrCmd.Stdout = &val
	attrCmd.Start()
	attrCmd.Wait()

	cleaned := strings.TrimSpace(val.String())
	n := len(cleaned)
	// TODO: was only checking "true" before this; I doubt there's ever a
	// situation where it returns "true" instead of "set"
	return cleaned[n-4:n] == "true" || cleaned[n-3:n] == "set"
}

func IsCoverageSkipped(filePath string) bool {
	if hasGitAttr(gitAttrLinguistGenerated, filePath) {
		log.Println("Skipping as file is linguist-generated: ", filePath)
		return true
	} else if hasGitAttr(gitAttrCoverageExcluded, filePath) {
		log.Println("Skipping as file is coverage-excluded: ", filePath)
		return true
	}
	return false
}
