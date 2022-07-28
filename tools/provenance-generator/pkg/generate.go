package pkg

import (
	"os"
	"strings"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"sigs.k8s.io/bom/pkg/provenance"
)

func GenerateAttestation(config Config) {
	// Build the arguments RawMessage:

	cloneRecord := config.CloneRecords[0] // Safe as earlier step deletes all other entries

	// Get the k/k commit we are building
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
	p.Metadata.Completeness.Parameters = true  // The parameters are complete as we know the from prow
	p.Metadata.Completeness.Materials = true   // The materials are complete as we only use the github repo
	p.Metadata.Completeness.Environment = true // We don't use environment values to build images/binaries
	startTime := config.StartTime.UTC()
	endTime := time.Now().UTC()
	p.Metadata.BuildStartedOn = &startTime
	p.Metadata.BuildFinishedOn = &endTime
	p.Invocation.ConfigSource.EntryPoint = repo + "/blob/" + cloneRecord.Refs.BaseRef + strings.Trim(config.EntryPointOpts.Args[1], ".")
	p.BuildType = "https://cloudbuild.googleapis.com/CloudBuildYaml@v1"
	p.Invocation.Parameters = config.EntryPointOpts.Args
	// TODO (upodroid) switch to config.Arguments but it needs to be sorted, config.ArgumentsInsertOrder has the key insert order

	p.AddMaterial(("git+" + repo), slsa.DigestSet{"sha1": commitSHA})

	// Create the new attestation and attach the predicate
	attestation := &provenance.Statement{}
	attestation = provenance.NewSLSAStatement()
	attestation.Predicate = p
	attestation.Subject = config.Subject

	attestation.Write("attestation.json")

}
