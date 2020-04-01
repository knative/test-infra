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

package cluster

import (
	"log"

	"github.com/spf13/cobra"
	clm "knative.dev/pkg/testutils/clustermanager/e2e-tests"

	"knative.dev/test-infra/kntest/pkg/cluster/ops"
)

// AddCommands adds cluster subcommands.
func AddCommands(topLevel *cobra.Command) {
	var clusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Cluster related commands.",
	}

	rw := &ops.RequestWrapper{
		Request: clm.GKERequest{},
		NoWait:  false,
	}
	addOptions(clusterCmd, rw)
	addCreate(clusterCmd, rw)
	addDelete(clusterCmd, rw)
	addGet(clusterCmd, rw)
	topLevel.AddCommand(clusterCmd)
}

func addCreate(cc *cobra.Command, rw *ops.RequestWrapper) {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a GKE cluster.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := ops.Create(rw); err != nil {
				log.Fatalf("error creating the cluster: %v", err)
			}
		},
	}
	cc.AddCommand(createCmd)
}

func addDelete(clusterCmd *cobra.Command, rw *ops.RequestWrapper) {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the current GKE cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := ops.Delete(rw); err != nil {
				log.Fatalf("error deleting the cluster: %v", err)
			}
		},
	}
	clusterCmd.AddCommand(deleteCmd)
}

func addGet(clusterCmd *cobra.Command, rw *ops.RequestWrapper) {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the existing cluster from kubeconfig or gcloud.",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := ops.Get(rw); err != nil {
				log.Fatalf("error getting the cluster: %v", err)
			}
		},
	}
	clusterCmd.AddCommand(getCmd)
}
