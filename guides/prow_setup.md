# Prow setup

All Prow related config files for [prow.knative.dev](https://prow.knative.dev)
are under [config/prod/prow](../config/prod/prow). There is also a staging Prow in
[prow-staging.knative.dev](https://prow-staging.knative.dev) for testing the new
config changes before rolling out to production, and all its config files are
under [config/staging/prow](../config/staging/prow).

## Modify Prow cluster configs

Prow is a collection of microservices that are deployed on a Kubernetes cluster.
All the Kubernetes object `.yaml` files (except secrets) for Knative production
Prow are under [config/prod/prow/cluster](../config/prod/prow/cluster). These
files are auto-synced from
[config/staging/prow/cluster](../config/staging/prow/cluster) by the
[prow-config-updater](../tools/prow-config-updater) tool in the Prow staging
process, so developers should _NOT_ modify these files directly.

To make changes to production Prow cluster configs, please follow the staging
process:

1. Check files under
   [config/staging/prow/cluster](../config/staging/prow/cluster) and make
   changes.

1. Create a pull request with the changes.

1. After the PR is approved and merged, a staging process will be started to
   make sure the changes do not break Prow. If it passes, the bot will create a
   new PR to roll out the changes to
   [config/prod/prow/cluster](../config/prod/prow/cluster). The new PR will be
   automatically merged without needing manual approval.

> Note all the config changes will be automatically updated on the Prow cluster
> by a postsubmit Prow job, once the PR is merged.

## Expanding Boskos pool

We use [GKE](https://cloud.google.com/kubernetes-engine) to run integration
tests for Knative projects. To create a GKE cluster, a GCP project is needed.
[boskos](https://github.com/kubernetes/test-infra/tree/master/boskos) is the
resource manager service we use to manage a pool of GCP projects and handle the
transition between different states for these projects.

1. All projects and permissions can be created by running
   `./config/prod/prow/create_boskos_projects.sh`. For example, to create 10
   extra projects run
   `./config/prod/prow/create_boskos_projects 10 0X0X0X-0X0X0X-0X0X0X`. You will
   need to substitute the actual billing ID for the second argument. In the
   event the script fails, it should be easy to follow along with in the GUI or
   run on the CLI. Projects are created with a numeric, incremental prefix
   automatically, based on the contents of
   [prow/boskos_resources.yaml](../config/prod/prow/boskos/boskos_resources.yaml),
   which is automatically updated.

1. Increase the compute CPU quota for the project to 200. Go to
   <https://console.cloud.google.com/iam-admin/quotas?service=compute.googleapis.com&metric=CPUs&project=PROJECT_NAME>
   (replace `PROJECT_NAME` with the real project name) and click `Edit Quota`.
   Select at least five regions to increase the quota
   (`us-central1, us-west1, us-east1, europe-west1, asia-east1`). This needs to
   be done manually and should get automatically approved once the request is
   submitted. For the reason, enter _Need more resources for running tests_.

1. Create a pull request with the changes. Once it's merged the configs will be
   automatically updated by a postsubmit job.
