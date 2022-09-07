package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

// App is a modscope commandline application.
type App struct{}

// Command returns a cobra command for the application.
func (a App) Command() *cobra.Command {
	fl := Flags{}
	c := &cobra.Command{
		Use:          "modscope",
		Short:        "modscope is a tool to show information about Go modules",
		SilenceUsage: true,
	}
	c.AddCommand(currentCmd(&fl))
	c.AddCommand(listCmd(&fl))
	c.SetOut(os.Stdout)
	return fl.configure(c)
}
