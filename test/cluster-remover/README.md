## cluster-remover tool

cluster-remover tool is a wrapper of clustermanager lib from knative/pkg, which
can be consumed by command line and remove a GKE cluster and/or release Boskos
project if running in Prow:

- [In Prow] Release Boskos project
- [Not in Prow] No-op

## Prerequisite

- `GOOGLE_APPLICATION_CREDENTIALS` set
- `kubeconfig` is set to point to the cluster to be removed

## Usage

This tool can be invoked from command line directly without any argument