<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> adding go mod parsing
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

package gomod

import (
	"errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"

<<<<<<< HEAD
	"k8s.io/apimachinery/pkg/util/sets"

=======
>>>>>>> adding go mod parsing
	"golang.org/x/mod/modfile"
)

// Modules returns a map of given given modules to their direct dependencies,
// and a list of unique dependencies.
func Modules(gomod []string, domain string) (map[string][]string, []string, error) {
	if len(gomod) == 0 {
		return nil, nil, errors.New("no go module files provided")
	}

<<<<<<< HEAD
<<<<<<< HEAD
	packages := make(map[string][]string, 1)
	cache := make(sets.String, 1)
=======
	packages := make(map[string][]string, 0)
	dependencies := make([]string, 0)
	cache := make(map[string]bool)
>>>>>>> adding go mod parsing
=======
	packages := make(map[string][]string, 1)
	cache := make(sets.String, 1)
>>>>>>> feedback, use sets.
	for _, gm := range gomod {
		name, pkgs, err := Module(gm, domain)
		if err != nil {
			return nil, nil, err
		}
		packages[name] = pkgs
		for _, pkg := range pkgs {
			if cache.Has(pkg) {
				continue
			}
			cache.Insert(pkg)
		}
	}

	return packages, cache.List(), nil
}

<<<<<<< HEAD
// Module returns the name and a list of direct dependencies for a given module.
=======
// Module returns the name and a list of dependencies for a given module.
>>>>>>> adding go mod parsing
// TODO: support url and gopath at some point for the gomod string.
func Module(gomod string, domain string) (string, []string, error) {
	domain = strings.TrimSpace(domain)
	if len(domain) == 0 {
		return "", nil, errors.New("no domain provided")
	}

	b, err := ioutil.ReadFile(gomod)
	if err != nil {
		return "", nil, err
	}

<<<<<<< HEAD

=======
>>>>>>> feedback, use sets.
	file, err := modfile.Parse(gomod, b /*VersionFixer func*/, nil)
	if err != nil {
		return "", nil, err
	}

	packages := make(sets.String, 0)
	for _, r := range file.Require {
<<<<<<< HEAD
		// Do not include indirect dependencies.
		if r.Indirect {
			continue
		}
		// Look for requirements that have the prefix of domain.
		if strings.HasPrefix(r.Mod.Path, domain) && !packages.Has(r.Mod.Path) {
			packages.Insert(r.Mod.Path)
		}
	}

	return file.Module.Mod.Path, packages.List(), nil
=======
		// Look for requirements that have the prefix of domain.
		if strings.HasPrefix(r.Mod.Path, domain) && !packages.Has(r.Mod.Path) {
			packages.Insert(r.Mod.Path)
		}
	}

<<<<<<< HEAD
	return file.Module.Mod.Path, packages, nil
>>>>>>> adding go mod parsing
=======
	return file.Module.Mod.Path, packages.List(), nil
>>>>>>> feedback, use sets.
}
