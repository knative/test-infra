# Resources Clean Up Tool

This tool is designed to clean up stale resources from gcr, for now it only deletes old images created during testing.

## Basic Usage

Directly invoke [cleanup.sh](cleanup.sh) script with certain flags. There is no-op if invoking or sourcing this script without arguments.

By default the current gcloud credentials are used to delete the images. If necessary, use the flag `--service-account _key-file.json_` to specify a service account that will be performing the access to the gcr.

### Clean up old images from multiple gcrs

Projects to be cleaned up are expected to be defined in a `resources.yaml` file. To remove old images from them, call [cleanup.sh](cleanup.sh) with action "delete-old-gcr-images" and following flags:
- "--project-resource-yaml" as path of `resources.yaml` file - Mandatory
- "--re-project-name" for regex matching projects names - Optional, defaults to `knative-boskos-[a-zA-Z0-9]+`
- "--days-to-keep" - Optional, default `365`

Example:

```./cleanup.sh "delete-old-gcr-images" --project-resource-yaml "ci/prow/boskos/resources.yaml" --days-to-keep 90```

### Clean up old images from a specific gcr

Cleaning up from a specific gcr is supported, except for some special ones (_knative-release_ and _knative-nightly_). Call [cleanup.sh](cleanup.sh) with action "delete-old-images-from-gcr" and following flags:
- "--gcr-to-cleanup" as name of gcr, e.g. "gcr.io/foo" - Mandatory
- "--days-to-keep" - Optional, default `365`

Example:

```./cleanup.sh "delete-old-images-from-gcr" --gcr-to-cleanup "gcr.io/foo" --days-to-keep 90```

## Prow Job

There is a weekly prow job that triggers this tool runs at 11:00/12:00PM(Day light saving) PST every Monday. This tool scans all gcr projects defined in [ci/prow/boskos/resources.yaml](/ci/prow/boskos/resources.yaml) and deletes images older than 90 days.
