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

package commands

import "github.com/spf13/cobra"

// New creates a new buoy cli command set.
func New() *cobra.Command {
	var buoyCmd = &cobra.Command{
		Use:   "buoy",
		Short: "Introspect go module dependencies.",
	}

	addFloatCmd(buoyCmd)
	addNeedsCmd(buoyCmd)
	addCheckCmd(buoyCmd)
	addExistsCmd(buoyCmd)

	return buoyCmd
}
