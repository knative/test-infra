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
	"strings"

	"golang.org/x/mod/modfile"
)

// Modules returns a map of given given modules to their dependencies, and a
// list of unique dependencies.
func Modules(gomod []string, domain string) (map[string][]string, []string, error) {
	if len(gomod) == 0 {
		return nil, nil, errors.New("no go module files provided")
	}

	packages := make(map[string][]string, 0)
	dependencies := make([]string, 0)
	cache := make(map[string]bool)
	for _, gm := range gomod {
		name, pkgs, err := Module(gm, domain)
		if err != nil {
			return nil, nil, err
		}
		packages[name] = pkgs
		for _, pkg := range pkgs {
			if _, seen := cache[pkg]; seen {
				continue
			}
			cache[pkg] = true
			dependencies = append(dependencies, pkg)
		}
	}

	return packages, dependencies, nil
}

// Module returns the name and a list of dependencies for a given module.
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

	file, err := modfile.Parse(gomod, b, nil)
	if err != nil {
		return "", nil, err
	}

	packages := make([]string, 0)
	for _, r := range file.Require {
		// Look for requirements that have the prefix of domain.
		if strings.HasPrefix(r.Mod.Path, domain) {
			packages = append(packages, r.Mod.Path)
		}
	}

	return file.Module.Mod.Path, packages, nil
}
