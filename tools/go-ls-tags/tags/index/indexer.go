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
	"errors"
	"go/build/constraint"
	"log"
	"sync"

	"github.com/hashicorp/go-multierror"
	"knative.dev/test-infra/tools/go-ls-tags/files"
)

var ErrIndexingFailed = errors.New("indexing failed")

// Index can index files to collect Go tags from them.
type Index struct {
	Files []string
}

// Tags will search files for Go build tags.
func (i *Index) Tags() ([]string, error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(i.Files))
	results := make(chan []string, 1)
	errCh := make(chan error)
	errResult := make(chan error)
	tagCh := make(chan string)
	go collectTags(tagCh, results)
	go collectErrors(errCh, errResult)
	for _, file := range i.Files {
		go lookupTags(file, wg, tagCh, errCh)
	}
	wg.Wait()
	close(tagCh)
	close(errCh)

	tags := <-results
	err := <-errResult
	return tags, err
}

func lookupTags(path string, wg *sync.WaitGroup, tagCh chan<- string, errCh chan<- error) {
	defer wg.Done()
	err := files.ReadLines(path, func(line string) {
		expr, err := constraint.Parse(line)
		if err != nil {
			return
		}
		for _, tag := range extractTags(expr) {
			tagCh <- tag
		}
	})
	if err != nil {
		errCh <- err
	}
}

func collectTags(tagCh <-chan string, results chan<- []string) {
	tags := new(set)
	for tag := range tagCh {
		tags.add(tag)
	}

	results <- tags.list()
}

func collectErrors(errs <-chan error, results chan<- error) {
	causes := make([]error, 0)
	for err := range errs {
		causes = append(causes, err)
	}
	var err error
	if len(causes) > 0 {
		err = multierror.Append(ErrIndexingFailed, causes...)
	}
	results <- err
}

func extractTags(expr constraint.Expr) []string {
	tags := make([]string, 0, 1)
	switch v := expr.(type) {
	case *constraint.TagExpr:
		tags = append(tags, v.Tag)
	case *constraint.NotExpr:
		tags = append(tags, extractTags(v.X)...)
	case *constraint.AndExpr:
		tags = append(tags, extractTags(v.X)...)
		tags = append(tags, extractTags(v.Y)...)
	case *constraint.OrExpr:
		tags = append(tags, extractTags(v.X)...)
		tags = append(tags, extractTags(v.Y)...)
	default:
		log.Fatalf("unsupported type %T", v)
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
