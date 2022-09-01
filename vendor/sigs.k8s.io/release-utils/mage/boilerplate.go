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
	"github.com/carolynvs/magex/shx"

	kpath "k8s.io/utils/path"
	"sigs.k8s.io/release-utils/command"
)

const (
	// repo-infra (used for boilerplate script)
	defaultRepoInfraVersion = "v0.2.5"
	repoInfraURLBase        = "https://raw.githubusercontent.com/kubernetes/repo-infra"
)

// EnsureBoilerplateScript downloads copyright header boilerplate script, if
// not already present in the repository.
func EnsureBoilerplateScript(version, boilerplateScript string, forceInstall bool) error {
	found, err := kpath.Exists(kpath.CheckSymlinkOnly, boilerplateScript)
	if err != nil {
		return fmt.Errorf(
			"checking if copyright header boilerplate script (%s) exists: %w",
			boilerplateScript, err,
		)
	}

	if !found || forceInstall {
		if version == "" {
			log.Printf(
				"A verify_boilerplate.py version to install was not specified. Using default version: %s",
				defaultRepoInfraVersion,
			)

			version = defaultRepoInfraVersion
		}

		if !strings.HasPrefix(version, "v") {
			return fmt.Errorf(
				"repo-infra version (%s) must begin with a 'v'",
				version,
			)
		}

		if _, err := semver.ParseTolerant(version); err != nil {
			return fmt.Errorf(
				"%s was not SemVer-compliant. Cannot continue.: %w",
				version, err,
			)
		}

		binDir := filepath.Dir(boilerplateScript)
		if err := os.MkdirAll(binDir, 0o755); err != nil {
			return fmt.Errorf("creating binary directory: %w", err)
		}

		file, err := os.Create(boilerplateScript)
		if err != nil {
			return fmt.Errorf("creating file: %w", err)
		}

		defer file.Close()

		installURL, err := url.Parse(repoInfraURLBase)
		if err != nil {
			return fmt.Errorf("parsing URL: %w", err)
		}

		installURL.Path = path.Join(
			installURL.Path,
			version,
			"hack",
			"verify_boilerplate.py",
		)

		installCmd := command.New(
			"curl",
			"-sSfL",
			installURL.String(),
			"-o",
			boilerplateScript,
		)

		err = installCmd.RunSuccess()
		if err != nil {
			return fmt.Errorf("installing verify_boilerplate.py: %w", err)
		}
	}

	if err := os.Chmod(boilerplateScript, 0o755); err != nil {
		return fmt.Errorf("making script executable: %w", err)
	}

	return nil
}

// VerifyBoilerplate runs copyright header checks
func VerifyBoilerplate(version, binDir, boilerplateDir string, forceInstall bool) error {
	if _, err := kpath.Exists(kpath.CheckSymlinkOnly, boilerplateDir); err != nil {
		return fmt.Errorf(
			"checking if copyright header boilerplate directory (%s) exists: %w",
			boilerplateDir, err,
		)
	}

	boilerplateScript := filepath.Join(binDir, "verify_boilerplate.py")

	if err := EnsureBoilerplateScript(version, boilerplateScript, forceInstall); err != nil {
		return fmt.Errorf("ensuring copyright header script is installed: %w", err)
	}

	if err := shx.RunV(
		boilerplateScript,
		"--boilerplate-dir",
		boilerplateDir,
	); err != nil {
		return fmt.Errorf("running copyright header checks: %w", err)
	}

	return nil
}
