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
	"go/build/constraint"
	"log"
	"sync"

	"knative.dev/test-infra/tools/go-ls-tags/files"
)

type Index struct {
	Files []string
	tags  map[string]struct{}
}

func (i *Index) Tags() ([]string, error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(i.Files))
	tags := new(set)
	tagChan := make(chan string)
	go collectTags(tagChan, tags)
	for _, file := range i.Files {
		go lookupTags(file, wg, tagChan)
	}
	wg.Wait()

	return tags.list(), nil
}

func lookupTags(path string, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	files.ReadLines(path, func(line string) {
		expr, err := constraint.Parse(line)
		if err != nil {
			return
		}
		for _, tag := range extractTags(expr) {
			ch <- tag
		}
	})
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

func collectTags(ch <-chan string, tags *set) {
	for tag := range ch {
		tags.add(tag)
	}
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
