package gomod

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

import (
	"github.com/blang/semver/v4"

	"knative.dev/test-infra/pkg/git"
)

func Float(gomod, release, domain string, ruleset git.RulesetType) ([]string, error) {
	_, packages, err := Modules([]string{gomod}, domain)
	if err != nil {
		return nil, err
	}

	this, err := semver.ParseTolerant(release)

	refs := make([]string, 0)
	for _, pkg := range packages {
		repo, err := moduleToRepo(pkg)
		if err != nil {
			return nil, err
		}

		if ref, refType := repo.BestRefFor(this, ruleset); refType != git.NoRef {
			refs = append(refs, ref)
		}
	}
	return refs, nil
}
