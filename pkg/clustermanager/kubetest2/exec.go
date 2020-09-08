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

package kubetest2

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func runWithOutput(command *exec.Cmd) (string, error) {
	var buf bytes.Buffer
	command.Stdout = io.MultiWriter(&buf, os.Stdout)
	command.Stderr = io.MultiWriter(&buf, os.Stderr)
	err := command.Run()
	return buf.String(), err
}
