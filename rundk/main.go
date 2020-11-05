package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"knative.dev/test-infra/pkg/helpers"
	"knative.dev/test-infra/pkg/interactive"
)

var (
	image                     string
	entrypoint                string
	dindEnabled               bool
	useLocalGcloudCredentials bool
	useLocalKubeconfig        bool
	mounts                    []string
	mandatoryEnvVars          []string
	optionalEnvVars           []string
)

func main() {
	pflag.StringVar(&image, "test-image", "gcr.io/knative-tests/test-infra/prow-tests:stable", "The image we use to run the test flow.")
	pflag.StringVar(&entrypoint, "entrypoint", "runner.sh", "The entrypoint executable that runs the test commands.")
	pflag.BoolVar(&dindEnabled, "enable-docker-in-docker", false, "Enable running docker commands in the test flow. "+
		"By enabling this the container will share the same docker daemon in the host machine, so be careful when using it.")
	pflag.BoolVar(&useLocalGcloudCredentials, "use-local-gcloud-credentials", false, "Use the same gcloud credentials as local, which can be set "+
		"either by setting env var GOOGLE_CLOUD_APPLICATION_CREDENTIALS or from ~/.config/gcloud")
	pflag.BoolVar(&useLocalKubeconfig, "use-local-kubeconfig", false, "Use the same kubeconfig as local, which can be set "+
		"either by setting env var KUBECONFIG or from ~/.kube/config")
	pflag.StringSliceVar(&mounts, "mounts", []string{}, "A list of extra folders or files separated by comma that need to be mounted to run the test flow."+
		"It must be in the format of `source1:target1,source2:target2,source3:target3`.")
	pflag.StringSliceVar(&mandatoryEnvVars, "mandatory-env-vars", []string{}, "A list of env vars separated by comma that must be set on local.")
	pflag.StringSliceVar(&optionalEnvVars, "optional-env-vars", []string{}, "A list of env vars separated by comma that optionally need to be set on local.")
	pflag.Parse()

	cmd, cancel := setup()
	defer cancel()

	run(cmd, pflag.Args()...)
}

func setup() (interactive.Docker, func()) {
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
	envs := interactive.Env{}

	// Setup temporary directory to save the artifacts
	tmpDir, err := ioutil.TempDir("/tmp", "prow-docker.")
	if err != nil {
		log.Fatal("Error setting up the temporary directory: ", err)
	}

	// Setup docker-in-docker
	if dindEnabled {
		cmd.AddMount("bind", "/var/run/docker.sock", "/var/run/docker.sock")
		// On a MAC machine, set to host networking to make DinD work. Not
		// required on Linux.
		if runtime.GOOS == "darwin" {
			cmd.AddArgs("--network", "host")
		}
	}

	// Setup gcloud credentials
	if useLocalGcloudCredentials {
		gcloudKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		// If GOOGLE_APPLICATION_CREDENTIALS is not empty, also mount the key file
		// to the container
		if gcloudKey != "" {
			envs["GOOGLE_APPLICATION_CREDENTIALS"] = gcloudKey
			if _, err := os.Stat(gcloudKey); !os.IsNotExist(err) {
				cleanup := cmd.CopyAndAddMount("bind", tmpDir, gcloudKey, gcloudKey)
				builtUpDefers = append(builtUpDefers, cleanup)
			}
		}
		// Copy and mount for gcloud default credentials to be available
		defaultGcloudConfigPath := path.Join(os.Getenv("HOME"), ".config/gcloud")
		if _, err := os.Stat(defaultGcloudConfigPath); !os.IsNotExist(err) {
			cleanup := cmd.CopyAndAddMount("bind", tmpDir, defaultGcloudConfigPath, "/root/.config/gcloud")
			builtUpDefers = append(builtUpDefers, cleanup)
		}
	}

	// Setup kubeconfig
	if useLocalKubeconfig {
		kubeconfig := os.Getenv("KUBECONFIG")
		// If KUBECONFIG is not empty, also mount the kubeconfig files to the
		// container
		if kubeconfig != "" {
			envs["KUBECONFIG"] = kubeconfig
			for _, f := range strings.Split(kubeconfig, string(os.PathListSeparator)) {
				if _, err := os.Stat(f); !os.IsNotExist(err) {
					cleanup := cmd.CopyAndAddMount("bind", tmpDir, f, f)
					builtUpDefers = append(builtUpDefers, cleanup)
				}
			}
		}
		// Copy and mount for default kube context to be available
		defaultKubeconfigPath := path.Join(os.Getenv("HOME"), ".kube")
		if _, err := os.Stat(defaultKubeconfigPath); !os.IsNotExist(err) {
			cleanup := cmd.CopyAndAddMount("bind", tmpDir, path.Join(os.Getenv("HOME"), ".kube"), "/root/.kube")
			builtUpDefers = append(builtUpDefers, cleanup)
		}
	}

	// Setup args to promote env vars from local to the Docker container
	if err = envs.PromoteFromEnv(mandatoryEnvVars...); err != nil {
		log.Fatal("Missing mandatory argument: ", err)
	}
	envs.PromoteFromEnv(optionalEnvVars...) // Optional, so don't check error

	// Setup args to mount the required files and directories on the host
	// machine to the container
	for _, m := range mounts {
		sourceAndTarget := strings.Split(m, ":")
		if len(sourceAndTarget) != 2 {
			log.Fatalf("The mount string %q must be in the format of [source:target]", m)
		}
		cmd.AddMount("bind", sourceAndTarget[0], sourceAndTarget[1])
	}

	// Copy and mount source code root dir
	repoRoot, err := helpers.GetRootDir()
	if err != nil {
		log.Fatal("Error getting the repo's root directory: ", err)
	}
	cleanup := cmd.CopyAndAddMount("bind", tmpDir, repoRoot, repoRoot)
	builtUpDefers = append(builtUpDefers, cleanup)

	// Setup logging
	cmd.LogFile = path.Join(tmpDir, "build-log.txt")
	log.Print("Logging to ", cmd.LogFile)

	// Setup and mount ARTIFACTS directory
	extArtifacts := os.Getenv("ARTIFACTS")
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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting the working directory: ", err)
	}
	// Set starting directory to be the same as the current working directory.
	cmd.AddArgs("-w=" + wd)

	return cmd, func() {
		for _, ff := range builtUpDefers {
			if ff != nil {
				ff()
			}
		}
	}
}

func run(cmd interactive.Docker, commandAndArgsOpt ...string) error {
	// Finally add the image then command to run (if any)
	cmd.AddArgs(image)
	log.Println(cmd)
	if len(commandAndArgsOpt) != 0 {
		cmd.AddArgs(entrypoint)
		cmd.AddArgs(commandAndArgsOpt...)
		log.Println("👆 Starting in 3 seconds, ^C to abort!")
	}

	time.Sleep(time.Second * 3)
	return cmd.Run()
}
