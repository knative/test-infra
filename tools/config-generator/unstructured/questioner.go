/*
Copyright 2021 The Knative Authors

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

package unstructured

import "strings"

// Questioner can be used to look up sub elements of unstructured objects, like
// those created by yaml.Unmarshal or json.Unmarshal.
type Questioner interface {
	// Query will look up sub element, by provided query string. The query string
	// is in format of dot-separated queries like: "foo.bar.42.fizz". In given
	// example, we will be searching for: map value of key "foo", then map value
	// of key "bar", then slice value of index 42, and so on.
	Query(query string) (interface{}, error)
	// Dig will look up sub element, by provided list of Digger's.
	Dig(diggers []Digger) (interface{}, error)
}

// NewQuestioner creates new Questioner object.
func NewQuestioner(un interface{}) Questioner {
	return &defaultQuestioner{un: un}
}

type defaultQuestioner struct {
	un interface{}
}

func (d defaultQuestioner) Query(query string) (interface{}, error) {
	digrs := toDiggers(strings.Split(query, "."))
	return d.Dig(digrs)
}

func (d defaultQuestioner) Dig(diggers []Digger) (interface{}, error) {
	var err error
	un := d.un
	for _, dig := range diggers {
		un, err = dig(un)
		if err != nil {
			return nil, err
		}
	}
	return un, nil
}
