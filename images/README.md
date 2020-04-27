# Prow Job Images

This directory contains custom Docker images used by our Prow jobs.

## Building and publishing a new image

To build and push a new image, you can run `make push` from an image directory. This should only be
done once a PR is approved and from the master branch.

Users typically keep their working directory at the repo root and run relative make commands, like
`make -C images/prow-tests push`, but for brevity's sake this document just mentions the short way.

For testing purposes, you can build with --no-cache, which does not rerun a step in the Dockerfile
if the text in the Dockerfile does not change (which, note, if you are doing `go install` or
similar, and you update files, the image will not change! Either `make build` to rebuild everything
or add/remove a no-op flag to the command you care about). Do this by running `make
iterative-build`. Run `make iterative-shell` to get a shell in this image.

Note that you must have proper permission in the `knative-tests` project to push new images to the
GCR.
