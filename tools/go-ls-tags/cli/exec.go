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
	"context"
	"io"
	"os"

	"knative.dev/test-infra/pkg/golang/retcode"
)

// OsExitFunc is a os.Exit like function.
type OsExitFunc func(int)

// ExecuteOption is a option that can configure the ExecuteContext.
type ExecuteOption func(ctx *ExecuteContext)

// ExecuteContext holds a execution context.
type ExecuteContext struct {
	OsExitFunc
	Args   []string
	Out    io.Writer
	ErrOut io.Writer
	context.Context
}

// Execute the go-ls-tags tool.
func Execute(opts ...ExecuteOption) {
	ctx := ExecuteContext{
		OsExitFunc: os.Exit,
		Args:       os.Args[1:],
		Out:        os.Stdout,
		ErrOut:     os.Stderr,
		Context:    context.Background(),
	}
	for _, opt := range opts {
		opt(&ctx)
	}
	newExecution(ctx).invoke()
}

func (c ExecuteContext) report(err error) {
	c.OsExitFunc(retcode.Calc(err))
}
