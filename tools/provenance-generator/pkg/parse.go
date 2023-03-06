package pkg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/in-toto/in-toto-golang/in_toto"
	prowapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	"k8s.io/test-infra/prow/entrypoint"
	"k8s.io/test-infra/prow/pod-utils/clone"
	"sigs.k8s.io/yaml"
)

type Config struct {
	CloneLogPath           string
	ImageReferencePath     string
	EntryPointOptsVariable string
	FileCheckSumPath       string
	Arguments              map[string]string
	ArgumentsInsertOrder   []string
	EntryPointOpts         *entrypoint.Options
	CloneRecords           []clone.Record
	StartTime              time.Time
	Subject                []in_toto.Subject
	ProwJob                *prowapi.ProwJob
	ProwUrl                string
}

func LoadParameters(config Config) Config {
	// Parse Entrypoint of the prowjob to pull args, path, etc
	config.Arguments = map[string]string{}
	keys := []string{}

	config = ParseEntryPoint(config)
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
		log.Fatalf("failed to open file: %v", err)
	}
	config.StartTime = fileInfo.ModTime()

	byteValue, _ := os.ReadFile(config.CloneLogPath)
	cloneLogs := []clone.Record{}

	log.Printf("log path %v", config.CloneLogPath)
	if err := json.Unmarshal(byteValue, &cloneLogs); err != nil {
		log.Fatalf("failed to unmarshal clone-logs file: %v", err)
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
	config = GenerateSubject(config)
	// Get ProwJob info and add it to the statement
	url := config.ProwUrl + "/prowjob?prowjob=" + os.Getenv("PROW_JOB_ID")
	config.ProwJob, _ = FetchProwJob(url)
	return config
}

func FetchProwJob(url string) (*prowapi.ProwJob, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("status code not 2XX: %v", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prowJob *prowapi.ProwJob
	// For some reason, server returns yaml instead of json, yikes https://github.com/kubernetes/test-infra/blob/master/prow/cmd/deck/main.go#L1430
	if err := yaml.Unmarshal(data, &prowJob); err != nil {
		log.Printf("cannot unmarshal data from deck: %v", err)
		return nil, err
	}
	return prowJob, nil
}

func GenerateSubject(config Config) Config {
	subjects := []in_toto.Subject{}
	if config.ImageReferencePath != "" {
		// Parse Image References
		imageRefFile, err := os.ReadFile(config.ImageReferencePath)
		if err != nil {
			log.Fatalf("failed to open file %v", err)
		}
		lines := strings.Split(string(imageRefFile), "\n")

		for _, elem := range lines {
			name, err := name.ParseReference(elem)
			if err != nil {
				break // Intentional to avoid processing empty lines, etc
			}
			algo, digest, _ := strings.Cut(name.Identifier(), ":")
			subject := in_toto.Subject{
				Name: name.Context().Name(),
				Digest: map[string]string{
					algo: digest,
				},
			}
			subjects = append(subjects, subject)
		}
	} else if config.FileCheckSumPath != "" {
		// Parse Image References
		checkSumFile, err := os.ReadFile(config.FileCheckSumPath)
		if err != nil {
			log.Fatalf("failed to open file %v", err)
		}
		lines := strings.Split(string(checkSumFile), "\n")
		// https://github.com/hashicorp/go-getter/blob/main/checksum.go
		for _, elem := range lines {
			parts := strings.Fields(elem)
			switch len(parts) {
			case 2:
				// GNU-style:
				//  <checksum>  file1
				//  <checksum> *file2
				algo, err := parseAlgorithm(parts[0])
				if err != nil {
					break // checksum wasn't identified properly so don't add it to the subject
				}
				subject := in_toto.Subject{
					Name: parts[1],
					Digest: map[string]string{
						algo: parts[0],
					},
				}
				subjects = append(subjects, subject)
			case 0:
				break // We didn't encounter a checksum
			}
		}

	}
	config.Subject = subjects

	return config
}

func parseAlgorithm(checksumValue string) (string, error) {
	var algo string
	bytes, err := hex.DecodeString(checksumValue)
	if err != nil {
		return "", fmt.Errorf("invalid checksum: %s", err)
	}
	switch len(bytes) {
	case md5.Size:
		algo = "md5"
	case sha1.Size:
		algo = "sha1"
	case sha256.Size:
		algo = "sha256"
	case sha512.Size:
		algo = "sha512"
	default:
		return "", fmt.Errorf("unknown type for checksum: %s", checksumValue)
	}
	return algo, nil
}

func ParseEntryPoint(config Config) Config {
	entryPointOpsRaw := os.Getenv(config.EntryPointOptsVariable)
	entryPointOpts := &entrypoint.Options{}
	if err := json.Unmarshal([]byte(entryPointOpsRaw), entryPointOpts); err != nil {
		log.Fatalf("failed to unmarshal ENTRYPOINTS_OPTIONS env variable %v", err)
	}
	config.EntryPointOpts = entryPointOpts

	if config.EntryPointOpts.Args[0] != "runner.sh" {
		log.Fatal("this prowjob is misconfigured, expecting runner.sh to be called first")
	}

	return config
}
