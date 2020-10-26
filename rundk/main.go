package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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

	builtUpDefers := make([]func(), 0)

	envs := interactive.Env{}
	if err = envs.PromoteFromEnv(mandatoryEnvVars...); err != nil {
		log.Fatal("Missing mandatory argument: ", err)
	}
	envs.PromoteFromEnv(optionalEnvVars...) // Optional, so don't check error

	// Setup command
	cmd := interactive.NewDocker()

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
	// Mount source code dir
	// Add overlay mount over the user's git repo, so the flow doesn't mess it up
	cancel := cmd.AddRWOverlay(repoRoot, repoRoot)
	builtUpDefers = append(builtUpDefers, cancel)

	// Add overlay mount for kube context to be available (if reusing an existing cluster)
	// If future use needs other directories, mounting the whole home directory could be a pain
	//  because our prow-tests image will be installing Go in /root/.gvm
	cancel = cmd.AddRWOverlay(path.Join(os.Getenv("HOME"), ".kube"), "/root/.kube")
	builtUpDefers = append(builtUpDefers, cancel)

	// Starting directory
	cmd.AddArgs("-w=" + repoRoot)

	// Setup temporary directory
	tmpDir, err := ioutil.TempDir("", "prow-docker.")
	if err != nil {
		log.Fatal("Error setting up the temporary directory: ", err)
	}
	log.Print("Logging to ", tmpDir)
	cmd.LogFile = path.Join(tmpDir, "build-log.txt")

	extArtifacts := os.Getenv("ARTIFACTS")
	// Artifacts directory
	if extArtifacts == "" {
		log.Print("Setting local ARTIFACTS directory to ", tmpDir)
		extArtifacts = tmpDir
	}
	cmd.AddMount("bind", extArtifacts, extArtifacts)
	envs["ARTIFACTS"] = extArtifacts
	builtUpDefers = append(builtUpDefers, func() {
		log.Print("Artifacts found at ", extArtifacts)
	})
	cmd.AddEnv(envs)

	// Until everyone is using 20.xx version of docker (see https://github.com/docker/cli/pull/1498) adding the --pull flag,
	//  need to separately pull the image first to be sure we have the latest
	pull := interactive.NewCommand("docker", "pull", image)
	pull.Run()

	return cmd, func() {
		for _, ff := range builtUpDefers {
			ff()
		}
	}
}

func run(cmd interactive.Docker, image string, commandAndArgsOpt ...string) error {
	// Finally add the image then command to run (if any)
	cmd.AddArgs(image)
	cmd.AddArgs("runner.sh")
	cmd.AddArgs(commandAndArgsOpt...)
	fmt.Println(cmd)
	fmt.Println("Starting in 3 seconds, ^C to abort!")
	time.Sleep(time.Second * 3)
	return cmd.Run()
}
