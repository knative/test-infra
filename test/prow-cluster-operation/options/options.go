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

package options

import (
	"strings"

	"github.com/spf13/cobra"
	"knative.dev/pkg/testutils/clustermanager"
)

type RequestWrapper struct {
	Request          clustermanager.GKERequest
	BackupRegionsStr string
	AddonsStr        string
	NoWait           bool
}

func NewRequestWrapper() RequestWrapper {
	return RequestWrapper{
		Request: clustermanager.GKERequest{},
	}
}

func (rw *RequestWrapper) Prep() {
	if rw.BackupRegionsStr != "" {
		rw.Request.BackupRegions = strings.Split(rw.BackupRegionsStr, ",")
	}
	if rw.AddonsStr != "" {
		rw.Request.Addons = strings.Split(rw.AddonsStr, ",")
	}
}

func AddOptions(cmd *cobra.Command, rw *RequestWrapper) {
	cmd.Flags().Int64Var(&rw.Request.MinNodes, "min-nodes", 0, "minimal number of nodes")
	cmd.Flags().Int64Var(&rw.Request.MaxNodes, "max-nodes", 0, "maximal number of nodes")
	cmd.Flags().StringVar(&rw.Request.NodeType, "node-type", "", "node type")
	cmd.Flags().StringVar(&rw.Request.Region, "region", "", "GCP region")
	cmd.Flags().StringVar(&rw.Request.Zone, "zone", "", "GCP zone")
	cmd.Flags().StringVar(&rw.Request.Project, "project", "", "GCP project")
	cmd.Flags().StringVar(&rw.Request.ClusterName, "name", "", "cluster name")
	cmd.Flags().StringVar(&rw.BackupRegionsStr, "backup-regions", "", "GCP regions as backup, separated by comma")
	cmd.Flags().StringVar(&rw.AddonsStr, "addons", "", "addons to be added, separated by comma")
	cmd.Flags().BoolVar(&rw.Request.SkipCreation, "skip-creation", false, "should skip creation or not")
}
