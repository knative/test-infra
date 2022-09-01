/*
Copyright 2022 The Kubernetes Authors.

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

package mage

import (
	"fmt"
	"log"
	"runtime"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/downloads"
)

const defaultCosignVersion = "v1.7.1"

// EnsureCosign makes sure that the specified cosign version is available
func EnsureCosign(version string) error {
	if version == "" {
		version = defaultCosignVersion
	}

	log.Printf("Checking if `cosign` version %s is installed\n", version)
	found, err := pkg.IsCommandAvailable("cosign", version, "version")
	if err != nil {
		return err
	}

	if !found {
		fmt.Println("`cosign` not found")
		return InstallCosign(version)
	}

	fmt.Println("`cosign` is installed!")
	return nil
}

// InstallCosign installs the required cosign version
func InstallCosign(version string) error {
	fmt.Println("Will install `cosign`")
	target := "cosign"
	if runtime.GOOS == "windows" {
		target = "cosign.exe"
	}

	opts := downloads.DownloadOptions{
		UrlTemplate: "https://github.com/sigstore/cosign/releases/download/{{.VERSION}}/cosign-{{.GOOS}}-{{.GOARCH}}",
		Name:        target,
		Version:     version,
		Ext:         "",
	}

	return downloads.DownloadToGopathBin(opts)
}
