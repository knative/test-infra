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
