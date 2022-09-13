# Modscope

A tool to display Go module information like current module name, path, and
with support to list the Go workspace modules.

## Usage

### List workspace modules

```shell
$ modscope ls
knative.dev/hack
knative.dev/hack/schema
knative.dev/hack/test/e2e
```

### List workspace module paths

```shell
$ modscope ls -p
/home/example/git/knative-hack
/home/example/git/knative-hack/schema
/home/example/git/knative-hack/test/e2e
```

### Current module

```shell
$ pwd
/home/example/git/knative-hack/schema/vendor
$ modscope curr
knative.dev/hack/schema
```

## Installation

```shell
$ go install knative.dev/test-infra/tools/modscope@latest
```

or use the `go run` command:

```shell
$ go run knative.dev/test-infra/tools/modscope@latest --help
```
