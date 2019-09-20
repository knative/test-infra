## cluster-creator tool

cluster-creator tool is a wrapper of clustermanager lib from knative/pkg, which
can be consumed by command line and acquire a GKE cluster in following order:

1. Acquiring cluster if kubeconfig already points to it
1. Get GCP project name if not provided as a parameter:
    - [In Prow] Acquire from Boskos
    - [Not in Prow] Read from gcloud config

    Failed obtaining project name will fail the tool
1. Get default cluster name if not provided as a parameter
1. Delete cluster if cluster with same name and location already exists in GKE
1. Create cluster

## Prerequisite

- `GOOGLE_APPLICATION_CREDENTIALS` set

## Usage

This tool can be invoked from command line with following parameters:

- `--min-nodes`: minumum number of nodes, default 1
- `--max-nodes`: maximum number of nodes, default 3
- `--node-type`: GCE node type, default "n1-standard-4"
- `--region`: GKE region, default "us-central1"
- `--zone`: GKE zone, default empty
- `--project`" GCP project, default empty
- `--name`: cluster name, default empty
- `--backup-regions`: backup regions to be used if cluster creation in primary
  region failed, comma separated list, default "us-west1,us-east1"
- `--addons`: GKE addons, comma separated list, default empty
