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

import "istio.io/test-infra/tools/prowgen/pkg/spec"

const (
	testgridDashboardAnnotation   = "testgrid-dashboards"
	testgridDashboardTabAnnoation = "testgrid-tab-name"
)

// addAnnotations adds extra annotations for generating TestGrid config.
func addAnnotations(jobsConfig spec.JobsConfig) spec.JobsConfig {
	for i, job := range jobsConfig.Jobs {
		if hasPeriodic(job.Types) {
			if job.Annotations == nil {
				job.Annotations = map[string]string{}
			}

			if jobsConfig.Branches[0] == "main" {
				// main branch Prow jobs need a dashboard for each repo.
				job.Annotations[testgridDashboardAnnotation] = jobsConfig.Repo
				job.Annotations[testgridDashboardTabAnnoation] = job.Name
			} else {
				// release branch Prow jobs are aggregated under a single dashboard.
				job.Annotations[testgridDashboardAnnotation] = jobsConfig.Org + "-" + jobsConfig.Branches[0]
				job.Annotations[testgridDashboardTabAnnoation] = jobsConfig.Repo + "-" + job.Name
			}
		}
		jobsConfig.Jobs[i] = job
	}
	return jobsConfig
}
