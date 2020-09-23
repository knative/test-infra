# Prow Test Job Image

This directory contains the custom Docker image used by our Prow test jobs.

The `prow-tests` image is pinned on a specific `kubekins` image; update
`Dockerfile` if you need to use a newer/different image. This will basically
define the versions of `bazel`, `go`, `kubectl` and other build tools.

The Prow jobs are configured to use the `prow-tests` image tagged with `stable`.

## Building and publishing a new image

See parent README.md

## Testing and Deploying Images

When `make cloud_build` is run for `prow-tests`, it publishes your image with
the `:beta` label. At night Pacific time, a set of duplicate jobs are run using
this new image and viewable at https://testgrid.knative.dev/beta-prow-tests

A normal release process would involve first editing Dockerfile and/or programs
it installs from test-infra, then running `make iterative-build` to build a
local copy and ensure it still builds. Next, if adding a new tool or doing
something complicated, do some rudimentary exploration in the image by running
`make iterative-shell` to ensure it's working; most usage of prow-tests is via
`runner.sh` so you might check a new version of Go by running
`GO_VERSION=go1.30 runner.sh bash` and running `go version` or
`go get something`.

With an image you feel comfortable with deploying, create a PR to
knative/test-infra and get approval. Once merged, pull upstream into your master
branch and run `make cloud_build`. This will upload your image to the registry
at http://gcr.io/knative-tests/test-infra/prow-tests with a date-commit_hash,
`latest`, and `beta` tags. Let the image run at least a single overnight and
ensure the jobs in the https://testgrid.knative.dev/beta-prow-tests testgrid are
as good as the jobs at https://testgrid.knative.dev/knative for master branch
and releases under "knative-X.Y" at https://testgrid.knative.dev ; it's known
that some are just bad so you will either learn to know them or just do a lot of
clicking; it's unlikely to take more than ten minutes to review all the jobs
unless something is actually wrong.

Once you are happy the :beta image is working, go to
https://gcr.io/knative-tests/test-infra/prow-tests and add the label `stable` to
your image with the `beta` tag. Monitor #test on Slack to see if any new
complaints appear.
