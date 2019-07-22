# Overview

Code coverage tool has two major features.

1. As a pre-submit tool, it runs code coverage on every single commit to Github and reports coverage changes back to the PR as a comment by a robot account. If made an required job on repo, it can be used to block a PR from merging if coverage falls below threshold.
1. As a periodical running job, it outputs a Junit XML that can be read from other tools like [testgrid](http://testgrid.knative.dev/serving#coverage) to get overall coverage metrics.

## Design

See [design.md](design.md)

## Build and Release

Run `make coverage-dev-image` to build and upload a staging version, intended for testing and debugging.
Staging version can be triggered on PR through command
`/test pull-knative-<repository name>-go-coverage-dev`.

Validate the staging version

- To verify pre-submit workflow, add the comment `/test pull-knative-<repository name>-go-coverage-dev` on a PR and see if it produces the same result as `/test pull-knative-<repository name>-go-coverage`.
- To verity periodic workflow, rerun a `post-knative-serving-go-coverage-dev` job and see if the junit xml produced is correct.

After verification, use this command to upload to production version

`make coverage-image`
