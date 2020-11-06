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
	"github.com/spf13/cobra"
	"knative.dev/test-infra/pkg/ghutil"
)

func addReposCmd(root *cobra.Command) {

	var tokenPath string

	var cmd = &cobra.Command{
		Use:   "repos org1 [org2 org3...]",
		Short: "List the repos for a given org.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgs := args
			gh, err := ghutil.NewGithubClient(tokenPath)
			if err != nil {
				return err
			}

			// for all given orgs, list the repos.
			for _, org := range orgs {
				repos, err := gh.ListRepos(org)
				if err != nil {
					return err
				}
				for _, repo := range repos {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), org+"/"+repo)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&tokenPath, "token-path", "t", "", "GitHub token file path.")

	root.AddCommand(cmd)
}
