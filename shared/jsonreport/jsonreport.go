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

package jsonreport

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/knative/test-infra/shared/common"
	"github.com/knative/test-infra/shared/prow"
)

const (
	filename = "flaky-tests.json"
)

type Report struct {
	Repo  string   `json:"-"`
	Flaky []string `json:"flaky"`
}

func (r *Report) WriteToArtifactsDir() error {
	artifactsDir := prow.GetLocalArtifactsDir()
	err := common.CreateDir(path.Join(artifactsDir, r.Repo))
	if nil != err {
		return err
	}
	outFilePath := path.Join(artifactsDir, r.Repo, filename)
	contents, err := json.Marshal(r)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(outFilePath, contents, 0644)
}

func GetReportForRepo(repo string) (*Report, error) {
	var report Report
	artifactsDir := prow.GetLocalArtifactsDir()
	content, err := ioutil.ReadFile(path.Join(artifactsDir, repo, filename))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func NewReport(repo string, flaky []string) *Report {
	return &Report{
		Repo:  repo,
		Flaky: flaky,
	}
}
