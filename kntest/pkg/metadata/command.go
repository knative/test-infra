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

package metadata

import (
	"log"

	"github.com/spf13/cobra"
	"knative.dev/pkg/testutils/metahelper/client"
)

func AddCommands(topLevel *cobra.Command) {
	var metadataCmd = &cobra.Command{
		Use:   "metadata",
		Short: "Commands for manipulating metadata.json file in Prow job artifacts.",
	}

	var key string
	metadataCmd.PersistentFlags().StringVar(&key, "key", "", "meta info key")

	// Create with default path of metahelper/client, so that the path is
	// consistent with all other consumers of metahelper/client that run within
	// the same context of this tool
	c, err := client.New("")
	if err != nil {
		log.Fatal(err)
	}
	addSetCommand(metadataCmd, c, &key)
	addGetCommand(metadataCmd, c, &key)
	topLevel.AddCommand(metadataCmd)
}

func addSetCommand(metadataCmd *cobra.Command, c *client.Client, key *string) {
	var value string

	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set the meta info key to the given value.",
		Run: func(cmd *cobra.Command, args []string) {
			if *key == "" {
				log.Fatal("meta info key cannot be empty")
			}

			if err := c.Set(*key, value); err != nil {
				log.Fatalf("error setting meta info for %q=%q: %v", *key, value, err)
			}
		},
	}
	setCmd.Flags().StringVar(&value, "value", "", "meta info value")
	metadataCmd.AddCommand(setCmd)
}

func addGetCommand(metadataCmd *cobra.Command, c *client.Client, key *string) {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the meta info value for the given key.",
		Run: func(cmd *cobra.Command, args []string) {
			if *key == "" {
				log.Fatal("meta info key cannot be empty")
			}

			res, err := c.Get(*key)
			if err != nil {
				log.Fatalf("error getting meta info for %q: %v", *key, err)
			}

			log.Print(res)
		},
	}
	metadataCmd.AddCommand(getCmd)
}
