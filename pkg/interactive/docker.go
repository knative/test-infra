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

// Helper functions for running interactive docker CLI sessions from Go
package interactive

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var defaultDockerCommands []string

func init() {
	defaultDockerCommands = []string{"docker", "run", "-it", "--rm", "--entrypoint", "bash"}
}

// Env represents a collection of environment variables and their values
type Env map[string]string

// Docker is mostly an Command preloaded with arguments which setup Docker for running an image interactively.
type Docker struct {
	Command
}

// PromoteFromEnv pulls the named environment variables from the environment and puts them in the Env.
// It does not stop on error and returns an error listing all the failed values
func (e Env) PromoteFromEnv(envVars ...string) error {
	var err error
	for _, env := range envVars {
		v := os.Getenv(env)
		if v == "" {
			err = fmt.Errorf("environment variable %q is not set; %v", env, err)
		} else {
			e[env] = v
		}
	}
	return err
}

// NewDocker creates a Docker with default Docker command arguments for running interactively
func NewDocker() Docker {
	return Docker{NewCommand(defaultDockerCommands...)}
}

// AddEnv adds arguments so all the environment variables present in e become part of the docker run's environment
func (d *Docker) AddEnv(e Env) {
	for k, v := range e {
		d.AddArgs("-e", fmt.Sprintf("%s=%s", k, v))
	}
}

// AddMount add arguments for the --mount command
func (d *Docker) AddMount(typeStr, source, target string, optAdditionalArgs ...string) {
	addl := ""
	if len(optAdditionalArgs) != 0 {
		addl = "," + strings.Join(optAdditionalArgs, ",")
	}
	d.AddArgs("--mount", fmt.Sprintf("type=%s,source=%s,target=%s%s", typeStr, source, target, addl))
}

// CopyAndAddMount copies the source files into a temp directory, and then
// mounts them as the target. It also returns a function to remove the temp
// directory for cleaning up.
func (d *Docker) CopyAndAddMount(typeStr, parentDir, source, target string, optAdditionalArgs ...string) (func(), error) {
	fi, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("error getting the FileInfo for %q: %w", source, err)
	}

	tempDir, err := ioutil.TempDir(parentDir, fi.Name())
	if err != nil {
		return nil, fmt.Errorf("error creating a temporary directory: %w", err)
	}
	if err := os.Chmod(tempDir, 0777); err != nil {
		return nil, fmt.Errorf("error changing file mode for %q: %w", tempDir, err)
	}
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	cmd := NewCommand("cp", "-r", source, tempDir)
	if err := cmd.Run(); err != nil {
		return cleanup, fmt.Errorf("error copying %q to %q: %w", source, tempDir, err)
	}

	if fi.IsDir() {
		d.AddMount(typeStr, tempDir, target, optAdditionalArgs...)
	} else {
		d.AddMount(typeStr, tempDir+string(os.PathSeparator)+fi.Name(), target, optAdditionalArgs...)
	}

	return cleanup, nil
}
