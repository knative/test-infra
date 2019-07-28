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
	"github.com/knative/test-infra/shared/ghutil/fakeghutil"
	"github.com/knative/test-infra/tools/monitoring/prowapi"
)

var (
	oldCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
The following tests are currently flaky. Running them again to verify...

Test name | Retries
--- | ---
fakejob0 | 0/3
fakejob1 | 1/3

Automatically retrying...
/test fakejob1`
	retryCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
The following tests are currently flaky. Running them again to verify...

Test name | Retries
--- | ---
fakejob0 | 1/3
fakejob1 | 1/3

Automatically retrying...
/test fakejob0`
	noMoreRetriesCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
The following tests are currently flaky. Running them again to verify...

Test name | Retries
--- | ---
fakejob0 | 3/3
fakejob1 | 1/3

Job fakejob0 expended all 3 retries without success.`
	failedShortCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
The following tests are currently flaky. Running them again to verify...

Test name | Retries
--- | ---
fakejob0 | 0/3
fakejob1 | 1/3

Failed non-flaky tests preventing automatic retry of fakejob0:

` + "```\ntest0\ntest1\ntest2\ntest3\n```"
	failedLongCommentBody = `<!--AUTOMATED-FLAKY-RETRYER-->
The following tests are currently flaky. Running them again to verify...

Test name | Retries
--- | ---
fakejob0 | 0/3
fakejob1 | 1/3

Failed non-flaky tests preventing automatic retry of fakejob0:

` + "```\ntest0\ntest1\ntest2\ntest3\ntest4\ntest5\ntest6\ntest7\n```\n\nand 2 more."

	fakeOrg       = "fakeorg"
	fakeRepo      = "fakerepo"
	fakeUserLogin = "fakelogin"
	fakePullID    = 127
	fakeCommentID = int64(1)
	fakeUserID    = int64(99)
	fakeUser      = &github.User{
		ID:    &fakeUserID,
		Login: &fakeUserLogin,
	}
	fakeOldComment = &github.IssueComment{
		ID:   &fakeCommentID,
		Body: &oldCommentBody,
		User: fakeUser,
	}
	fakeJob = JobData{
		&prowapi.ReportMessage{
			JobName: "fakejob0",
			Refs: []prowapi.Refs{{
				Org:  fakeOrg,
				Repo: fakeRepo,
				Pulls: []prowapi.Pull{{
					Number: fakePullID,
				}},
			}},
		},
		nil,
		nil,
	}

	fakeFailedTests = []string{"test0", "test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9"}
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

// note: this test is not a standard table-driven test, each case depends on the previous. We add a comment
// to the PR each iteration so we can cover the 0 comments, 1 comment, and 2 comments cases.
func TestGetOldComment(t *testing.T) {
	gc := getFakeGithubClient()
	cases := []struct {
		client      *GithubClient
		org, repo   string
		pull        int
		wantComment *github.IssueComment
		wantErr     error
	}{
		{gc, fakeOrg, fakeRepo, fakePullID, nil, nil},
		{gc, fakeOrg, fakeRepo, fakePullID, fakeOldComment, nil},
		{gc, fakeOrg, fakeRepo, fakePullID, nil, fmt.Errorf("more than one comment on PR")},
	}

	for _, test := range cases {
		actualComment, actualErr := test.client.getOldComment(test.org, test.repo, test.pull)
		if !reflect.DeepEqual(actualComment, test.wantComment) {
			t.Errorf("get old comment: got comment '%v', want comment '%v'", actualComment, test.wantComment)
		}
		if !reflect.DeepEqual(actualErr, test.wantErr) {
			t.Errorf("get old comment: got err '%v', want err '%v'", actualErr, test.wantErr)
		}
		test.client.CreateComment(test.org, test.repo, test.pull, oldCommentBody)
	}
}

func TestParseEntries(t *testing.T) {
	cases := []struct {
		input *github.IssueComment
		want  map[string]int
	}{
		{fakeOldComment, map[string]int{"fakejob0": 0, "fakejob1": 1}},
	}
	for _, data := range cases {
		actual, _ := parseEntries(data.input)
		if !reflect.DeepEqual(actual, data.want) {
			t.Fatalf("parse entries: got '%v', want '%v'", actual, data.want)
		}
	}
}

func TestBuildNewComment(t *testing.T) {
	cases := []struct {
		jd       *JobData
		entries  map[string]int
		outliers []string
		wantBody string
	}{
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, nil, retryCommentBody},
		{&fakeJob, map[string]int{"fakejob0": 3, "fakejob1": 1}, nil, noMoreRetriesCommentBody},
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, fakeFailedTests[:4], failedShortCommentBody},
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, fakeFailedTests, failedLongCommentBody},
	}

	for _, test := range cases {
		gotBody := buildNewComment(test.jd, test.entries, test.outliers)
		if gotBody != test.wantBody {
			t.Fatalf("build new comment: got body \n'%v'\n, want \n'%v'", gotBody, test.wantBody)
		}
	}
}
