/*
Copyright 2022 The Knative Authors

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

package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"knative.dev/test-infra/tools/go-ls-tags/tags"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

type App struct{}

func (a App) Command() *cobra.Command {
	arg := &args{}
	root := &cobra.Command{
		Use:           "go-ls-tags",
		Short:         "List build tags within a Go source tree",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			lister := &tags.Lister{
				Extension:  arg.extension,
				Exclude:    arg.exclude,
				Directory:  arg.directory,
				IgnoreFile: arg.ignoreFile,
				Context:    cmd.Context(),
			}
			return presenter{cmd}.present(lister.List())
		},
	}
	root.SetOut(os.Stdout)
	root.SetContext(contextWithLogger())
	return withArgs(root, arg)
}

var _ commandline.CobraProvider = App{}
