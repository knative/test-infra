package index

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
	"errors"
	"fmt"
	"go/build/constraint"
	"strings"
	"sync"

	"knative.dev/test-infra/pkg/logging"
	"knative.dev/test-infra/tools/go-ls-tags/files"
)

// ErrIndexingFailed when indexing failed.
var ErrIndexingFailed = errors.New("indexing failed")

// Index can index files to collect Go tags from them.
type Index struct {
	Files []string
}

// Tags will search files for Go build tags.
func (i *Index) Tags(ctx context.Context) ([]string, error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(i.Files))
	results := make(chan result, 1)
	ch := make(chan result)
	go collect(ctx, ch, results)
	for _, file := range i.Files {
		go lookupTags(ctx, file, wg, ch)
	}
	wg.Wait()
	close(ch)

	r := <-results
	return r.tags, r.err
}

type result struct {
	tags []string
	err  error
}

func lookupTags(ctx context.Context, path string, wg *sync.WaitGroup, ch chan<- result) {
	log := logging.FromContext(ctx)
	log.Debugf("Lookup tags for file %s", path)
	defer wg.Done()
	err := files.ReadLines(ctx, path, func(line string) error {
		expr, err := constraint.Parse(line)
		if err != nil {
			return nil
		}
		ch <- result{tags: extractTags(expr)}
		return files.SkipRest
	})
	if err != nil {
		ch <- result{err: err}
	}
}

func collect(ctx context.Context, ch <-chan result, resultCh chan<- result) {
	tags := new(set)
	causes := make([]error, 0)
	log := logging.FromContext(ctx)
	for res := range ch {
		log.Infof("result: %#v", res)
		if res.err != nil {
			causes = append(causes, res.err)
		}
		for _, tag := range res.tags {
			tags.add(tag)
		}
	}
	var err error
	if len(causes) > 0 {
		sb := new(strings.Builder)
		for _, cause := range causes {
			sb.WriteString("\n - " + cause.Error())
		}
		err = fmt.Errorf("%w: %s", ErrIndexingFailed, sb)
	}

	resultCh <- result{tags.list(), err}
}

func extractTags(expr constraint.Expr) []string {
	tags := make([]string, 0, 1)
	switch v := expr.(type) {
	case *constraint.TagExpr:
		tags = append(tags, v.Tag)
	case *constraint.NotExpr:
		return tags
	case *constraint.AndExpr:
		tags = append(tags, extractTags(v.X)...)
		tags = append(tags, extractTags(v.Y)...)
	case *constraint.OrExpr:
		tags = append(tags, extractTags(v.X)...)
		tags = append(tags, extractTags(v.Y)...)
	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}
	return tags
}

type set struct {
	elements map[string]struct{}
}

func (s *set) list() []string {
	s.init()
	val := make([]string, 0, len(s.elements))
	for element := range s.elements {
		val = append(val, element)
	}
	return val
}

func (s *set) add(val string) {
	s.init()
	s.elements[val] = struct{}{}
}

func (s *set) init() {
	if s.elements == nil {
		s.elements = make(map[string]struct{})
	}
}
