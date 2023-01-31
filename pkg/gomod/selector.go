package gomod

import (
	"log"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"knative.dev/test-infra/pkg/gowork"
)

var (
	// ErrNoDomain is returned when no domain is provided.
	ErrNoDomain = errors.New("no domain provided")
)

// Matcher matches a module name.
type Matcher interface {
	// Match returns true if the module name matches.
	Match(modname string) bool
}

// MatcherFunc is a function that implements the Matcher interface.
func MatcherFunc(fn matcherFunc) Matcher {
	return fn
}

type matcherFunc func(modname string) bool

func (fn matcherFunc) Match(modname string) bool {
	return fn(modname)
}

// Selector is a set of matchers that can be used to select a subset of modules.
type Selector struct {
	Includes []Matcher
	Excludes []Matcher
}

// Select returns true if the module should be selected.
func (s Selector) Select(modname string) bool {
	selected := false
	for _, m := range s.Includes {
		if m.Match(modname) {
			selected = true
			break
		}
	}
	if !selected {
		return false
	}
	for _, m := range s.Excludes {
		if m.Match(modname) {
			return false
		}
	}
	return true
}

// CurrentModulesMatcher matches the current project modules.
type CurrentModulesMatcher struct {
	once sync.Once
	mods []gowork.Module
	err  error
}

func (c *CurrentModulesMatcher) Match(modname string) bool {
	c.once.Do(func() {
		os := gowork.RealSystem{}
		c.mods, c.err = gowork.List(os, os)
		if c.err != nil && errors.Is(c.err, gowork.ErrInvalidGowork) {
			var m *gowork.Module
			m, c.err = gowork.Current(os, os)
			if m != nil {
				c.mods = []gowork.Module{*m}
			}
		}
	})
	if c.err != nil {
		log.Fatal(c.err)
	}
	for _, m := range c.mods {
		if m.Name == modname {
			return true
		}
	}
	return false
}

// DefaultSelector returns a selector that includes modules with given domain,
// but excludes the current project modules.
func DefaultSelector(domain string) (Selector, error) {
	domain = strings.TrimSpace(domain)
	if len(domain) == 0 {
		return Selector{}, ErrNoDomain
	}

	return Selector{
		Includes: []Matcher{MatcherFunc(func(modname string) bool {
			return strings.HasPrefix(modname, domain)
		})},
		Excludes: []Matcher{&CurrentModulesMatcher{}},
	}, nil
}
