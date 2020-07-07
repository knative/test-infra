# Knative Test Infrastructure

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/knative/test-infra)
[![Go Report Card](https://goreportcard.com/badge/knative/test-infra)](https://goreportcard.com/report/knative/test-infra)
[![LICENSE](https://img.shields.io/github/license/knative/test-infra.svg)](https://github.com/knative/test-infra/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://knative.slack.com/archives/CCSNR4FCH)

The `test-infra` repository contains a collection of tools for testing Knative,
collecting metrics and displaying test results.

## High level architecture

Knative uses [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)
to schedule testing and update issues.

### Gubernator

Knative uses
[gubernator](https://github.com/kubernetes/test-infra/tree/master/gubernator) to
provide a [PR dashboard](https://gubernator.knative.dev/pr) for contributions in
the Knative github organization, and
[Spyglass](https://github.com/kubernetes/test-infra/tree/master/prow/spyglass)
to visualize test details.

### TestGrid

Knative provides a [health dashboard](https://testgrid.knative.dev/) to show
test, code and release health for each repo. It covers key areas such as
continuous integration, code coverage, nightly release, conformance and etc.

### E2E Testing

Our E2E testing uses
[kubetest](https://github.com/kubernetes/test-infra/blob/master/kubetest) to
build/deploy/test Knative clusters.

## Contributing

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).
