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

package git

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Repo is a simplified git remote, containing only the list of tags, default
// branch and branches.
type Repo struct {
	Ref           string
	DefaultBranch string
	Tags          []string
	Branches      []string
}

// GetRepo will fetch a git repo and process it into a Repo object.
func GetRepo(ref, url string) (*Repo, error) {
	repo := new(Repo)
	repo.Ref = ref

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ref := range refs {
		if ref.Name().IsTag() {
			repo.Tags = append(repo.Tags, ref.Name().Short())
		} else if ref.Name().IsBranch() {
			repo.Branches = append(repo.Branches, ref.Name().Short())
		} else if ref.Name() == "HEAD" { // Default branch.
			repo.DefaultBranch = ref.Target().Short()
		}
	}

	return repo, nil
}

type RefType int

const (
	// DefaultBranchRef - default branch
	DefaultBranchRef RefType = iota
	// ReleaseBranchRef - release branch
	ReleaseBranchRef
	// ReleaseRef - tagged release
	ReleaseRef
	// NoRef - ref not found
	NoRef
)

var refTypeString = []string{"Default Branch", "Release Branch", "Release", "No Ref"}

// String returns the string of RefType in human readable form.
func (rt RefType) String() string {
	if rt >= DefaultBranchRef && rt <= NoRef {
		return refTypeString[rt]
	}
	return ""
}

// BestRefFor Returns module@ref, isRelease based on the provided ruleset for
// a this release.
func (r *Repo) BestRefFor(this semver.Version, ruleset RulesetType) (string, RefType) {

	switch ruleset {
	case AnyRule, ReleaseOrReleaseBranchRule, ReleaseRule:
		var largest *semver.Version
		// Look for a release.
		for _, t := range r.Tags {
			if sv, ok := normalizeTagVersion(t); ok {
				v, _ := semver.Make(sv)
				// Go does not understand how to fetch semver tags with pre or build tags, skip those.
				if v.Pre != nil || v.Build != nil {
					continue
				}
				if v.Major == this.Major && v.Minor == this.Minor {
					if largest == nil || largest.LT(v) {
						largest = &v
					}
				}
			}
		}
		if largest != nil {
			return fmt.Sprintf("%s@%s", r.Ref, tagVersion(*largest)), ReleaseRef
		}
	}

	switch ruleset {
	case AnyRule, ReleaseOrReleaseBranchRule, ReleaseBranchRule:
		var largest *semver.Version
		// Look for a release branch.
		for _, b := range r.Branches {
			if bv, ok := normalizeBranchVersion(b); ok {
				v, _ := semver.Make(bv)

				if v.Major == this.Major && v.Minor == this.Minor {
					if largest == nil || largest.LT(v) {
						largest = &v
					}
				}
			}
		}
		if largest != nil {
			return fmt.Sprintf("%s@%s", r.Ref, branchVersion(*largest)), ReleaseBranchRef
		}
	}

	switch ruleset {
	case AnyRule:
		// Look for a Return default branch.
		return fmt.Sprintf("%s@%s", r.Ref, r.DefaultBranch), DefaultBranchRef
	}

	// No ref found with the provided rule
	return r.Ref, NoRef
}

func normalizeTagVersion(v string) (string, bool) {
	if strings.HasPrefix(v, "v") {
		// No need to account for unicode widths.
		return v[1:], true
	}
	return v, false
}

func tagVersion(v semver.Version) string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func normalizeBranchVersion(v string) (string, bool) {
	if strings.HasPrefix(v, "release-") {
		// No need to account for unicode widths.
		return v[len("release-"):] + ".0", true
	}
	return v, false
}

func branchVersion(v semver.Version) string {
	return fmt.Sprintf("release-%d.%d", v.Major, v.Minor)
}
