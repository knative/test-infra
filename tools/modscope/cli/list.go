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
func List(os OS, fl *Flags, prt Printer) error {
	pr := presenter{os, fl, prt}
	mods, err := modules.List(os, os)
	if err != nil && errors.Is(err, modules.ErrInvalidGowork) {
		var mod *modules.Module
		mod, err = modules.Current(os, os)
		var m modules.Module
		if mod != nil {
			m = *mod
		}
		return pr.presentModule(m, err)
	}
	return pr.presentList(mods, err)
}
