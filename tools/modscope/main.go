package main

import (
	"github.com/wavesoftware/go-commandline"
	"knative.dev/test-infra/tools/modscope/cli"
)

func main() {
	commandline.New(cli.App{}).ExecuteOrDie(cli.Options...)
}

// RunMain is for testing purposes.
func RunMain(opts ...commandline.Option) { //nolint:deadcode
	prev := cli.Options
	cli.Options = append(make([]commandline.Option, 0, len(prev)), opts...)
	defer func(p []commandline.Option) {
		cli.Options = p
	}(prev)
	main()
}
