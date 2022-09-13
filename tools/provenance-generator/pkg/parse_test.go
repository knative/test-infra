package pkg

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/in-toto/in-toto-golang/in_toto"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	prowapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	"sigs.k8s.io/yaml"
)

func TestParseSubjectImages(t *testing.T) {
	config := Config{
		ImageReferencePath: "testdata/image-refs.txt",
	}
	expectedSubject := []in_toto.Subject{
		{
			Name: "gcr.io/knative-releases/knative.dev/serving/cmd/controller",
			Digest: slsa.DigestSet{
				"sha256": "bac158dfb0c73d13ed42266ba287f1a86192c0ba581e23fbe012d30a1c34837c",
			},
		},
		{
			Name: "gcr.io/knative-releases/knative.dev/serving/cmd/queue",
			Digest: slsa.DigestSet{
				"sha256": "83f6888ea9561495f67334d044ffa8ad067d251ad953358dda7ea5183390cc69",
			},
		},
	}
	config = GenerateSubject(config)

	if !reflect.DeepEqual(config.Subject, expectedSubject) {
		t.Errorf("expected '%s', got '%s'", expectedSubject, config.Subject)
	}
}

func TestParseSubjectFiles(t *testing.T) {
	config := Config{
		FileCheckSumPath: "testdata/shasums.txt",
	}
	expectedSubject := []in_toto.Subject{
		{
			Name: "func_darwin_amd64",
			Digest: slsa.DigestSet{
				"sha256": "f13541b1dc1ff6c1f61c974b6905fbde36846e85f20095ba7d7e4d39ffcb3f05",
			},
		},
		{
			Name: "func_darwin_arm64",
			Digest: slsa.DigestSet{
				"sha256": "7fbb8869db3f30eb16e76998e798a6027f410b09e9721136f991afd164bffa68",
			},
		},
		{
			Name: "func_linux_amd64",
			Digest: slsa.DigestSet{
				"sha256": "b072fa3d18aab0ef281e260bd76263605159ee60589449aec72f8433d55517cf",
			},
		},
		{
			Name: "func_windows_amd64.exe",
			Digest: slsa.DigestSet{
				"sha256": "943750331d8b090b2f0b90212a15a7bca0a82c791b697803edbec0572af112e3",
			},
		},
	}
	config = GenerateSubject(config)
	if !reflect.DeepEqual(config.Subject, expectedSubject) {
		t.Errorf("expected '%s', got '%s'", expectedSubject, config.Subject)
	}
}

func TestGetProwJob(t *testing.T) {
	prowJob := &prowapi.ProwJob{}
	testPJfile, _ := os.ReadFile("testdata/prowjob.yaml")
	if err := yaml.Unmarshal(testPJfile, prowJob); err != nil {
		t.Errorf("cannot unmarshal prowjob file: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testPJfile)
	}))
	defer server.Close()
	fetchedProwJob, err := FetchProwJob(server.URL)
	if err != nil {
		t.Errorf("returned err: %v", err)
	}
	if !reflect.DeepEqual(fetchedProwJob, prowJob) {
		t.Errorf("expected '%v', got '%v'", prowJob, fetchedProwJob)
	}
}
