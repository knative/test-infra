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
	"github.com/spf13/cobra"

	"knative.dev/test-infra/kntest/pkg/cluster/ops"
)

func addOptions(clusterCmd *cobra.Command, rw *ops.RequestWrapper) {
	pf := clusterCmd.PersistentFlags()
	req := &rw.Request
	pf.Int64Var(&req.MinNodes, "min-nodes", 0, "minimal number of nodes")
	pf.Int64Var(&req.MaxNodes, "max-nodes", 0, "maximal number of nodes")
	pf.StringVar(&req.NodeType, "node-type", "", "node type")
	pf.StringVar(&req.Region, "region", "", "GCP region")
	pf.StringVar(&req.Zone, "zone", "", "GCP zone")
	pf.StringVar(&req.Project, "project", "", "GCP project")
	pf.StringVar(&req.ClusterName, "name", "", "cluster name")
	pf.StringVar(&req.ReleaseChannel, "release-channel", "", "GKE release channel")
	pf.StringVar(&req.ResourceType, "resource-type", "", "Boskos Resource Type")
	pf.StringSliceVar(&req.BackupRegions, "backup-regions", []string{}, "GCP regions as backup, separated by comma")
	pf.StringSliceVar(&req.Addons, "addons", []string{}, "addons to be added, separated by comma")
	pf.BoolVar(&rw.Request.SkipCreation, "skip-creation", false, "should skip creation or not")
}
