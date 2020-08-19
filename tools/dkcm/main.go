package main

import (
	"flag"

	"knative.dev/test-infra/tools/dkcm/mainservice"
)

func initFlags() *mainservice.Options {
	var o mainservice.Options
	flag.StringVar(&o.Region, "default-region", mainservice.DefaultRegion, "")
	flag.Parse()
	return &o
}

func main() {
	o := initFlags()
	mainservice.Start(o)
}
