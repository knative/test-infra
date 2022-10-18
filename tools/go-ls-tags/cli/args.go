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
	"github.com/spf13/cobra"
	"knative.dev/test-infra/tools/go-ls-tags/files"
	"knative.dev/test-infra/tools/go-ls-tags/tags"
)

type args struct {
	directory  string
	ignoreFile string
	extension  string
	exclude    []string
	joiner     string
	sort       bool
}

func withArgs(root *cobra.Command, a *args) *cobra.Command {
	pf := root.PersistentFlags()
	pf.StringVar(&a.extension,
		"extension",
		tags.DefaultExtension,
		"A Go file extension")
	pf.StringVar(&a.ignoreFile,
		"ignore-file",
		tags.DefaultIgnoreFile,
		"An ignore file used to filter out tags")
	pf.StringVar(&a.directory,
		"directory",
		files.WorkingDirectoryOrDie(),
		"A directory to start from")
	pf.StringSliceVar(&a.exclude,
		"exclude",
		tags.DefaultExcludes,
		"Directories to exclude")
	pf.StringVar(&a.joiner,
		"joiner",
		"\n",
		"Tags will be joined with this string on output")
	pf.BoolVar(&a.sort,
		"sort",
		true,
		"If true then tags will be sorted",
	)
	return root
}
