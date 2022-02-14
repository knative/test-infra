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

package gke

import (
	"log"

	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
)

// AddCommand adds gke subcommands.
func AddCommand(kubetest2Cmd *cobra.Command, kubetest2Opts *kubetest2.Options) {
	clusterConfig := &kubetest2.GKEClusterConfig{}

	var gkeCmd = &cobra.Command{
		Use:   "gke",
		Short: "gke related commands for kubetest2.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := kubetest2.Run(kubetest2Opts, clusterConfig); err != nil {
				log.Fatalf("Failed to run tests with kubetest2: %v", err)
			}
		},
	}
	addOptions(gkeCmd, clusterConfig)

	kubetest2Cmd.AddCommand(gkeCmd)
}
