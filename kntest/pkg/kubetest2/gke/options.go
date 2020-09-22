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

	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
)

func addOptions(gkeCmd *cobra.Command, cfg *kubetest2.GKEClusterConfig) {
	f := gkeCmd.Flags()
	f.StringVar(&cfg.GCPServiceAccount, "gcp-service-account", "", "The GCP service account used for creating the cluster.")
	f.StringVar(&cfg.Environment, "environment", "prod", "The GKE environment, must be one of prod, staging, staging2 and test.")
	f.StringVar(&cfg.CommandGroup, "command-group", "beta", "The gcloud command group, must be alpha, beta or empty.")
	f.StringVar(&cfg.GCPProjectID, "gcp-project-id", "", "GCP project ID for creating the cluster")

	f.StringVar(&cfg.Name, "name", "e2e-cls", "The GKE cluster name.")
	f.StringVar(&cfg.Region, "region", "us-central1", "The region to create the GKE cluster.")
	f.StringSliceVar(&cfg.BackupRegions, "backup-regions", []string{"us-west1", "us-east1"}, "The backup regions if the cluster creation runs into stockout issue in the primary region.")
	f.StringVar(&cfg.Machine, "machine", "e2-standard-4", "The machine type for the GKE cluster.")
	f.IntVar(&cfg.MinNodes, "min-nodes", 1, "The minimum number of nodes.")
	f.IntVar(&cfg.MaxNodes, "max-nodes", 3, "The maximum number of nodes.")
	f.StringVar(&cfg.Network, "network", "e2e-network", "The network name for the GKE cluster.")
	f.StringVar(&cfg.ReleaseChannel, "release-channel", "", "The release channel of the GKE cluster, can be one of (empty, rapid, regular, stable). When it's used, --cluster-version must be empty or a valid version in the channel."+
		"Reference: https://cloud.google.com/kubernetes-engine/docs/concepts/release-channels")
	f.StringVar(&cfg.Version, "cluster-version", "latest", "The version of the GKE cluster. If --release-channel is not empty, it must be empty or a valid version in the channel")
	f.StringVar(&cfg.Scopes, "scopes", "cloud-platform", "Scopes for the GKE cluster, should be comma-separated.")
	f.StringVar(&cfg.Addons, "addons", "", "Addons for the GKE cluster, should be comma-separated.")
	f.BoolVar(&cfg.EnableWorkloadIdentity, "enable-workload-identity", false, "Whether to enable workload identity for this cluster or not.")
	f.StringVar(&cfg.PrivateClusterAccessLevel, "private-cluster-access-level", "", "Private cluster access level, if not empty, must be one of 'no', 'limited' or 'unrestricted'")
	f.StringVar(&cfg.PrivateClusterMasterIPSubnetRange, "private-cluster-master-ip-subnet-range", "172.16.0", "The master IP subnet range for the private cluster. The last digit must be left empty to allow retrying cluster creation in the backup regions.")
	f.StringVar(&cfg.PrivateClusterMasterIPSubnetMask, "private-cluster-master-ip-subnet-mask", "28", "The master IP subnet mask for the private cluster.")

	f.StringVar(&cfg.ExtraGcloudFlags, "extra-gcloud-flags", "", "The extra gcloud flags that will be used for cluster creation.")
}
