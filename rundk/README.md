# rundk

`rundk` is a tool to run a test command from the test image, by using it
developers can reproduce the test flow as run in the CI environment.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) must be installed
- If you run this tool on MAC, `/tmp` directory must be added to Docker's File
  Sharing settings to allow Docker accessing the temporary files created for
  mount, see the [comment](https://github.com/docker/docker.github.io/issues/4709#issuecomment-639596451).

## Installation

`rundk` can be installed and upgraded by running:

```shell
go get knative.dev/test-infra/rundk
```

## Usage

```shell
Usage of rundk:
  --test-image string
      The image we use to run the test flow. (default "gcr.io/knative-tests/test-infra/prow-tests:stable")
  --entrypoint string
      The entrypoint executable that runs the test commands. (default "runner.sh")
  --enable-docker-in-docker
      Enable running docker commands in the test fllow.
      By enabling this the container will share the same docker daemon in the host machine, so be careful when using it.
  --use-local-gcloud-credentials
      Use the same gcloud credentials as local.
      It can be set either by setting env var GOOGLE_CLOUD_APPLICATION_CREDENTIALS or from ~/.config/gcloud
  --use-local-kubeconfig
      Use the same kubeconfig as local.
      It can be set either by setting env var KUBECONFIG or from ~/.kube/config
  --mounts source1:target1,source2:target2,source3:target3
      A list of extra folders or files separated by comma that need to be mounted to run the test flow.
      It must be in the format of source1:target1,source2:target2,source3:target3.
  --mandatory-env-vars string
      A list of env vars separated by comma that must be set on local, which will then be promoted to the container.
  --optional-env-vars string
      A list of env vars separated by comma that optionally need to be set on local, which will then be promoted to the container.
```

### Example

Run E2E tests for a Knative repository:

```shell
rundk --use-local-gcloud-credentials ./test/e2e-tests.sh --gcp-project-id=one-project-for-testing
```

```shell
rundk --use-local-kubeconfig ./test/e2e-tests.sh --run-tests
```

> Note: the `rundk` command must be run under the root or sub directory of your
> local Knative repository.
