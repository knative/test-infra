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
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/tools/flaky-test-reporter/ghutil/fakeghutil"
)

var tsMapForTest = map[string]TestStat{
	"passed": TestStat{
		TestName: "a",
		Passed:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Failed:   []int{},
		Skipped:  []int{},
	},
	"flaky": TestStat{
		TestName: "a",
		Passed:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Failed:   []int{0},
		Skipped:  []int{},
	},
	"failed": TestStat{
		TestName: "a",
		Passed:   []int{},
		Failed:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Skipped:  []int{},
	},
	"notenoughdata": TestStat{
		TestName: "a",
		Passed:   []int{0, 1, 2, 3, 4, 5, 6},
		Failed:   []int{},
		Skipped:  []int{7, 8, 9},
	},
}

var (
	fakeOrg    = "fakeorg"
	fakeRepo   = "fakerepo"
	fakeUserID = int64(99)
	fakeUser   = &github.User{
		ID: &fakeUserID,
	}

	dryrun = false
)

func getFakeGithubIssueClient() *GithubIssue {
	fg := fakeghutil.NewFakeGithubClient()
	fg.Repos = []string{fakeRepo}
	fg.User = fakeUser
	return &GithubIssue{
		user:   fakeUser,
		client: fg,
	}
}

func createFlakyIssue(fgc *GithubIssue, title, body string) (*github.Issue, *github.IssueComment) {
	return createNewIssue(fgc, title, body, "Flaky")
}

func createPassedIssue(fgc *GithubIssue, title, body string) (*github.Issue, *github.IssueComment) {
	return createNewIssue(fgc, title, body, "Passed")
}

func createNewIssue(fgc *GithubIssue, title, body, testStat string) (*github.Issue, *github.IssueComment) {
	issue, _ := fgc.client.CreateIssue(fakeOrg, fakeRepo, title, body)
	commentBody := fmt.Sprintf("Latest result for this test: %s", testStat)
	comment, _ := fgc.client.CreateComment(fakeOrg, fakeRepo, *issue.Number, commentBody)
	return issue, comment
}

func createRepoData(passed, flaky, failed, notenoughdata int, startTime int64) *RepoData {
	config := &JobConfig{
		Repo: fakeRepo,
	}
	tss := map[string]*TestStat{}
	for i := 0; i < passed; i++ {
		ts := tsMapForTest["passed"]
		tss[fmt.Sprintf("passedtest_%d", i)] = &ts
	}
	for i := 0; i < flaky; i++ {
		ts := tsMapForTest["flaky"]
		tss[fmt.Sprintf("flakytest_%d", i)] = &ts
	}
	for i := 0; i < failed; i++ {
		ts := tsMapForTest["failed"]
		tss[fmt.Sprintf("failedtest_%d", i)] = &ts
	}
	for i := 0; i < notenoughdata; i++ {
		ts := tsMapForTest["notenoughdata"]
		tss[fmt.Sprintf("notenoughdatatest_%d", i)] = &ts
	}
	return &RepoData{
		Config:             config,
		TestStats:          tss,
		LastBuildStartTime: &startTime,
	}
}

func TestCreateIssue(t *testing.T) {
	datas := []struct {
		passed, flaky, failed, notenoughdata int
		wantIssues                           int
	}{
		{197, 2, 0, 0, 1}, // flaky rate > 1%, create only 1 issue
		{200, 2, 0, 0, 2}, // flaky rate < 1%, create issue for each
	}

	for _, d := range datas {
		fgc := getFakeGithubIssueClient()
		repoData := createRepoData(d.passed, d.flaky, d.failed, d.notenoughdata, int64(0))
		fgc.processGithubIssueForRepo(repoData, make(map[string][]*flakyIssue), fakeRepo, dryrun)
		issues, _ := fgc.client.ListIssuesByRepo(fakeOrg, fakeRepo, []string{})
		if len(issues) != d.wantIssues {
			t.Fatalf("2%% tests failed, got %d issues, want %d issue", len(issues), d.wantIssues)
		}
	}
}

func TestExistingIssue(t *testing.T) {
	fgc := getFakeGithubIssueClient()
	repoData := createRepoData(200, 2, 0, 0, int64(0))
	flakyIssuesMap, _ := fgc.getFlakyIssues()
	fgc.processGithubIssueForRepo(repoData, flakyIssuesMap, fakeRepo, dryrun)
	existIssues, _ := fgc.client.ListIssuesByRepo(fakeOrg, fakeRepo, []string{})
	flakyIssuesMap, _ = fgc.getFlakyIssues()

	*repoData.LastBuildStartTime++
	fgc.processGithubIssueForRepo(repoData, flakyIssuesMap, fakeRepo, dryrun)
	issues, _ := fgc.client.ListIssuesByRepo(fakeOrg, fakeRepo, []string{})
	if len(existIssues) != len(issues) {
		t.Fatalf("issues already exists, got %d new issues, want 0 new issues", len(issues)-len(existIssues))
	}
}

func TestUpdateIssue(t *testing.T) {
	dataForTest := []struct {
		issueState     string
		ts             TestStat
		passedLastTime bool
		appendComment  bool
		wantStatus     string
		wantErr        error
	}{
		{"open", tsMapForTest["flaky"], false, true, "open", nil},
		{"open", tsMapForTest["flaky"], true, true, "open", nil},
		{"open", tsMapForTest["passed"], false, true, "open", nil},
		{"open", tsMapForTest["passed"], true, true, "closed", nil},
		{"open", tsMapForTest["failed"], false, true, "open", nil},
		{"open", tsMapForTest["failed"], true, true, "open", nil},
		{"open", tsMapForTest["notenoughdata"], false, true, "open", nil},
		{"open", tsMapForTest["notenoughdata"], true, true, "open", nil},
		{"closed", tsMapForTest["flaky"], false, true, "open", nil},
		{"closed", tsMapForTest["flaky"], true, true, "open", nil},
		{"closed", tsMapForTest["passed"], false, true, "closed", nil},
		{"closed", tsMapForTest["passed"], true, true, "closed", nil},
		{"closed", tsMapForTest["failed"], false, true, "closed", nil},
		{"closed", tsMapForTest["failed"], true, true, "closed", nil},
		{"closed", tsMapForTest["notenoughdata"], false, true, "closed", nil},
		{"closed", tsMapForTest["notenoughdata"], true, true, "closed", nil},
	}

	title := "c"
	body := "d"

	for _, data := range dataForTest {
		fgc := getFakeGithubIssueClient()
		issue, comment := createFlakyIssue(fgc, title, body)
		if data.passedLastTime {
			issue, comment = createPassedIssue(fgc, title, body)
		}
		commentBody := comment.GetBody()

		fi := flakyIssue{
			issue:   issue,
			comment: comment,
		}

		dryrun := false
		gotErr := fgc.updateIssue(&fi, "new", &data.ts, dryrun)
		if nil == data.wantErr {
			if nil != gotErr {
				t.Fatalf("update %v, got err: '%v', want err: '%v'", data, gotErr, data.wantErr)
			}
		} else {
			if !strings.HasPrefix(gotErr.Error(), data.wantErr.Error()) {
				t.Fatalf("update %v, got err start with: '%s', want err: '%s'", data, gotErr.Error(), data.wantErr.Error())
			}
		}

		gotComment, _ := fgc.client.GetComment(fakeOrg, fakeRepo, *comment.ID)
		if data.appendComment && gotComment.GetBody() == commentBody {
			t.Fatalf("update comment %v, got: '%s', want: 'new' on top of existing comment", data, gotComment.GetBody())
		}
		if !data.appendComment && gotComment.GetBody() != commentBody {
			t.Fatalf("update comment %v, got: '%s', want: '%s'", data, gotComment.GetBody(), commentBody)
		}
	}
}
