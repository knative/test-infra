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

package main

import (
  "fmt"
  "reflect"
	"testing"

	"github.com/google/go-github/github"
	"github.com/TrevorFarrelly/test-infra/shared/ghutil/fakeghutil"
  "github.com/knative/test-infra/tools/monitoring/prowapi"
)

var (
  fakeCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
  The following tests are currently flaky. Running them again to verify...

  Test name | Retries
  --- | ---
  fakejob | 0/3
  fakejob2 | 1/3

  /test fakejob2`
  fakeOrg = "fakeorg"
  fakeRepo = "fakerepo"
  fakePullNumber = 127
  fakeUserID = int64(99)
  fakeUserLogin = "fakelogin"
  fakeUser = &github.User{
    ID: &fakeUserID,
    Login: &fakeUserLogin,
  }
  fakeJob = JobData{
   &prowapi.ReportMessage{
     JobName: "test-job-2",
     Refs: []prowapi.Refs{{
       Org: fakeOrg,
       Repo: fakeRepo,
       Pulls: []prowapi.Pull{{
         Number: fakePullNumber,
       },},
     },},
   },
   nil,
   nil,
  }
  fakeOldComment = &github.IssueComment{
    Body: &fakeCommentBody,
    User: fakeUser,
  }
)


func getFakeGithubClient() *GithubClient {
  gc := fakeghutil.NewFakeGithubClient()
  gc.Repos = []string{fakeRepo}
	gc.User = fakeUser
  return &GithubClient{
    gc,
    *gc.User.Login,
    false,
  }
}

func TestGetOldComment(t *testing.T) {
  gc := getFakeGithubClient()
  cases := []struct{
    client *GithubClient
    org, repo string
    pull int
    wantComment *github.IssueComment
    wantErr error
  }{
    {gc, fakeOrg, fakeRepo, fakePullNumber, nil, nil},
    {gc, fakeOrg, fakeRepo, fakePullNumber, fakeOldComment, nil},
    {gc, fakeOrg, fakeRepo, fakePullNumber, nil, fmt.Errorf("more than one comment on PR")},
  }

  for _, test := range cases {
    actualComment, actualErr := test.client.getOldComment(test.org, test.repo, test.pull)
    if actualComment != test.wantComment {
      t.Errorf("get old comment: got comment %v, want comment %v", actualComment, test.wantComment)
    }
    if !reflect.DeepEqual(actualErr, test.wantErr) {
      t.Errorf("get old comment: got err %v, want err %v", actualErr, test.wantErr)
    }
    // add a comment so we test 0, 1, and 2 comments on the PR respectively
    test.client.CreateComment(test.org, test.repo, test.pull, fakeCommentBody)
  }
}

func TestParseEntries(t *testing.T) {
  cases := []struct{
    input *github.IssueComment
    want map[string]int
  }{
    {fakeOldComment, map[string]int{"fakejob":0, "fakejob2":1}},
  }
  for _, data := range cases {
    actual, _ := parseEntries(data.input)
    if !reflect.DeepEqual(data.want, actual) {
      t.Fatalf("parse entries: got '%v', want '%v'", actual, data.want)
    }
  }
}

func TestBuildNewComment(t *testing.T) {
  cases := []struct {
    jd *JobData
    entries map[string]int
    outliers []string
    want string
  }{

  }
}
