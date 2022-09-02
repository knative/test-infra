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

package kubetest2

import (
	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
	"knative.dev/test-infra/tools/kntest/pkg/kubetest2/gke"
)

// AddCommand add the command for running kubetest2.
func AddCommand(topLevel *cobra.Command) {
	var kubetest2Cmd = &cobra.Command{
		Use:   "kubetest2",
		Short: "Simple wrapper of kubetest2 commands for Knative testing.",
	}

	kubetest2Options := &kubetest2.Options{}
	addOptions(kubetest2Cmd, kubetest2Options)

	gke.AddCommand(kubetest2Cmd, kubetest2Options)

	topLevel.AddCommand(kubetest2Cmd)
}
