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
	"runtime"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/archive"
	"github.com/carolynvs/magex/pkg/downloads"
)

const defaultKoVersion = "0.11.2"

// EnsureKO
func EnsureKO(version string) error {
	if version == "" {
		version = defaultKoVersion
	}

	fmt.Printf("Checking if `ko` version %s is installed\n", version)
	found, err := pkg.IsCommandAvailable("ko", version, "version")
	if err != nil {
		return err
	}

	if !found {
		fmt.Println("`ko` not found")
		return InstallKO(version)
	}

	fmt.Println("`ko` is installed!")
	return nil
}

// Maybe we can  move this to release-utils
func InstallKO(version string) error {
	fmt.Println("Will install `ko`")
	target := "ko"
	if runtime.GOOS == "windows" {
		target = "ko.exe"
	}

	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/google/ko/releases/download/v{{.VERSION}}/ko_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}{{.EXT}}",
			Name:        "ko",
			Version:     version,
			OsReplacement: map[string]string{
				"darwin":  "Darwin",
				"linux":   "Linux",
				"windows": "Windows",
			},
			ArchReplacement: map[string]string{
				"amd64": "x86_64",
			},
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".tar.gz",
		},
		TargetFileTemplate: target,
	}

	return archive.DownloadToGopathBin(opts)
}
