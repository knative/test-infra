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

package main

import (
	"github.com/spf13/pflag"

	"knative.dev/test-infra/tools/dkcm/mainservice"
)

func initFlags() *mainservice.Options {
	var o mainservice.Options
	var regions []string
	pflag.StringSliceVar(&regions, "region", []string{}, "")
	pflag.Parse()
	if len(regions) > 0 {
		o.Region = regions[0]
		if len(regions) > 1 {
			o.BackupRegions = regions[1:]
		}
	} else {
		o.Region = mainservice.DefaultRegion
	}
	return &o
}

func main() {
	o := initFlags()
	mainservice.Start(o)
}
