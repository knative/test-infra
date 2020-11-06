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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"knative.dev/test-infra/pkg/ghutil"
)

func addActionsCmd(root *cobra.Command) {
	var cmd = &cobra.Command{
		Use:   "actions",
		Short: "Interact with GitHub Actions.",
	}

	addActionsListCmd(cmd)
	addActionsRunCmd(cmd)

	root.AddCommand(cmd)
}

func addActionsListCmd(root *cobra.Command) {
	var (
		tokenPath      string
		query          string
		short          bool
		onlyWorkflowID bool
	)

	var cmd = &cobra.Command{
		Use:   "list org/repo",
		Short: "List GitHub Actions workflows for a given repo.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repos := args

			gh, err := ghutil.NewGithubClient(tokenPath)
			if err != nil {
				return err
			}

			for _, r := range repos {
				or := strings.Split(r, "/")
				if len(or) != 2 {
					return fmt.Errorf("unexpected format %q, expected %q", r, "org/repo")
				}
				org := or[0]
				repo := or[1]

				workflows, err := gh.ListWorkflows(org, repo)
				if err != nil {
					return err
				}
				for _, w := range workflows {
					if !queryByName(w, query) {
						continue
					}

					if short || onlyWorkflowID {
						if onlyWorkflowID {
							_, _ = fmt.Fprintln(cmd.OutOrStdout(), w.GetID())
						} else {
							_, _ = fmt.Fprintln(cmd.OutOrStdout(), w.GetURL())
						}
					} else {
						_, _ = fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("%s (%s) %s", w.GetName(), w.GetState(), w.GetURL()))
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&tokenPath, "token-path", "t", "", "GitHub token file path.")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Search for a workflow by name.")
	cmd.Flags().BoolVar(&short, "short", false, "Short output will only print the workflow url.")
	cmd.Flags().BoolVarP(&onlyWorkflowID, "only-workflow-id", "w", false, "Only output the workflow ID.")

	root.AddCommand(cmd)
}

func addActionsRunCmd(root *cobra.Command) {
	var (
		tokenPath  string
		query      string
		ref        string
		inputs     string
		workflowID int64
		// TODO: interactive inputs based on workflow file config.
	)

	var cmd = &cobra.Command{
		Use:   "run org/repo --query OneResult",
		Short: "Run a GitHub Actions workflow for a given repo.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			or := strings.Split(args[0], "/")
			if len(or) != 2 {
				return fmt.Errorf("unexpected format %q, expected %q", args[0], "org/repo")
			}
			org := or[0]
			repo := or[1]

			gh, err := ghutil.NewGithubClient(tokenPath)
			if err != nil {
				return err
			}

			var jsonInputs map[string]interface{}
			if inputs != "" {
				if err := json.Unmarshal([]byte(inputs), &jsonInputs); err != nil {
					return err
				}
			}

			if workflowID == 0 {
				workflows, err := gh.ListWorkflows(org, repo)
				if err != nil {
					return err
				}
				for _, w := range workflows {
					if !queryByName(w, query) {
						continue
					}

					if workflowID == 0 {
						workflowID = w.GetID()
					} else {
						return fmt.Errorf("query %q matched more than one workflow, cancelling", query)
					}
				}
			}

			if workflowID == 0 {
				return errors.New("unable to locate the workflow requested")
			}

			opts := github.CreateWorkflowDispatchEventRequest{
				Ref:    ref,
				Inputs: jsonInputs,
			}

			_, err = gh.Client.Actions.CreateWorkflowDispatchEvent(context.Background(), org, repo, workflowID, opts)
			return err
		},
	}

	cmd.Flags().StringVarP(&tokenPath, "token-path", "t", "", "GitHub token file path.")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Search for a workflow by name.")
	cmd.Flags().StringVar(&ref, "ref", "master", "Ref to run workflow from.") // This should be the default branch... but for now we use mostly master.
	cmd.Flags().Int64Var(&workflowID, "id", 0, "Workflow ID.")
	cmd.Flags().StringVar(&inputs, "inputs", "", "Workflow inputs.")

	root.AddCommand(cmd)
}

// queryByName returns true if the name of the workflow contains the query.
// Query is case insensitive.
func queryByName(workflow *github.Workflow, query string) bool {
	if query == "" {
		return true
	}
	query = strings.ToLower(query)
	name := strings.ToLower(workflow.GetName())
	return strings.Contains(name, query)
}
