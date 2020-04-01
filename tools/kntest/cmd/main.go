package main

import (
    "log"

    "github.com/spf13/cobra"
)

func main() {
    // Parent command to which all subcommands are added.
    cmds := &cobra.Command{
        Use:   "kntest",
        Short: "Tool used in Knative testing, implemented with Go.",
        Run: func(cmd *cobra.Command, args []string) {
            cmd.Help()
        },
    }

    // TODO(chizhg): add subcommands.

    if err := cmds.Execute(); err != nil {
        log.Fatalf("error during command execution: %v", err)
    }
}
