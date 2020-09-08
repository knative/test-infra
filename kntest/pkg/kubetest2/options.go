package kubetest2

import (
	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
)

func addOptions(kubetest2Cmd *cobra.Command, opts *kubetest2.Options) {
	pf := kubetest2Cmd.PersistentFlags()
	pf.StringVar(&opts.ExtraKubetest2Flags, "extra-kubetest2-flags", "", "extra flags for kubetest2")
	pf.StringVar(&opts.TestCommand, "test-command", "", "test command for running the tests")
	pf.BoolVar(&opts.SaveMetaData, "save-meta-data", true, "whether or not to save cluster info into metadata.json")
}
