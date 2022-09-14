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

package tags

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"knative.dev/test-infra/pkg/logging"
	"knative.dev/test-infra/tools/go-ls-tags/files"
)

func (l Lister) filterIgnored(ctx context.Context, tags []string) ([]string, error) {
	log := logging.FromContext(ctx)
	ignoreFile := l.IgnoreFile
	if !path.IsAbs(l.IgnoreFile) {
		ignoreFile = path.Join(l.Directory, l.IgnoreFile)
	}
	if _, err := os.Stat(ignoreFile); errors.Is(err, os.ErrNotExist) {
		return tags, nil
	}
	log.Infof("Using ignore file: %s", ignoreFile)
	ignored := make([]string, 0, len(tags))
	err := files.ReadLines(ctx, ignoreFile, func(line string) error {
		line = strings.Trim(line, " \t")
		if line == "" || strings.HasPrefix(line, "#") {
			return nil
		}
		ignored = append(ignored, line)
		return nil
	})
	if err != nil {
		return nil, errwrap(err)
	}

	log.Infof("Ignoring tags: %q", ignored)
	result := make([]string, 0, len(tags))
OUTER:
	for _, tag := range tags {
		for _, ignore := range ignored {
			if tag == ignore {
				continue OUTER
			}
		}
		result = append(result, tag)
	}
	return result, nil
}
