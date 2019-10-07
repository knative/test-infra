/*
Copyright 2019 The Knative Authors

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
	"github.com/spf13/cobra"
	"knative.dev/test-infra/test/prow-cluster-operation/actions"
	"knative.dev/test-infra/test/prow-cluster-operation/options"
)

func AddCommands(topLevel *cobra.Command) {
	addDelete(topLevel)
	addCreate(topLevel)
	addGet(topLevel)
}

func addDelete(topLevel *cobra.Command) {
	sharedOptions := &options.RequestWrapper{}
	subCommand := &cobra.Command{
		Use:   "delete",
		Short: "Delete cluster or release Boskos if running in Prow",
		Run: func(cmd *cobra.Command, args []string) {
			actions.Delete(sharedOptions)
		},
	}
	options.AddOptions(subCommand, sharedOptions)
	topLevel.AddCommand(subCommand)
}

func addCreate(topLevel *cobra.Command) {
	sharedOptions := &options.RequestWrapper{}
	subCommand := &cobra.Command{
		Use:   "create",
		Short: "Create cluster and acquire Boskos if running in Prow, stores metadata under ${ARTIFACT}",
		Run: func(cmd *cobra.Command, args []string) {
			actions.Create(sharedOptions)
		},
	}
	options.AddOptions(subCommand, sharedOptions)
	topLevel.AddCommand(subCommand)
}

func addGet(topLevel *cobra.Command) {
	sharedOptions := &options.RequestWrapper{}
	subCommand := &cobra.Command{
		Use:   "get",
		Short: "Get cluster based on input, stores metadata under ${ARTIFACT}",
		Run: func(cmd *cobra.Command, args []string) {
			actions.Get(sharedOptions)
		},
	}
	options.AddOptions(subCommand, sharedOptions)
	topLevel.AddCommand(subCommand)
}
