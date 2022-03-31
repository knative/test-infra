# Prow setup for Knative projects

All Prow config files for running Prow jobs for Knative projects are under
[prow/](../prow).

## Setting up Prow for a new organization

1. Install [Knative Prow GitHub App](https://github.com/apps/knative-prow) to
   the organization.

1. Create a team called _Knative Prow Robots_, and make it an Admin of the org
   (or repo).

1. Invite at least [knative-prow-robot](https://github.com/knative-prow-robot)
   for your org. Add it to the robots team you created. For automated releases
   and metrics reporting (e.g., code coverage) you'll need to also add
   [knative-prow-releaser-robot](https://github.com/knative-prow-releaser-robot)
   and [knative-metrics-robot](https://github.com/knative-metrics-robot).

## Setting up Prow for a new repo (reviewers assignment and auto merge)

1. Create the appropriate `OWNERS` and optional `OWNERS_ALIASES` files (at least
   for the root dir).

1. Make sure that _Knative Prow Robots_ team is an Admin of the repo.

1. Add the new repo to
   [jobs_config](../prow/jobs_config), the meta
   config file for generating Prow config and Prow job config. Check the
   config files for other repos for blueprints for what to add. Then run
   `./hack/generate-configs.sh` to regenerate
   [prow/jobs/generated](../prow/jobs/generated), otherwise the presubmit test
   in test-infra will fail. Create a PR with the changes. Once it's merged
   the configs will be automatically updated by a postsubmit job.

1. Wait a few minutes, check that Prow is working by entering `/woof` as a
   comment in any PR in the new repo.

### Setting up jobs for a repo

1. Have the test infrastructure in place (usually this means having at least
   `//test/presubmit-tests.sh` working, and optionally `//hack/release.sh`
   working for automated nightly releases).

1. Update [jobs_config](../prow/jobs_config)
   (usually, copy and update the existing configuration from another
   repository). Run `./hack/generate-configs.sh` to regenerate
   [prow/jobs/config.yaml](../prow/jobs/config.yaml),
   otherwise the presubmit test will fail. Create a pull request with the
   changes. Once it's merged the configs will be automatically updated by a
   postsubmit job.

1. Wait a few minutes, enter `/test [prow_job_name]` or `/test all` or `/retest`
   as a comment in any PR in the repo and ensure the test jobs are executed.
