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

	"sigs.k8s.io/release-utils/command"
)

const (
	gitConfigNameKey    = "user.name"
	gitConfigNameValue  = "releng-ci-user"
	gitConfigEmailKey   = "user.email"
	gitConfigEmailValue = "nobody@k8s.io"
)

func CheckGitConfigExists() bool {
	userName := command.New(
		"git",
		"config",
		"--global",
		"--get",
		gitConfigNameKey,
	)

	stream, err := userName.RunSilentSuccessOutput()
	if err != nil || stream.OutputTrimNL() == "" {
		// NB: We're intentionally ignoring the error here because 'git config'
		// returns an error (result code -1) if the config doesn't exist.
		return false
	}

	userEmail := command.New(
		"git",
		"config",
		"--global",
		"--get",
		gitConfigEmailKey,
	)

	stream, err = userEmail.RunSilentSuccessOutput()
	if err != nil || stream.OutputTrimNL() == "" {
		// NB: We're intentionally ignoring the error here because 'git config'
		// returns an error (result code -1) if the config doesn't exist.
		return false
	}

	return true
}

func EnsureGitConfig() error {
	exists := CheckGitConfigExists()
	if exists {
		return nil
	}

	if err := command.New(
		"git",
		"config",
		"--global",
		gitConfigNameKey,
		gitConfigNameValue,
	).RunSuccess(); err != nil {
		return fmt.Errorf("configuring git %s: %w", gitConfigNameKey, err)
	}

	if err := command.New(
		"git",
		"config",
		"--global",
		gitConfigEmailKey,
		gitConfigEmailValue,
	).RunSuccess(); err != nil {
		return fmt.Errorf("configuring git %s: %w", gitConfigEmailKey, err)
	}

	return nil
}
