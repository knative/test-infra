package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"knative.dev/test-infra/pkg/gowork"
)

func listCmd(fl *Flags) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List modules in current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return List(gowork.RealSystem{}, fl, cmd)
		},
	}
}

// List lists modules in current project and prints them using Printer.
func List(os OS, fl *Flags, prt Printer) error {
	pr := presenter{os, fl, prt}
	mods, err := gowork.List(os, os)
	if err != nil && errors.Is(err, gowork.ErrInvalidGowork) {
		var mod *gowork.Module
		mod, err = gowork.Current(os, os)
		var m gowork.Module
		if mod != nil {
			m = *mod
		}
		return pr.presentModule(m, err)
	}
	return pr.presentList(mods, err)
}
