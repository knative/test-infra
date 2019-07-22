# Overview

The code coverage tool has two major features:

1. As a pre-submit tool, it runs code coverage on every single commit to Github and reports coverage changes back to the PR as a comment by a robot account. If made a required job on the repository, it can be used to block a PR from merging if coverage falls below a certain threshold.
1. As a periodic running job, it outputs a Junit XML that can be read from other tools like [testgrid](http://testgrid.knative.dev/serving#coverage) to get overall coverage metrics.

## Design

See the [design document](design.md)

## Build and Release

We use a staging version to verify the a new build works, before it goes into production.

run `make coverage-dev-image` to build and upload staging version.
Staging version can be triggered on PR through command
`/test pull-knative-<repository name>-go-coverage-dev`.

Validate the staging version

- To verify pre-submit workflow, add the comment `/test pull-knative-<repository name>-go-coverage-dev` on a PR and see if it produces the same result as `/test pull-knative-<repository name>-go-coverage`.
- To verity periodic workflow, rerun a `post-knative-serving-go-coverage-dev` job and see if the junit xml produced is correct.

After verification, use this command to upload to production version

`make coverage-image`
