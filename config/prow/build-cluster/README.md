## Knative prow

This directory contains prow configs hosted for Knative. This prow is bumped by
[`bump.sh`](../bump.sh) script

## Create Build Cluster

`create-build-cluster.sh` wraps `create-build-cluster-helper.sh`,
with defined env vars overriding GCP project, cluster used by Knative.

Create build cluster:

```
# Creating build cluster, skip `Create a SA and secret for uploading results to GCS` as it's not needed
./create-build-cluster.sh

# Connect to knative-prow cluster
make -C .. get-cluster-credentials

# Add the kubeconfig of newly created cluster to kubeconfig secret in Prow cluster
python3 "${GOPATH}/src/k8s.io/test-infra/gencred/merge_kubeconfig_secret.py" --src-key config --dest-key config build-cluster-kubeconfig.yaml
```

The previous step created a build cluster, and registered its kubeconfig as part
of a secret named `kubeconfig`, and mapped it with a nick name of
`build-knative`, and this is what it will be referred to when being used by main
Prow cluster. To make this work, the deployments depending on these secrets need
to be restarted, so far these are `plank`, `deck`, `sinker`, and `crier`

Copy secrets from old cluster to build cluster:

```
# Copy over secrets needed for jobs
./add-secrets.sh
```

### Prow Clusters

- [Prow control plane cluster(default prow cluster)](https://pantheon.corp.google.com/kubernetes/clusters/details/us-central1-f/prow?project=knative-tests) \[Only Google employees can access]
  - Prow deployments, core configs, and plugins are hosted here
  - Prow job configs are hosted in knative/test-infra/config/prow/jobs
- [Build cluster](https://pantheon.corp.google.com/kubernetes/clusters/details/us-central1-f/knative-prow-build-cluster) \[Only Google employees can access]
  - Build cluster deployments are hosted here
  - Boskos resources are hosted in the boskos folder
- [Trusted cluster](https://pantheon.corp.google.com/kubernetes/clusters/details/us-central1-a/prow-trusted) \[Only Google employees can access]
  - This is a very basic cluster with important secrets

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
service cluster are stored under [`/prow/knative/cluster`](/prow/knative/cluster)
- Secrets for prow build cluster are stored under [`/prow/knative/cluster/build`](/prow/knative/cluster/build)

Please make sure
granting service account
`kubernetes-external-secrets-sa@knative-tests.iam.gserviceaccount.com`
permission for accessing secret manager in the project(GCP allows setting
permission on individual secret level) see more detailed instruction on how to
do so at [Prow
Secret](https://github.com/kubernetes/test-infra/blob/master/prow/prow_secrets.md).
