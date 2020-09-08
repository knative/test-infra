package kubetest2

import (
	"github.com/spf13/cobra"

	"knative.dev/test-infra/kntest/pkg/kubetest2/gke"
	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
)

func AddCommands(topLevel *cobra.Command) {
	var kubetest2Cmd = &cobra.Command{
		Use:   "kubetest2",
		Short: "Simple wrapper of kubetest2 commands for Knative testing.",
	}

	kubetest2Options := &kubetest2.Options{}
	addOptions(kubetest2Cmd, kubetest2Options)

	gke.AddCommand(kubetest2Cmd, kubetest2Options)

	topLevel.AddCommand(kubetest2Cmd)
}
