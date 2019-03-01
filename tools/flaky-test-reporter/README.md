# flaky-test-reporter

flaky-test-reporter is a tool that identifies flaky tests by retrospectively analyzing continuous flows, tracking flaky tests with Github issues, and sends summary of flaky tests to Slack channels.

## Basic Usage

Flags for this tool are:

* `--service-account` specifies the path of file containing service account for GCS access.
* `--github-token` specifies the path of file containing Github token for Github API calls.
* `--slack-token` specifies the path of file containing Slack token for Slack web API calls.
* `--dry-run` enables dry-run mode.

**IMPORTANT: This tool is _NOT_ intended to run locally, as this could interfere with Github issues and potentially flooding Slack channels**

For debugging purpose it's highly recommended to start with `--dry-run` flag, by passing this flag this tool will only collect information from GCS/Github/Slack, and all the manipulations of Github/Slack resources are faked.

## Prow Job

There is a prow job that triggers this tool runs at 4:00/5:00AM(Day light saving) everyday.
