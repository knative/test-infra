# README

## Knative prow

This directory contains prow configs hosted for Knative. This prow is bumped by knative-autobump-config.yaml which uses [Kubernetes generic-autobumper](https://github.com/kubernetes/test-infra/tree/master/prow/cmd/generic-autobumper).

- `Makefile` Commands to interact with the Prow instance regarding configs and
  updates.
- `cluster/*.yaml` Deployments of the Prow cluster.
- `jobs/generated` Generated configuration of the Prow jobs.
- `jobs_config` Input configuration for `configgen` tool.
- `jobs/run_job.sh` Convenience script to start a Prow job from command-line.
- `jobs/pj-on-kind.sh` Convenience script to start a Prow job on kind from
  command-line.
- `cluster/boskos` Just Boskos resource definition and helper scripts; deployments in
  `cluster/build/*`.

### Prow Clusters

- [Prow control plane cluster(default prow cluster)](https://console.cloud.google.com/kubernetes/clusters/details/us-central1/prow?project=knative-tests)
  - Prow deployments, core configs, and plugins are hosted in this repo
  - Prow job configs are hosted in knative/test-infra repo
- [Build cluster](https://console.cloud.google.com/kubernetes/clusters/details/us-central1/prow-build?project=knative-tests)
  - Build cluster deployments are hosted in this repo
  - Boskos resources are hosted in knative/test-infra repo
- [Trusted cluster](https://console.cloud.google.com/kubernetes/clusters/details/us-central1-a/prow-trusted?project=knative-tests)
  - This is a small cluster with important secrets and runs sensitive prow jobs.

### Manually Deploy

Manual deployments are defined as in [Makefile](./Makefile), specifically:

- `make -C prow/knative deploy`: deploys all yamls under [cluster](./cluster)
- `make -C prow/knative deploy-build`: deploys all yamls under [cluster/build](./cluster/build)

## Prow Secrets

Some of the prow secrets are managed by kubernetes external secrets, which
allows prow cluster creating secrets based on values from google secret manager
(Not necessarily the same GCP project where prow is located). Secrets are
declared in this repositories:

- Secrets for prow
service cluster are stored under [`/prow/cluster`](/prow/cluster)
- Secrets for prow build cluster are stored under [`/prow/cluster/build`](/prow/cluster/build)

Please make sure
granting service account
`kubernetes-external-secrets-sa@knative-tests.iam.gserviceaccount.com`
permission for accessing secret manager in the project(GCP allows setting
permission on individual secret level) see more detailed instruction on how to
do so at [Prow
Secret](https://github.com/kubernetes/test-infra/blob/master/prow/prow_secrets.md).
