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
	"github.com/spf13/cobra"

	clm "knative.dev/test-infra/pkg/clustermanager/e2e-tests"
)

func addCommonOptions(clusterCmd *cobra.Command, rw *clm.RequestWrapper) {
	pf := clusterCmd.PersistentFlags()
	req := &rw.Request
	// The default values set here are not used in the final operations,
	// they will further be defaulted in
	// https://github.com/knative/pkg/blob/7727cb37e05d6c6dd2abadbc3ab01ab748f12561/testutils/clustermanager/e2e-tests/gke.go#L73-L114
	pf.StringVar(&req.Project, "project", "", "GCP project")
	pf.StringVar(&req.ClusterName, "name", "", "cluster name")
	pf.StringSliceVar(&rw.Regions, "region", []string{}, "GCP regions, separated by comma or multiple args")
	pf.StringVar(&req.ResourceType, "resource-type", "", "Boskos Resource Type")
}

func addCreateOptions(clusterCmd *cobra.Command, rw *clm.RequestWrapper) {
	pf := clusterCmd.Flags()
	req := &rw.Request
	// The default values set here are not used in the final operations,
	// they will further be defaulted in
	// https://github.com/knative/pkg/blob/7727cb37e05d6c6dd2abadbc3ab01ab748f12561/testutils/clustermanager/e2e-tests/gke.go#L73-L114
	pf.Int64Var(&req.MinNodes, "min-nodes", 0, "minimal number of nodes")
	pf.Int64Var(&req.MaxNodes, "max-nodes", 0, "maximal number of nodes")
	pf.StringVar(&req.NodeType, "node-type", "", "node type")
	pf.StringVar(&req.ReleaseChannel, "release-channel", "", "GKE release channel")
	pf.StringVar(&req.GKEVersion, "version", "", "GKE version")
	pf.StringSliceVar(&req.Addons, "addons", []string{}, "addons to be added, separated by comma")
}
