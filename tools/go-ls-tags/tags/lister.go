package tags

/*
Copyright 2022 The Knative Authors

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
	"context"
	"sort"

	"knative.dev/test-infra/pkg/logging"
	"knative.dev/test-infra/tools/go-ls-tags/tags/index"
)

const (
	// DefaultExtension is a default extension for Go source files.
	DefaultExtension = "go"
	// DefaultIgnoreFile is a default name of ignore file, that can be used to
	// ignore Go build tags.
	DefaultIgnoreFile = ".gotagsignore"
)

// DefaultExcludes is list of default excludes to be used.
var DefaultExcludes = []string{"vendor", "third_party", "hack", ".git"}

// Lister can list tags of Go source files.
type Lister struct {
	Directory  string
	IgnoreFile string
	Extension  string
	Exclude    []string
	Sort       bool
}

// List all used Go build tags, or errors.
func (l Lister) List(ctx context.Context) ([]string, error) {
	log := logging.FromContext(ctx)
	ff, err := l.files(ctx)
	if err != nil {
		return nil, err
	}
	ix := index.Index{Files: ff}
	var tags []string
	tags, err = ix.Tags(ctx)
	if err != nil {
		return nil, errwrap(err)
	}
	tags, err = l.filterIgnored(ctx, tags)
	if l.Sort {
		sort.Strings(tags)
	}
	log.Infof("Found tags: %q", tags)
	return tags, err
}
