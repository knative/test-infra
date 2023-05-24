# ⚠ ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠ ⚠
# ⚠ This repo is no longer in use please go to [infra repo](https://github.com/knative/infra) ⚠
# ⚠ ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠ ⚠

## Knative Test Infrastructure

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/knative/test-infra)
[![Go Report Card](https://goreportcard.com/badge/knative/test-infra)](https://goreportcard.com/report/knative/test-infra)
[![LICENSE](https://img.shields.io/github/license/knative/test-infra.svg)](https://github.com/knative/test-infra/blob/main/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://knative.slack.com/archives/CCSNR4FCH)

The `test-infra` repository contains a collection of tools for testing Knative,
collecting metrics and displaying test results.
This is the main repo for the [Productivity Working Group](https://github.com/knative/community/blob/main/working-groups/WORKING-GROUPS.md#productivity)

The Productivity Working Group also has other repos:
- [knative/actions](https://github.com/knative/actions)

  Reusable github workflows and actions

- [knative/hack](https://github.com/knative/hack)

  Shellscripts used across the repos placed in a separate repo to avoid
  dependency cycles

- [knative/release](https://github.com/knative/release)

  Release documentation and tools

- [knative-sandbox/.github](https://github.com/knative-sandbox/.github)

  Tools for github actions

- [knative-sandbox/kperf](https://github.com/knative-sandbox/kperf)

  A performance test framework

- [knative-sandbox/knobots](https://github.com/knative-sandbox/knobots)

  Automated pull requests to fix up the code (based on github actions)

## Tools we use

We use two big platforms for running automation:
- [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)

  To schedule testing and update issues. Prow handles the merge queue
  and makes sure every commit passes tests. Prow builds releases from release branches.

- [Github Actions](https://docs.github.com/en/actions)

  We use github actions for some automated tests, coordinating releases
  and syncronizing files between repos

<!-- TODO: As an improvement for the architecture section maybe mention how
the tools fit together -->

### Spyglass

Knative uses
[Spyglass](https://github.com/kubernetes/test-infra/tree/master/prow/spyglass)
to visualize test details.

### TestGrid

Knative provides a [health dashboard](https://testgrid.knative.dev/) to show
test, code and release health for each repo. It covers key areas such as
continuous integration, nightly release, conformance and etc.

### E2E Testing

Our E2E testing uses
[kubetest2](https://github.com/kubernetes-sigs/kubetest2) to
build/deploy/test Knative clusters (managed by Prow).

## Contributing

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).

## Guides

To setup the CI/CD flow for a knative project, see [guides](./guides/README.md).
