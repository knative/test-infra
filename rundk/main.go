package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"knative.dev/test-infra/pkg/helpers"
	"knative.dev/test-infra/pkg/interactive"
)

func main() {
	image := flag.String("test-image", "gcr.io/knative-tests/test-infra/prow-tests:stable", "The image we use to run the test flow.")
	mounts := flag.String("mounts", "", "A list of extra folders or files separated by comma that need to be mounted to run the test flow.")
	mandatoryEnvVars := flag.String("mandatory-env-vars", "GOOGLE_APPLICATION_CREDENTIALS", "A list of env vars separated by comma that must be set on local.")
	optionalEnvVars := flag.String("optional-env-vars", "", "A list of env vars separated by comma that optionally need to be set on local.")
	flag.Parse()

	cmd, cancel := setup(*image, strings.Split(*mounts, ","),
		strings.Split(*mandatoryEnvVars, ","), strings.Split(*optionalEnvVars, ","))
	defer cancel()

	run(cmd, *image, flag.Args()...)
}

func setup(image string, mounts, mandatoryEnvVars, optionalEnvVars []string) (interactive.Docker, func()) {
	var err error

	// Until everyone is using 20.xx version of docker (see https://github.com/docker/cli/pull/1498) adding the --pull flag,
	//  need to separately pull the image first to be sure we have the latest
	pull := interactive.NewCommand("docker", "pull", image)
	if err = pull.Run(); err != nil {
		log.Fatal("Error pulling the test image: ", err)
	}

	builtUpDefers := make([]func(), 0)

	// Setup command
	cmd := interactive.NewDocker()

	// Setup args to promote env vars from local to the Docker container
	envs := interactive.Env{}
	if err = envs.PromoteFromEnv(mandatoryEnvVars...); err != nil {
		log.Fatal("Missing mandatory argument: ", err)
	}
	envs.PromoteFromEnv(optionalEnvVars...) // Optional, so don't check error

	gcloudKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	// If GOOGLE_APPLICATION_CREDENTIALS is not empty, also mount the key file
	// to the container
	if gcloudKey != "" {
		cmd.AddMount("bind", gcloudKey, gcloudKey)
	}
	// Mount the required files and directories on the host machine to the container
	for _, m := range mounts {
		if strings.TrimSpace(m) != "" {
			cmd.AddMount("bind", m, m)
		}
	}

	repoRoot, err := helpers.GetRootDir()
	if err != nil {
		log.Fatal("Error getting the repo's root directory: ", err)
	}

	// Setup temporary directory
	tmpDir, err := ioutil.TempDir("/tmp", "prow-docker.")
	if err != nil {
		log.Fatal("Error setting up the temporary directory: ", err)
	}

	// Copy and mount source code dir
	// Add overlay mount over the user's git repo, so the flow doesn't mess it
	// up
	cleanup := cmd.CopyAndAddMount("bind", tmpDir, repoRoot, filepath.Dir(repoRoot))
	builtUpDefers = append(builtUpDefers, cleanup)

	// Copy and mount for kube context to be available (if reusing an existing cluster)
	// If future use needs other directories, mounting the whole home directory could be a pain
	//  because our prow-tests image will be installing Go in /root/.gvm
	cleanup = cmd.CopyAndAddMount("bind", tmpDir, path.Join(os.Getenv("HOME"), ".kube"), "/root/.kube")
	builtUpDefers = append(builtUpDefers, cleanup)

	cmd.LogFile = path.Join(tmpDir, "build-log.txt")
	log.Print("Logging to ", cmd.LogFile)

	extArtifacts := os.Getenv("ARTIFACTS")
	// Artifacts directory
	if extArtifacts == "" {
		log.Print("Setting local ARTIFACTS directory to ", tmpDir)
		extArtifacts = tmpDir
	}
	cmd.AddMount("bind", extArtifacts, extArtifacts)
	envs["ARTIFACTS"] = extArtifacts
	builtUpDefers = append(builtUpDefers, func() {
		log.Print("💡 Artifacts can be found at ", extArtifacts)
	})
	cmd.AddEnv(envs)

	// Starting directory
	cmd.AddArgs("-w=" + repoRoot)

	return cmd, func() {
		for _, ff := range builtUpDefers {
			if ff != nil {
				ff()
			}
		}
	}
}

func run(cmd interactive.Docker, image string, commandAndArgsOpt ...string) error {
	// Finally add the image then command to run (if any)
	cmd.AddArgs(image)
	cmd.AddArgs("runner.sh")
	cmd.AddArgs(commandAndArgsOpt...)
	log.Println(cmd)
	log.Println("👆 Starting in 3 seconds, ^C to abort!")
	time.Sleep(time.Second * 3)
	return cmd.Run()
}
