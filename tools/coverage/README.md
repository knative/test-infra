# Overview

Code coverage tool has two major features.

1. As a pre-submit tool, it runs code coverage on every single commit to Github and reports coverage change back to the PR as a comment by a robot account. It also has the ability to block a PR from merging if coverage falls below threshold
1. As a periodical running job, it reports on TestGrid to show users how coverage changes over time.

## Design

See design.md

## Build & Release

run `make coverage-dev-image` to build and upload dev version.
Dev version can be triggered on PR through command
`/test pull-knative-$REPONAME-go-coverage-dev`.

If dev version is running as expected,
use the command to upload to production version `make coverage-image`



