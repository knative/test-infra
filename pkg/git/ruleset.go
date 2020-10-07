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

import "strings"

type RulesetType int

const (
	AnyRule RulesetType = iota
	ReleaseOrReleaseBranchRule
	ReleaseRule
	ReleaseBranchRule
	InvalidRule
)

func (rt RulesetType) String() string {
	return [...]string{"Any", "ReleaseOrBranch", "Release", "Branch", "Invalid"}[rt]
}

func Ruleset(rule string) RulesetType {
	switch strings.ToLower(rule) {
	case strings.ToLower(AnyRule.String()):
		return AnyRule
	case strings.ToLower(ReleaseOrReleaseBranchRule.String()):
		return ReleaseOrReleaseBranchRule
	case strings.ToLower(ReleaseRule.String()):
		return ReleaseRule
	case strings.ToLower(ReleaseBranchRule.String()):
		return ReleaseBranchRule
	default:
		return InvalidRule
	}
}

func Rulesets() []string {
	return []string{
		AnyRule.String(),
		ReleaseOrReleaseBranchRule.String(),
		ReleaseRule.String(),
		ReleaseBranchRule.String(),
	}
}
