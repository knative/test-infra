package main

import (
	"context"
	"fmt"
	"github.com/knative/test-infra/coverage/artifacts/artsTest"
	"github.com/knative/test-infra/coverage/gcs"
	"github.com/knative/test-infra/coverage/gcs/gcsFakes"
	"github.com/knative/test-infra/coverage/githubUtil/githubFakes"
	"github.com/knative/test-infra/coverage/githubUtil/githubPr"
	"github.com/knative/test-infra/coverage/test"
	"log"
	"testing"
)

const (
	testPresubmitBuild = 787
)

func repoDataForTest() *githubPr.GithubPr {
	ctx := context.Background()
	log.Printf("creating fake repo data \n")

	return &githubPr.GithubPr{
		RepoOwner:     "fakeRepoOwner",
		RepoName:      "fakeRepoName",
		Pr:            7,
		RobotUserName: "fakeCovbot",
		GithubClient:  githubFakes.FakeGithubClient(),
		Ctx:           ctx,
	}
}

func gcsArtifactsForTest() *gcs.GcsArtifacts {
	return &gcs.GcsArtifacts{
		Ctx:       context.Background(),
		Bucket:    "fakeBucket",
		Client:    gcsFakes.NewFakeStorageClient(),
		Artifacts: artsTest.LocalArtsForTest("gcsArts-").Artifacts,
	}
}

func preSubmitForTest() (data *gcs.PreSubmit) {
	repoData := repoDataForTest()
	build := gcs.GcsBuild{
		StorageClient: gcsFakes.NewFakeStorageClient(),
		Bucket:        gcsFakes.FakeGcsBucketName,
		Job:           gcsFakes.FakePreSubmitProwJobName,
		Build:         testPresubmitBuild,
	}
	pbuild := gcs.PresubmitBuild{
		GcsBuild:      build,
		Artifacts:     *gcsArtifactsForTest(),
		PostSubmitJob: gcsFakes.FakePostSubmitProwJobName,
	}
	data = &gcs.PreSubmit{
		GithubPr:       *repoData,
		PresubmitBuild: pbuild,
	}
	log.Println("finished preSubmitForTest()")
	return
}

func TestRunPresubmit(t *testing.T) {
	log.Println("Starting TestRunPresubmit")
	arts := artsTest.LocalArtsForTest("TestRunPresubmit")
	arts.ProduceProfileFile("./" + test.CovTargetRelPath)
	p := preSubmitForTest()
	RunPresubmit(p, arts)
	if !test.FileOrDirExists(arts.LineCovFilePath()) {
		t.Fatalf("No line cov file found in %s\n", arts.LineCovFilePath())
	}
}

// tests the construction of gcs url from PreSubmit
func TestK8sGcsAddress(t *testing.T) {
	data := preSubmitForTest()
	data.Build = 1286
	actual := data.UrlGcsLineCovLinkWithMarker(3)

	expected := fmt.Sprintf("https://storage.cloud.google.com/%s/pr-logs/pull/"+
		"%s_%s/%s/%s/%s/artifacts/line-cov.html#file3",
		gcsFakes.FakeGcsBucketName, data.RepoOwner, data.RepoName, data.PrStr(), gcsFakes.FakePreSubmitProwJobName, data.BuildStr())
	if actual != expected {
		t.Fatal(test.StrFailure("", expected, actual))
	}
	fmt.Printf("line cov link=%s", actual)
}
