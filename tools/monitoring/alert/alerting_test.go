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

package alert

import (
	"testing"
)

func TestToGcsLink(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "no extra link",
			arg:  "gs://knative-prow/logs/ci-knative-serving-go-coverage/1144047857179824128/",
			want: "gs://knative-prow/logs/ci-knative-serving-go-coverage/1144047857179824128/",
		},
		{
			name: "extra gubernator link",
			arg:  "gs://https://gubernator.knative.dev/build/knative-prow/logs/ci-knative-docs-continuous/1132539579983728640/",
			want: "gs://knative-prow/logs/ci-knative-docs-continuous/1132539579983728640/",
		},
		{
			name: "extra spyglass link",
			arg:  "gs://https://prow.knative.dev/view/gcs/knative-prow/logs/ci-knative-serving-go-coverage/1144047857179824128/",
			want: "gs://knative-prow/logs/ci-knative-serving-go-coverage/1144047857179824128/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toGcsLink(tt.arg)
			if got != tt.want {
				t.Errorf("toGcsLink(%v) = %v, want: %v", tt.arg, got, tt.want)
			}
		})
	}
}
