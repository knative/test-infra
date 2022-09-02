# Golang List Tags

This tool will find the Golang tags used within the project reliably using
Go implementation instead of wonky bash scripts.

## Usage

This tool is often used to search for tags in order to perform a Go build test:

```bash
go test \
  -tags "$(go run knative.dev/test-infra/tools/go-ls-tags@latest)" \
  -vet=off \
  -exec echo \
  ./...
```

## Options

### `--extension`

The extension to search for. Defaults to `.go`.

### `--exclude`

The directories to exclude from the search. Defaults to `vendor`, `third_party`,
`hack`, `.git`.

### `--directory`

The directory to search in. Defaults to `.` (current directory).

### `--ignore-file`

An ignore file used to filter out the tags. It is a newline separated list of
tags to ignore. Defaults to `./.gotagsignore`.
