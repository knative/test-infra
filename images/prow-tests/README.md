# Prow Test Job Image

This directory contains the custom Docker image used by our Prow test jobs.
Some forks of this images exist under images/.

The `prow-tests` image is pinned on a specific `kubekins` image; update
`Dockerfile` if you need to use a newer/different image. This will basically
define the versions of `bazel`, `go`, `kubectl` and other build tools.

The Prow jobs are configured to use the `prow-tests` image tagged with `stable`.
This tag must be manually set in GCR using the Cloud Console.

## Building and publishing a new image

See parent README.md
