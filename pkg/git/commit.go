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

package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"knative.dev/test-infra/pkg/cmd"
	"knative.dev/test-infra/pkg/helpers"
)

func call(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// MakeCommit adds the changed files and create a new Git commit.
func MakeCommit(gi Info, message string, dryrun bool) (bool, error) {
	if gi.Head == "" {
		log.Fatal("pushing to empty branch ref is not allowed")
	}
	var (
		statusCmd = "git status --porcelain"
		addCmd    = "git add -A"
		commitCmd = fmt.Sprintf("git commit -m %q", message)
		pushCmd   = fmt.Sprintf("git push -f git@github.com:%s/%s.git HEAD:%s",
			gi.UserID, gi.Repo, gi.Head)
	)

	changes, err := cmd.RunCommand(statusCmd)
	if err != nil {
		return false, fmt.Errorf("Failed running %q:\nOutput: %q\nError: %v",
			statusCmd, changes, err)
	}
	log.Print("diff output:", changes)
	if strings.TrimSpace(changes) == "" {
		log.Print("No changes to commit, skipping")
		return false, nil
	}

	if gi.UserName != "" && gi.Email != "" {
		commitCmd = fmt.Sprintf("%s --author '%s <%s>'",
			commitCmd, strings.Trim(gi.UserName, "'"), gi.Email)
	}

	cmds := []string{addCmd, commitCmd, pushCmd}

	return true, helpers.Run(
		fmt.Sprintf("Running %q", strings.Join(cmds, "; ")),
		func() error {
			out, err := cmd.RunCommands(cmds...)
			if err != nil {
				return fmt.Errorf("Failed running %q:\nOutput: %q\nError: %v",
					strings.Join(cmds, "; "), out, err)
			}
			return nil
		},
		dryrun,
	)
}
