package pkg

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	// prowapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	"github.com/docker/distribution/reference"
	"github.com/in-toto/in-toto-golang/in_toto"
	"k8s.io/test-infra/prow/entrypoint"
)

type Config struct {
	CloneLogPath           string
	ImageReferencePath     string
	EntryPointOptsVariable string
	Arguments              map[string]string
	ArgumentsInsertOrder   []string
	EntryPointOpts         *entrypoint.Options
	CloneRecords           []Record
	StartTime              time.Time
	Subject                []in_toto.Subject
}

type Subject struct {
}

type Refs struct {
	// Org is something like kubernetes or k8s.io
	Org string `json:"org"`
	// Repo is something like test-infra
	Repo string `json:"repo"`
	// RepoLink links to the source for Repo.
	RepoLink string `json:"repo_link,omitempty"`

	BaseRef string `json:"base_ref,omitempty"`
	BaseSHA string `json:"base_sha,omitempty"`
	// BaseLink is a link to the commit identified by BaseSHA.
	BaseLink string `json:"base_link,omitempty"`

	// PathAlias is the location under <root-dir>/src
	// where this repository is cloned. If this is not
	// set, <root-dir>/src/github.com/org/repo will be
	// used as the default.
	PathAlias string `json:"path_alias,omitempty"`

	// WorkDir defines if the location of the cloned
	// repository will be used as the default working
	// directory.
	WorkDir bool `json:"workdir,omitempty"`

	// CloneURI is the URI that is used to clone the
	// repository. If unset, will default to
	// `https://github.com/org/repo.git`.
	CloneURI string `json:"clone_uri,omitempty"`
	// SkipSubmodules determines if submodules should be
	// cloned when the job is run. Defaults to false.
	SkipSubmodules bool `json:"skip_submodules,omitempty"`
	// CloneDepth is the depth of the clone that will be used.
	// A depth of zero will do a full clone.
	CloneDepth int `json:"clone_depth,omitempty"`
	// SkipFetchHead tells prow to avoid a git fetch <remote> call.
	// Multiheaded repos may need to not make this call.
	// The git fetch <remote> <BaseRef> call occurs regardless.
	SkipFetchHead bool `json:"skip_fetch_head,omitempty"`
}

type Record struct {
	Refs     Refs      `json:"refs"`
	Commands []Command `json:"commands,omitempty"`
	Failed   bool      `json:"failed,omitempty"`

	// FinalSHA is the SHA from ultimate state of a cloned ref
	// This is used to populate RepoCommit in started.json properly
	FinalSHA string `json:"final_sha,omitempty"`
}

// Can't seem to import "k8s.io/test-infra/prow/pod-utils/clone" library :(
type Command struct {
	Command string `json:"command"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}

func LoadParameters(config Config) Config {
	// Parse Entrypoint of the prowjob to pull args, path, etc
	entryPointOpsRaw := os.Getenv(config.EntryPointOptsVariable)
	entryPointOpts := &entrypoint.Options{}
	err := json.Unmarshal([]byte(entryPointOpsRaw), entryPointOpts)
	if err != nil {
		log.Fatalf("failed to unmarshal ENTRYPOINTS_OPTIONS env variable %v", err)
	}
	config.EntryPointOpts = entryPointOpts

	// Sanity checks
	if config.EntryPointOpts.Args[0] != "runner.sh" {
		log.Fatal("this prowjob is misconfigured, expecting runner.sh to be called first")
	}

	config.Arguments = map[string]string{}
	keys := []string{}

	// Parse entrypoint flags
	for _, elem := range config.EntryPointOpts.Args {
		params := strings.Split(elem, "=") // expected values --something=foo
		if len(params) == 1 {
			config.Arguments[params[0]] = ""
			keys = append(keys, params[0])
		} else if len(params) == 2 {
			config.Arguments[params[0]] = params[1]
			keys = append(keys, params[0])
		}
	}
	config.ArgumentsInsertOrder = keys

	//Appoximate start time of job
	fileInfo, err := os.Stat(config.CloneLogPath)
	if err != nil {
		log.Fatalf("failed to open file %v", err)
	}
	config.StartTime = fileInfo.ModTime()

	byteValue, _ := os.ReadFile(config.CloneLogPath)
	cloneLogs := []Record{}

	err = json.Unmarshal(byteValue, &cloneLogs)
	if err != nil {
		log.Fatalf("failed to unmarshal clone-logs file %v", err)
	}
	config.CloneRecords = cloneLogs

	// Lets strip out the empty refs that is inserted by prow
	for i := len(config.CloneRecords) - 1; i >= 0; i-- {
		cloneRecord := config.CloneRecords[i]
		if cloneRecord.FinalSHA == "" {
			config.CloneRecords = append(config.CloneRecords[:i],
				config.CloneRecords[i+1:]...)
		}
	}

	// Parse Image References
	imageRefFile, err := os.ReadFile(config.ImageReferencePath)
	if err != nil {
		log.Fatalf("failed to open file %v", err)
	}
	lines := strings.Split(string(imageRefFile), "\n")
	subjects := []in_toto.Subject{}
	for _, elem := range lines {
		name, err := reference.ParseNamed(elem)
		if err != nil {
			break
		}
		digest := strings.TrimPrefix(elem, name.Name()+"@")
		digestComponents := strings.Split(digest, ":")
		subject := in_toto.Subject{
			Name: name.Name(),
			Digest: map[string]string{
				digestComponents[0]: digestComponents[1],
			},
		}
		subjects = append(subjects, subject)
	}

	config.Subject = subjects
	return config
}
