package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/in-toto/in-toto-golang/in_toto"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	prowapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	"k8s.io/test-infra/prow/pod-utils/clone"
	"sigs.k8s.io/yaml"
)

func TestGenerateAttestation(t *testing.T) {
	config := Config{
		EntryPointOptsVariable: "ENTRYPOINTS_OPTIONS",
		CloneRecords: []clone.Record{
			{
				Refs: prowapi.Refs{
					Org:     "knative",
					Repo:    "serving",
					BaseRef: "main",
				},
				FinalSHA: "c82be271867f137d0923be34acd18b6aca452446",
			},
		},
		Subject: []in_toto.Subject{
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
		},
	}

	os.Setenv("BUILD_ID", "1552221225705541632")
	os.Setenv("ENTRYPOINTS_OPTIONS", `{"timeout":10800000000000,"grace_period":15000000000,"artifact_dir":"/logs/artifacts","args":["runner.sh","./hack/release.sh","--publish","--tag-release"],"container_name":"test","process_log":"/logs/process-log.txt","marker_file":"/logs/marker-file.txt","metadata_file":"/logs/artifacts/metadata.json"}`)

	config = ParseEntryPoint(config)

	// Get ProwJob
	prowJob := &prowapi.ProwJob{}
	testPJfile, err := os.ReadFile("testdata/prowjob.yaml")
	if err = yaml.Unmarshal(testPJfile, prowJob); err != nil {
		t.Errorf("cannot unmarshal prowjob file: %v", err)
	}
	config.ProwJob = prowJob

	generatedAttestation := GenerateAttestation(config)

	// expectedAttestation := &provenance.Statement{}
	expectedAttestationFile, err := os.ReadFile("testdata/attestation.json")
	// if err = json.Unmarshal(expectedAttestationFile, expectedAttestation); err != nil {
	// 	t.Errorf("cannot unmarshal attestation.json: %v", err)
	// }

	// Comparing provenance.Statement is broken for some reason, 
	generatedAttestationJSON, err := generatedAttestation.ToJSON()
	if check, err := AreEqualJSON(generatedAttestationJSON, expectedAttestationFile); check == false {
		t.Errorf("expected '%v', got '%v'", string(generatedAttestationJSON), string(expectedAttestationFile))
	} else if err != nil {
		t.Error(err)
	}
}

func AreEqualJSON(b1, b2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	_ = json.Unmarshal(b1, &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling byte 1 :: %s", err.Error())
	}
	_ = json.Unmarshal(b2, &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling byte 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
