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
	fakeCommentBodyTemplate = "<!--AUTOMATED-FLAKY-RETRYER-->\n" +
		"The following tests are currently flaky. Running them again to verify...\n\n" +
		"Test name | Retries\n" +
		"--- | ---\n" +
		"fakejob0 | %d/3\n" +
		"fakejob1 | %d/3\n\n" +
		"%s"
	fakeFailedTestsShort = "Failed non-flaky tests preventing automatic retry of fakejob0:\n\n" +
		"```\ntest0\ntest1\ntest2\ntest3\n```"
	fakeFailedTestsLong = "Failed non-flaky tests preventing automatic retry of fakejob0:\n\n" +
		"```\ntest0\ntest1\ntest2\ntest3\ntest4\ntest5\ntest6\ntest7\n```\n\nand 2 more."

	fakeOldCommentBody = fmt.Sprintf(fakeCommentBodyTemplate, 0, 1, "Automatically retrying...\n/test fakejob1")

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
		Body: &fakeOldCommentBody,
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
			t.Errorf("get old comment: got comment %v, want comment %v", actualComment, test.wantComment)
		}
		if !reflect.DeepEqual(actualErr, test.wantErr) {
			t.Errorf("get old comment: got err %v, want err %v", actualErr, test.wantErr)
		}
		// add a comment so we test 0, 1, and 2 comments on the PR respectively
		test.client.CreateComment(test.org, test.repo, test.pull, fakeOldCommentBody)
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
		wantBool bool
	}{
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, nil, fmt.Sprintf(fakeCommentBodyTemplate, 1, 1, "Automatically retrying...\n/test fakejob0"), true},
		{&fakeJob, map[string]int{"fakejob0": 3, "fakejob1": 1}, nil, fmt.Sprintf(fakeCommentBodyTemplate, 3, 1, ""), false},
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, fakeFailedTests[:4], fmt.Sprintf(fakeCommentBodyTemplate, 0, 1, fakeFailedTestsShort), true},
		{&fakeJob, map[string]int{"fakejob0": 0, "fakejob1": 1}, fakeFailedTests, fmt.Sprintf(fakeCommentBodyTemplate, 0, 1, fakeFailedTestsLong), true},
	}

	for _, test := range cases {
		gotBody, gotBool := buildNewComment(test.jd, test.entries, test.outliers)
		if gotBody != test.wantBody {
			t.Fatalf("build new comment: got body \n'%v'\n, want \n'%v'", gotBody, test.wantBody)
		}
		if gotBool != test.wantBool {
			t.Fatalf("build new comment: got bool '%v', want '%v'", gotBool, test.wantBool)
		}
	}
}
