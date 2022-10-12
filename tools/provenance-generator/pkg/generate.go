package pkg

import (
	"os"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	prowapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	"k8s.io/test-infra/prow/pod-utils/clone"
	"sigs.k8s.io/bom/pkg/provenance"
)

type BuildConfig struct {
	command []string
	prowJob *prowapi.ProwJob
}

func GenerateAttestation(config Config) *provenance.Statement {
	var cloneRecord clone.Record
	if len(config.CloneRecords) == 1 {
		cloneRecord = config.CloneRecords[0]
	}

	// Get the commit we are building
	commitSHA := cloneRecord.FinalSHA
	// Create the predicate to populate it with the current
	// run metadata:
	p := provenance.NewSLSAPredicate()

	// SLSA v02, builder ID is a TypeURI
	p.Builder.ID = "https://prow.knative.dev"

	// Some of these fields have yet to be checked to assign the
	// correct values to them
	// This is commented as the in-toto go port does not have it
	repo := "https://github.com/" + cloneRecord.Refs.Org + "/" + cloneRecord.Refs.Repo
	p.Metadata.BuildInvocationID = os.Getenv("BUILD_ID")
	p.Metadata.Completeness.Parameters = true  // The parameters are complete as we know them from prow
	p.Metadata.Completeness.Materials = true   // The materials are complete as we only use the github repo
	p.Metadata.Completeness.Environment = true // We don't use environment values to build images/binaries
	startTime := config.ProwJob.CreationTimestamp.Time.UTC()
	endTime := time.Now().UTC() // The prowjob is still running when this timestamp is needed. Attestation happens after the images/binaries are built
	p.Metadata.BuildStartedOn = &startTime
	p.Metadata.BuildFinishedOn = &endTime
	p.Invocation.ConfigSource.EntryPoint = "https://github.com/knative/test-infra/tree/main/prow/jobs/generated/" + cloneRecord.Refs.Org
	p.BuildType = "https://prow.knative.dev/ProwJob@v1"
	var buildConfig interface{}
	// set config.ProwJob.Status to empty as it is not complete when the generator runs
	config.ProwJob.Status = prowapi.ProwJobStatus{}
	buildConfig = map[string]interface{}{
		"command":    config.EntryPointOpts.Args,
		"entrypoint": config.EntryPointOpts,
		"prowjob":    config.ProwJob,
	}
	p.BuildConfig = buildConfig
	// TODO (upodroid) switch to config.Arguments but it needs to be sorted, config.ArgumentsInsertOrder has the key insert order

	p.AddMaterial(("git+" + repo), slsa.DigestSet{"sha1": commitSHA})

	// Create the new attestation and attach the predicate
	attestation := &provenance.Statement{}
	attestation = provenance.NewSLSAStatement()
	attestation.Predicate = p
	attestation.Subject = config.Subject
	return attestation
}
