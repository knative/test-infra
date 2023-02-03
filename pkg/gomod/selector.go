package gomod

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"knative.dev/test-infra/pkg/gowork"
)

var (
	// ErrNoDomain is returned when no domain is provided.
	ErrNoDomain = errors.New("no domain provided")
)

// Matcher matches a module name.
type Matcher func(modname string) bool

// Selector is a set of matchers that can be used to select a subset of modules.
type Selector struct {
	Includes []Matcher
	Excludes []Matcher
}

// Select returns true if the module should be selected.
func (s Selector) Match(modname string) bool {
	for _, match := range s.Excludes {
		if match(modname) {
			return false
		}
	}

	for _, match := range s.Includes {
		if match(modname) {
			return true
		}
	}

	return false
}

// Assert that Selector.Match implements Matcher
var selector_type_assert = Selector{}
var _ Matcher = selector_type_assert.Match

// CurrentModulesMatcher matches the current project modules.
func CurrentModulesMatcher() (Matcher, error) {
	os := gowork.RealSystem{}
	knownMods := sets.NewString()
	mods, err := gowork.List(os, os)
	if err != nil {
		if !errors.Is(err, gowork.ErrInvalidGowork) {
			return nil, err
		}
		m, err := gowork.Current(os, os)
		if err != nil {
			return nil, err
		}
		if m != nil {
			mods = []gowork.Module{*m}
		}
	}

	for _, m := range mods {
		knownMods.Insert(m.Name)
	}

	return knownMods.Has, nil
}

// DefaultSelector returns a selector that includes modules with given domain,
// but excludes the current project modules.
func DefaultSelector(domain string) (Matcher, error) {
	domain = strings.TrimSpace(domain)
	if len(domain) == 0 {
		return func(string) bool { return true }, ErrNoDomain
	}

	currentModules, err := CurrentModulesMatcher()
	if err != nil {
		return nil, err
	}

	return Selector{
		Includes: []Matcher{func(modname string) bool {
			return strings.HasPrefix(modname, domain)
		}},
		Excludes: []Matcher{currentModules},
	}.Match, nil
}
