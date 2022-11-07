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

package files

import (
	"bufio"
	"context"
	"errors"
	"os"
)

var (
	// SkipRest can be used to skip the rest of the file.
	SkipRest = errors.New("skip rest")
)

// ReadLines will read lines from a file and pass each line to a func handler.
func ReadLines(ctx context.Context, path string, fn func(line string) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	return scanFile(ctx, f, fn)
}

func scanFile(ctx context.Context, f *os.File, fn func(line string) error) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(scanner.Text()); err != nil {
				if errors.Is(err, SkipRest) {
					return nil
				}
				return err
			}
		}
	}

	return scanner.Err()
}
