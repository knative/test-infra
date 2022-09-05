package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"knative.dev/test-infra/tools/modscope/modules"
)

func listCmd(fl *Flags) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List modules in current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return List(system{}, fl, cmd)
		},
	}
}

// List lists modules in current project and prints them using Printer.
func List(os OS, fl *Flags, print Printer) error {
	pr := presenter{os, fl, print}
	mods, err := modules.List(os, os)
	if err != nil && errors.Is(err, modules.ErrInvalidGowork) {
		var mod *modules.Module
		mod, err = modules.Current(os, os)
		return pr.presentModule(mod, err)
	}
	return pr.presentList(mods, err)
}
