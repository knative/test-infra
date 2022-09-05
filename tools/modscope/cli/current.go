package cli

import (
	"github.com/spf13/cobra"
	"knative.dev/test-infra/tools/modscope/modules"
)

func currentCmd(fl *Flags) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Aliases: []string{"curr"},
		Short:   "Prints current module",
		RunE: func(cmd *cobra.Command, args []string) error {
			return current(system{}, fl, cmd)
		},
	}
}

func current(os OS, fl *Flags, print Printer) error {
	pr := presenter{os, fl, print}
	curr, err := modules.Current(os, os)
	return pr.presentModule(curr, err)
}
