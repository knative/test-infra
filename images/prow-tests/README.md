# Prow Test Job Image

This directory contains the custom Docker image used by our Prow test jobs.

The `prow-tests` image is pinned on a specific `kubekins` image; update
`Dockerfile` if you need to use a newer/different image. This will basically
define the versions of `bazel`, `go`, `kubectl` and other build tools.

The Prow jobs are configured to use the `prow-tests` image tagged with `stable`.

## Building and publishing a new image

See parent README.md

## Testing and Deploying Images

A normal release process would involve first editing Dockerfile and/or program 
it installs from infra/toolbox, then running `make iterative-build` to build a
local copy and ensure it still builds. Next, if adding a new tool or doing
something complicated, do some rudimentary exploration in the image by running
`make iterative-shell` to ensure it's working; most usage of prow-tests is via
`runner.sh` so you might check a new version of Go by running
`GO_VERSION=go1.30 runner.sh bash` and running `go version` or
`go get something`.

With an image you feel comfortable with deploying, create a PR to
knative/infra and get approval. Once merged, a
[postsubmit Prow job](https://prow.knative.dev/?job=post-knative-infra-prow-tests-image-push)
will get triggered and build the new image into
[Google Cloud Console Artifact Registry](https://console.cloud.google.com/artifacts/docker/knative-tests/us/images).
After the job succeeds, confirm the image is available with the `beta` and
`latest` labels in the registry. If you feel comfortable with promoting the
image, set the `stable` label to the new image and assign `oldstable` to the
previous one. Finally, monitor #test on Slack to see if any new complaints appear.
