/*
Copyright 2019 The Knative Authors

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

// testgrid.go provides methods to perform action on testgrid.

package testgrid

import "fmt"

const (
	// BaseURL is Knative testgrid base URL
	BaseURL = "https://testgrid.knative.dev"
)

// jobNameTestgridURLMap contains harded coded mapping of job name: Testgrid tab URL relative to base URL
var jobNameTestgridURLMap = map[string]string{
	"ci-knative-serving-continuous":        "serving#continuous",
	"ci-knative-serving-istio-1.5-mesh":    "serving#istio-1.5-mesh",
	"ci-knative-serving-istio-1.5-no-mesh": "serving#istio-1.5-no-mesh",
	"ci-knative-serving-istio-1.4-mesh":    "serving#istio-1.4-mesh",
	"ci-knative-serving-istio-1.4-no-mesh": "serving#istio-1.4-no-mesh",
	"ci-knative-serving-gloo-0.17.1":       "serving#gloo-0.17.1",
	"ci-knative-serving-kourier-stable":    "serving#kourier-stable",
	"ci-knative-serving-contour-latest":    "serving#contour-latest",
	"ci-knative-serving-ambassador-latest": "serving#ambassador-latest",
}

// GetTestgridTabURL gets Testgrid URL for giving job and filters for Testgrid
func GetTestgridTabURL(jobName string, filters []string) (string, error) {
	url, ok := jobNameTestgridURLMap[jobName]
	if !ok {
		return "", fmt.Errorf("cannot find Testgrid tab for job '%s'", jobName)
	}
	for _, filter := range filters {
		url += "&" + filter
	}
	return fmt.Sprintf("%s/%s", BaseURL, url), nil
}
