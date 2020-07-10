# Prow setup for Knative projects

All Prow config files for running Prow jobs for Knative projects are under
[config/prod/prow](../config/prod/prow).

## Setting up Prow for a new organization

1. In GitHub, add the following
   [webhooks](https://developer.github.com/webhooks/) to the org (or repo), in
   `application/json` format and for all events. Ask one of the owners of
   _knative/test-infra_ for the webhook secrets.

   1. `http://prow.knative.dev/hook` (for Prow)
   1. `https://github-dot-knative-tests.appspot.com/webhook` (for Gubernator PR
      Dashboard)

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
   [config_knative.yaml](../config/prod/prow/config_knative.yaml), the meta
   config file for generating Prow config and Prow job config. Check the
   top-level section `presubmits:` and `periodics:` for blueprints for what to
   add. Then run `./hack/generate-configs.sh` to regenerate
   [config/prod/prow/jobs/config.yaml](../config/prod/prow/jobs/config.yaml) and
   [config/prod/prow/core](../config/prod/prow/core), otherwise the presubmit
   test in test-infra will fail. Create a PR with the changes. Once it's merged
   the configs will be automatically updated by a postsubmit job.

1. Wait a few minutes, check that Prow is working by entering `/woof` as a
   comment in any PR in the new repo.

1. Set **tide** as a required status check for the master branch.

   ![Branch Checks](branch_checks.png)

### Setting up jobs for a repo

1. Have the test infrastructure in place (usually this means having at least
   `//test/presubmit-tests.sh` working, and optionally `//hack/release.sh`
   working for automated nightly releases).

1. Update [config_knative.yaml](../config/prod/prow/config_knative.yaml)
   (usually, copy and update the existing configuration from another
   repository). Run `./hack/generate-configs.sh` to regenerate
   [config/prod/prow/jobs/config.yaml](../config/prod/prow/jobs/config.yaml),
   otherwise the presubmit test will fail. Create a pull request with the
   changes. Once it's merged the configs will be automatically updated by a
   postsubmit job.

1. Wait a few minutes, enter `/test [prow_job_name]` or `/test all` or `/retest`
   as a comment in any PR in the repo and ensure the test jobs are executed.

1. Optionally, set the new test jobs as required status checks for the master
   branch.

   ![Branch Checks](branch_checks.png)
