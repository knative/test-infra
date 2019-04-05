# Prow setup

## Creating the cluster

1. Create the GKE cluster, the role bindings and the GitHub secrets. You might
   need to update [Makefile](./prow/Makefile). For details, see <https://github.com/kubernetes/test-infra/blob/master/prow/getting_started.md>.

1. Ensure the GCP projects listed in [resources.yaml](./prow/boskos/resources.yaml)
   are created.

1. Apply [config_start.yaml](./prow/config_start.yaml) to the cluster.

1. Apply Boskos [config_start.yaml](./prow/boskos/config_start.yaml) to the cluster.

1. Run `make update-cluster`, `make update-boskos`, `make update-config`,
   `make update-plugins` and `make update-boskos-config`.

1. If SSL needs to be reconfigured, promote your ingress IP to static in Cloud
   Console, and [create the TLS secret](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls).

## Expanding Boskos pool

1. Create a new GCP project and add it to [resources.yaml](./prow/boskos/resources.yaml).

1. Make the following accounts editors of the project:
   * `knative-productivity-admins@googlegroups.com`
   * `knative-tests@appspot.gserviceaccount.com`
   * `prow-job@knative-tests.iam.gserviceaccount.com`
   * `prow-job@knative-nightly.iam.gserviceaccount.com`
   * `prow-job@knative-releases.iam.gserviceaccount.com`

1. Ensure that there is at least one other owner of the project. A good choice
   is one of the members of the `knative-productivity-admins@googlegroups.com`
   group.

1. Enable the Compute Engine API for the project (e.g., by visiting
   <https://console.developers.google.com/apis/api/compute.googleapis.com/overview?project=XXXXXXXX>).

1. Enable the Kubernetes Engine API for the project (e.g., by visiting
   <https://console.cloud.google.com/apis/api/container.googleapis.com/overview?project=XXXXXXXX>).

1. Run `make update-boskos-config`.

## Setting up Prow for a new repo (reviewers assignment and auto merge)

1. Create the appropriate `OWNERS` files (at least one for the root dir).

1. Make sure that *Knative Robots* is an Admin of the repo.

1. Add the repo to the [tide section](https://github.com/knative/test-infra/blob/6c1fc9978de156385ddbe431c3a5920d321d4382/ci/prow/make_config.go#L222)
   in the Prow config file generator and run `make config`. Create a PR with the
   changes to the generator and to the [config.yaml](./prow/config.yaml) file. Once
   the PR is merged, ask one of the owners of knative/test-infra to deploy the new
   config.

1. Wait a few minutes, check that Prow is working by entering `/woof` as a
   comment in any PR in the new repo.

1. Set **tide** as a required status check for the master branch.

   ![Branch Checks](branch_checks.png)

### Setting up jobs for a new repo

1. Have the test infrastructure in place (usually this means having at least
   `//test/presubmit-tests.sh` working, and optionally `//hack/release.sh` working
   for automated nightly releases).

1. Merge a pull request that:

   1. Updates [config_knative.yaml](./prow/config_knative.yaml), the Prow config
      file (usually, copy and update the existing configuration from another repository).
      Run `make config` to regenerate [config.yaml](./prow/config.yaml), otherwise
      the presubmit test will fail.

   1. Updates the Gubernator config with the new log dirs.

   1. Updates the Testgrid config with the new buckets, tabs and dashboard.

1. Ask one of the owners of *knative/test-infra* to:

   1. Run `make update-config` in `ci/prow`.

   1. Run `make deploy` in `ci/gubernator`.

   1. Run `make update-config` in `ci/testgrid`.

1. Wait a few minutes, enter `/retest` as a comment in any PR in the repo and
   ensure the test jobs are executed.

1. Set the new test jobs as required status checks for the master branch.

   ![Branch Checks](branch_checks.png)
