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

package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	prowgenpkg "istio.io/test-infra/tools/prowgen/pkg"
)

// GenerateProwJobsConfig will generate Prow jobs from prowJobsConfigInput, and write
// them to prowJobsConfigOutput.
func GenerateProwJobsConfig(prowJobsConfigInput, prowJobsConfigOutput string) error {

	bc := prowgenpkg.ReadBase(nil, filepath.Join(prowJobsConfigInput, ".base.yaml"))

	if err := filepath.WalkDir(prowJobsConfigInput, func(path string, d os.DirEntry, err error) error {
		log.Printf("Generating Prow jobs for %q", path)
		// Skip directory, base config file and other unrelated files.
		if d.IsDir() || d.Name() == ".base.yaml" || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		baseConfig := bc
		if _, err := os.Stat(filepath.Join(filepath.Dir(path), ".base.yaml")); !os.IsNotExist(err) {
			baseConfig = prowgenpkg.ReadBase(&baseConfig, filepath.Join(filepath.Dir(path), ".base.yaml"))
		}
		cli := &prowgenpkg.Client{
			BaseConfig:          baseConfig,
			LongJobNamesAllowed: true,
		}

		jobsConfig := cli.ReadJobsConfig(path)
		jobsConfig = addSchedule(jobsConfig)
		jobsConfig = addAnnotations(jobsConfig)
		output, err := cli.ConvertJobConfig(path, jobsConfig, jobsConfig.Branches[0])
		if err != nil {
			return fmt.Errorf("error generating Prow jobs config for %q: %w", path, err)
		}

		outputFile := filepath.Join(prowJobsConfigOutput,
			fmt.Sprintf("%s/%s-%s.gen.yaml", jobsConfig.Org, jobsConfig.Repo, jobsConfig.Branches[0]))
		log.Printf("Writing the generated Prow config to %q", outputFile)
		if err := prowgenpkg.Write(output, outputFile, bc.AutogenHeader); err != nil {
			return fmt.Errorf("error writing generated Prow jobs config to %q: %w", outputFile, err)
		}

		return nil

	}); err != nil {
		return fmt.Errorf("error walking dir %q: %w", prowJobsConfigInput, err)
	}

	return nil
}
