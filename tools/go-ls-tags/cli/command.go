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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"knative.dev/test-infra/tools/go-ls-tags/tags"
)

var ErrUnexpected = errors.New("unexpected")

type execution struct {
	ExecuteContext
	*cobra.Command
	*args
}

func (e execution) invoke() {
	err := e.Command.ExecuteContext(e.ExecuteContext)
	e.ExecuteContext.report(err)
}

func newExecution(ctx ExecuteContext) execution {
	ex := execution{
		ExecuteContext: ctx,
		args:           &args{},
	}
	root := &cobra.Command{
		Use:   "go-ls-tags",
		Short: "List build tags within a Go source tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			lister := &tags.Lister{
				Extension:  ex.extension,
				Exclude:    ex.exclude,
				Directory:  ex.directory,
				IgnoreFile: ex.ignoreFile,
				Context:    cmd.Context(),
			}
			return presenter{cmd}.present(lister.List())
		},
	}
	root.SetErr(ctx.ErrOut)
	root.SetOut(ctx.Out)
	root.SetArgs(ctx.Args)
	ex.Command = withArgs(root, ex.args)
	return ex
}

type presenter struct {
	printer
}

func (p presenter) present(tags []string, err error) error {
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnexpected, err)
	}
	for _, tag := range tags {
		p.Println(tag)
	}
	return nil
}

type printer interface {
	Println(i ...interface{})
}
