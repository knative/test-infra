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

package cli

import "strings"

type presenter struct {
	printer
	joiner string
}

func (p presenter) present(tags []string, err error) error {
	if err != nil {
		p.PrintErr("Error: ", err)
		return err
	}
	p.Println(strings.Join(tags, p.joiner))
	return nil
}

type printer interface {
	Println(i ...interface{})
	PrintErr(i ...interface{})
}
