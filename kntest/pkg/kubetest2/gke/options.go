package gke

import (
	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
)

func addOptions(gkeCmd *cobra.Command, cfg *kubetest2.GKEClusterConfig) {
	f := gkeCmd.Flags()
	f.StringVar(&cfg.Environment, "environment", "prod", "")
	f.StringVar(&cfg.CommandGroup, "command-group", "beta", "")
	f.StringVar(&cfg.GCPProjectID, "gcp-project-id", "", "GCP project ID for creating the cluster")
	f.StringVar(&cfg.Name, "name", "e2e-cls", "")
	f.StringVar(&cfg.Region, "region", "us-central1", "")
	f.StringSliceVar(&cfg.BackupRegions, "backup-regions", []string{"us-west1", "us-east1"}, "")
	f.StringVar(&cfg.Machine, "machine", "e2-standard-4", "")
	f.IntVar(&cfg.MinNodes, "min-nodes", 1, "")
	f.IntVar(&cfg.MaxNodes, "max-nodes", 3, "")
	f.StringVar(&cfg.Network, "network", "e2e-network", "")
	f.StringVar(&cfg.Version, "version", "latest", "")
	f.StringVar(&cfg.Scopes, "scopes", "cloud-platform", "")
	f.StringVar(&cfg.Addons, "addons", "", "")
	f.StringVar(&cfg.PrivateClusterAccessLevel, "private-cluster-access-level", "", "")
	f.StringVar(&cfg.PrivateClusterMasterIPRange, "private-cluster-master-ip-range", "172.16.0.%d/28", "")
}
