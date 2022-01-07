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

package files

import (
	"bufio"
	"log"
	"os"
)

func ReadLines(path string, fn func(line string)) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fn(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
