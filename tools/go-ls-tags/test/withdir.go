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

package test

import (
	"log"
	"os"
)

// WithDirectory executes a function with a current working directory set
// to a given directory.
func WithDirectory(dir string, fn func()) {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	err = os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}
	fn()
}
