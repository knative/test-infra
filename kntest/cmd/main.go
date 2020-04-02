/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"

	"github.com/spf13/cobra"

	"knative.dev/test-infra/kntest/pkg/cluster"
)

func main() {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "kntest",
		Short: "Tool used in Knative testing, implemented with Go.",
		Run: func(cmd *cobra.Command, args []string) {
			// Print out help info if parent command is run.
			cmd.Help()
		},
	}

	cluster.AddCommands(cmds)

	if err := cmds.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
