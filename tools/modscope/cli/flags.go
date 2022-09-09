package cli

import (
	"github.com/spf13/cobra"
)

// Flags are options for modscope command.
type Flags struct {
	DisplayFilepath bool
}

func (f *Flags) configure(c *cobra.Command) *cobra.Command {
	fl := c.PersistentFlags()
	fl.BoolVarP(&f.DisplayFilepath, "path", "p", false,
		"display module's filepath instead of name")
	return c
}
