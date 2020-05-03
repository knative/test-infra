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
	clm "knative.dev/pkg/testutils/clustermanager/e2e-tests"

	"knative.dev/test-infra/kntest/pkg/cluster/gke/ops"
)

// AddCommands adds gke subcommands.
func AddCommands(clusterCmd *cobra.Command) {
	var gkeCmd = &cobra.Command{
		Use:   "gke",
		Short: "gke related commands.",
	}

	rw := &ops.RequestWrapper{
		Request: clm.GKERequest{},
	}
	addCommonOptions(gkeCmd, rw)
	addCreate(gkeCmd, rw)
	addDelete(gkeCmd, rw)
	addGet(gkeCmd, rw)
	clusterCmd.AddCommand(gkeCmd)
}

func addCreate(cc *cobra.Command, rw *ops.RequestWrapper) {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a GKE cluster.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			regions := rw.Regions
			if len(regions) != 0 {
				rw.Request.Region = regions[0]
			}
			if len(regions) > 1 {
				rw.Request.BackupRegions = regions[1:]
			}
			if _, err := rw.Create(); err != nil {
				log.Fatalf("error creating the cluster: %v", err)
			}
		},
	}
	addCreateOptions(createCmd, rw)
	cc.AddCommand(createCmd)
}

func addDelete(clusterCmd *cobra.Command, rw *ops.RequestWrapper) {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the current GKE cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := rw.Delete(); err != nil {
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
			if _, err := rw.Get(); err != nil {
				log.Fatalf("error getting the cluster: %v", err)
			}
		},
	}
	clusterCmd.AddCommand(getCmd)
}
