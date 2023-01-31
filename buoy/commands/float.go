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

package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/git"
	"knative.dev/test-infra/pkg/gomod"
)

func addFloatCmd(root *cobra.Command) {
	var (
		domain        string
		release       string
		moduleRelease string
		rulesetFlag   string
		ruleset       git.RulesetType
	)

	var cmd = &cobra.Command{
		Use:   "float go.mod",
		Short: "Find latest versions of dependencies based on a release.",
		Long: `
The goal of the float command is to find the best reference for a given release.
Float will select a ref for found dependencies, in this order (for the Any
ruleset, default):

1. A release tag with matching major and minor; choosing the one with the
   highest patch version, ex: "v0.1.2"
2. If no tags, choose the release branch, ex: "release-0.1"
3. Finally, the default branch, ex: "master"

The selection process for float can be modified by providing a ruleset.

Rulesets,
  Any              tagged releases, release branches, default branch
  Release          tagged releases
  Branch           release branches
  ReleaseOrBranch  tagged releases, release branch

For rulesets that that restrict the selection process, no ref is selected.
`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validation
			ruleset = git.Ruleset(rulesetFlag)
			if ruleset == git.InvalidRule {
				return fmt.Errorf("invalid ruleset, please select one of: [%s]", strings.Join(git.Rulesets(), ", "))
			}
			if moduleRelease == "" {
				moduleRelease = release
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			gomodFile := args[0]

			selector, err := gomod.DefaultSelector(domain)
			if err != nil {
				return err
			}
			refs, err := gomod.Float(gomodFile, release, moduleRelease, selector, ruleset)
			if err != nil {
				return err
			}

			for _, r := range refs {
				if r != "" {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), r)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter (i.e. knative.dev) [required]")
	cmd.Flags().StringVarP(&release, "release", "r", "", "release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]")
	_ = cmd.MarkFlagRequired("release")
	cmd.Flags().StringVarP(&moduleRelease, "module-release", "m", "", "if the go modules are a different release set than the release, use --module-release, should be '<major>.<minor>' (i.e.: 0.12 or v0.12)")
	cmd.Flags().StringVar(&rulesetFlag, "ruleset", git.AnyRule.String(), fmt.Sprintf("The ruleset to evaluate the dependency refs. Rulesets: [%s]", strings.Join(git.Rulesets(), ", ")))

	root.AddCommand(cmd)
}
