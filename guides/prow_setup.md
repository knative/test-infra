# Prow setup

Prow is a collection of microservices that are deployed on a Kubernetes cluster.
All the Kubernetes object `.yaml` files (except secrets) for Knative prow are
under
[oss-test-infra repo](https://github.com/GoogleCloudPlatform/oss-test-infra/tree/master/prow/knative),
and managed by [oss oncall](go.k8s.io/oncall).

## Expanding Boskos pool

We use [GKE](https://cloud.google.com/kubernetes-engine) to run integration
tests for Knative projects. To create a GKE cluster, a GCP project is needed.
[boskos](https://github.com/kubernetes/test-infra/tree/master/boskos) is the
resource manager service we use to manage a pool of GCP projects and handle the
transition between different states for these projects.

1. All projects and permissions can be created by running
   `./config/prod/build-cluster/create_boskos_projects.sh`. For example, to create 10
   extra projects run
   `./config/prod/build-cluster/create_boskos_projects 10 0X0X0X-0X0X0X-0X0X0X`. You will
   need to substitute the actual billing ID for the second argument. In the
   event the script fails, it should be easy to follow along with in the GUI or
   run on the CLI. Projects are created with a numeric, incremental prefix
   automatically, based on the contents of
   [build-cluster/boskos_resources.yaml](../config/prod/build-cluster/boskos/boskos_resources.yaml),
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
