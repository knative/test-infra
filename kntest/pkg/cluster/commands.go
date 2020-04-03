package cluster

import (
    "github.com/spf13/cobra"

    "knative.dev/test-infra/kntest/pkg/cluster/gke"
)

func AddCommands(topLevel *cobra.Command) {
    var clusterCmd = &cobra.Command{
        Use:   "cluster",
        Short: "Cluster related commands.",
    }

    gke.AddCommands(clusterCmd)
}
