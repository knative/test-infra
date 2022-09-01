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

	"github.com/blang/semver"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/shx"
)

const (
	// zeitgeist
	defaultZeitgeistVersion = "v0.3.0"
	zeitgeistCmd            = "zeitgeist"
	zeitgeistModule         = "sigs.k8s.io/zeitgeist"
)

// Ensure zeitgeist is installed and on the PATH.
func EnsureZeitgeist(version string) error {
	if version == "" {
		log.Printf(
			"A zeitgeist version to install was not specified. Using default version: %s",
			defaultZeitgeistVersion,
		)

		version = defaultZeitgeistVersion
	}

	if _, err := semver.ParseTolerant(version); err != nil {
		return fmt.Errorf(
			"%s was not SemVer-compliant, cannot continue: %w",
			version, err,
		)
	}

	if err := pkg.EnsurePackage(zeitgeistModule, version); err != nil {
		return fmt.Errorf("ensuring package: %w", err)
	}

	return nil
}

// VerifyDeps runs zeitgeist to verify dependency versions
func VerifyDeps(version, basePath, configPath string, localOnly bool) error {
	if err := EnsureZeitgeist(version); err != nil {
		return fmt.Errorf("ensuring zeitgeist is installed: %w", err)
	}

	args := []string{"validate"}
	if localOnly {
		args = append(args, "--local")
	}

	if basePath != "" {
		args = append(args, "--base-path", basePath)
	}

	if configPath != "" {
		args = append(args, "--config", configPath)
	}

	if err := shx.RunV(zeitgeistCmd, args...); err != nil {
		return fmt.Errorf("running zeitgeist: %w", err)
	}

	return nil
}

/*
##@ Dependencies

.SILENT: update-deps update-deps-go update-mocks
.PHONY:  update-deps update-deps-go update-mocks

update-deps: update-deps-go ## Update all dependencies for this repo
	echo -e "${COLOR}Commit/PR the following changes:${NOCOLOR}"
	git status --short

update-deps-go: GO111MODULE=on
update-deps-go: ## Update all golang dependencies for this repo
	go get -u -t ./...
	go mod tidy
	go mod verify
	$(MAKE) test-go-unit
	./scripts/update-all.sh

update-mocks: ## Update all generated mocks
	go generate ./...
	for f in $(shell find . -name fake_*.go); do \
		cp scripts/boilerplate/boilerplate.generatego.txt tmp ;\
		cat $$f >> tmp ;\
		mv tmp $$f ;\
	done
*/
