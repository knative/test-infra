/*
Copyright 2021 The Kubernetes Authors.

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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/shx"

	kpath "k8s.io/utils/path"
	"sigs.k8s.io/release-utils/command"
	"sigs.k8s.io/release-utils/env"
)

const (
	// golangci-lint
	defaultGolangCILintVersion = "v1.45.2"
	golangciCmd                = "golangci-lint"
	golangciConfig             = ".golangci.yml"
	golangciURLBase            = "https://raw.githubusercontent.com/golangci/golangci-lint"
	defaultMinGoVersion        = "1.17"
)

// Ensure golangci-lint is installed and on the PATH.
func EnsureGolangCILint(version string, forceInstall bool) error {
	found, err := pkg.IsCommandAvailable(golangciCmd, version)
	if err != nil {
		return fmt.Errorf(
			"checking if %s is available: %w",
			golangciCmd, err,
		)
	}

	if !found || forceInstall {
		if version == "" {
			log.Printf(
				"A golangci-lint version to install was not specified. Using default version: %s",
				defaultGolangCILintVersion,
			)

			version = defaultGolangCILintVersion
		}

		if !strings.HasPrefix(version, "v") {
			return fmt.Errorf(
				"golangci-lint version (%s) must begin with a 'v'",
				version,
			)
		}

		if _, err := semver.ParseTolerant(version); err != nil {
			return fmt.Errorf(
				"%s was not SemVer-compliant. Cannot continue.: %w",
				version, err,
			)
		}

		installURL, err := url.Parse(golangciURLBase)
		if err != nil {
			return fmt.Errorf("parsing URL: %w", err)
		}

		installURL.Path = path.Join(installURL.Path, version, "install.sh")

		err = gopath.EnsureGopathBin()
		if err != nil {
			return fmt.Errorf("ensuring $GOPATH/bin: %w", err)
		}

		gopathBin := gopath.GetGopathBin()

		installCmd := command.New(
			"curl",
			"-sSfL",
			installURL.String(),
		).Pipe(
			"sh",
			"-s",
			"--",
			"-b",
			gopathBin,
			version,
		)

		err = installCmd.RunSuccess()
		if err != nil {
			return fmt.Errorf("installing golangci-lint: %w", err)
		}
	}

	return nil
}

// RunGolangCILint runs all golang linters
func RunGolangCILint(version string, forceInstall bool, args ...string) error {
	if _, err := kpath.Exists(kpath.CheckSymlinkOnly, golangciConfig); err != nil {
		return fmt.Errorf(
			"checking if golangci-lint config file (%s) exists: %w",
			golangciConfig, err,
		)
	}

	if err := EnsureGolangCILint(version, forceInstall); err != nil {
		return fmt.Errorf("ensuring golangci-lint is installed: %w", err)
	}

	if err := shx.RunV(golangciCmd, "version"); err != nil {
		return fmt.Errorf("getting golangci-lint version: %w", err)
	}

	if err := shx.RunV(golangciCmd, "linters"); err != nil {
		return fmt.Errorf("listing golangci-lint linters: %w", err)
	}

	runArgs := []string{"run"}
	runArgs = append(runArgs, args...)

	if err := shx.RunV(golangciCmd, runArgs...); err != nil {
		return fmt.Errorf("running golangci-lint linters: %w", err)
	}

	return nil
}

func TestGo(verbose bool, pkgs ...string) error {
	return testGo(verbose, "", pkgs...)
}

func TestGoWithTags(verbose bool, tags string, pkgs ...string) error {
	return testGo(verbose, tags, pkgs...)
}

func testGo(verbose bool, tags string, pkgs ...string) error {
	verboseFlag := ""
	if verbose {
		verboseFlag = "-v"
	}

	pkgArgs := []string{}
	if len(pkgs) > 0 {
		for _, p := range pkgs {
			pkgArg := fmt.Sprintf("./%s/...", p)
			pkgArgs = append(pkgArgs, pkgArg)
		}
	} else {
		pkgArgs = []string{"./..."}
	}

	cmdArgs := []string{"test"}
	cmdArgs = append(cmdArgs, verboseFlag)
	if tags != "" {
		cmdArgs = append(cmdArgs, "-tags", tags)
	}
	cmdArgs = append(cmdArgs, pkgArgs...)

	if err := shx.RunV(
		"go",
		cmdArgs...,
	); err != nil {
		return fmt.Errorf("running go test: %w", err)
	}

	return nil
}

// VerifyGoMod runs `go mod tidy` and `git diff --exit-code go.*` to ensure
// all module updates have been checked in.
func VerifyGoMod(scriptDir string) error {
	minGoVersion := env.Default("MIN_GO_VERSION", defaultMinGoVersion)
	if err := shx.RunV(
		"go", "mod", "tidy", fmt.Sprintf("-compat=%s", minGoVersion),
	); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}

	if err := shx.RunV("git", "diff", "--exit-code", "go.*"); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}

	return nil
}

// VerifyBuild builds the project for a chosen set of platforms
func VerifyBuild(scriptDir string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	scriptDir = filepath.Join(wd, scriptDir)

	buildScript := filepath.Join(scriptDir, "verify-build.sh")
	if err := shx.RunV(buildScript); err != nil {
		return fmt.Errorf("running go build: %w", err)
	}

	return nil
}
