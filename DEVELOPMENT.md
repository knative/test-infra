# Development

This doc explains how to setup a development environment so you can get started
[contributing](https://www.knative.dev/contributing/) to `Knative Test Infra`.
Also take a look at:

- [The pull request workflow](https://www.knative.dev/contributing/contributing/#pull-requests)
- [Iterating](#iterating)

## Prerequisites

Before submitting a PR, see also [CONTRIBUTING.md](./CONTRIBUTING.md).

### Sign up for GitHub

Start by creating [a GitHub account](https://github.com/join), then setup
[GitHub access via SSH](https://help.github.com/articles/connecting-to-github-with-ssh/).

### Checkout your fork

To check out this repository:

1. Create your own
   [fork of this repo](https://help.github.com/articles/fork-a-repo/)
1. Clone it to your machine:

```shell
mkdir -p ${GOPATH}/src/knative.dev
cd ${GOPATH}/src/knative.dev
git clone git@github.com:${YOUR_GITHUB_USERNAME}/test-infra.git
cd test-infra
git remote add upstream https://github.com/knative/test-infra.git
git remote set-url --push upstream no_push
```

_Adding the `upstream` remote sets you up nicely for regularly
[syncing your fork](https://help.github.com/articles/syncing-a-fork/)._

## Iterating

As you make changes to the code-base, there are two special cases to be aware
of:

- **If you change an input to generated code**, then you must run
  [`./hack/update-codegen.sh`](./hack/update-codegen.sh). Inputs include:

  - Prow configs templates
    [config/prod/prow/config_knative.yaml](./config/prod/prow/config_knative.yaml).

  - Prow configs generator [tools/config-generator](./tools/config-generator).

- **If you change a package's deps** (including adding an external dependency),
  then you must run [`./hack/update-deps.sh`](./hack/update-deps.sh).

These are both idempotent, and we expect that running these at `HEAD` to have no
diffs. Code generation and dependencies are automatically checked to produce no
diffs for each pull request.

update-deps.sh runs go get/mod command. In some cases, if newer dependencies are
required, you need to run "go get" manually.

### Updating existing dependencies

To update existing dependencies execute

```shell
./hack/update-deps.sh --upgrade && ./hack/update-codegen.sh
```
