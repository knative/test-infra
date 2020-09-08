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
)

func addOptions(kubetest2Cmd *cobra.Command, opts *kubetest2.Options) {
	pf := kubetest2Cmd.PersistentFlags()
	pf.StringVar(&opts.ExtraKubetest2Flags, "extra-kubetest2-flags", "", "extra flags for kubetest2")
	pf.StringVar(&opts.TestCommand, "test-command", "", "test command for running the tests")
	pf.BoolVar(&opts.SaveMetaData, "save-meta-data", true, "whether or not to save cluster info into metadata.json")
}
