## kntest cluster gke

`kntest cluster gke` command is used for creating, deleting or getting a GKE
cluster

## Prerequisite

- `GOOGLE_APPLICATION_CREDENTIALS` set

## Usage

This tool can be invoked from command line with following parameters:

- `--min-nodes`: minimum number of nodes, default 1
- `--max-nodes`: maximum number of nodes, default 3
- `--node-type`: GCE node type, default "e2-standard-4"
- `--region`: GKE region, default "us-central1"
- `--zone`: GKE zone, default empty
- `--project`: GCP project, default empty
- `--name`: cluster name, default empty
- `--release-channel`: GKE release channel, default empty
- `--resource-type`: Boskos resource type
- `--version`: GKE version
- `--backup-regions`: backup regions to be used if cluster creation in primary
  region failed, comma separated list, default "us-west1,us-east1"
- `--addons`: GKE addons, comma separated list, default empty
- `--skip-creation`: should skip creation or not

## Subcommands

### Create

`kntest cluster gke create` will always create a new cluster.

The flow is:

1. Get GCP project name if not provided as a parameter:

   - [In Prow] Acquire from Boskos
   - [Not in Prow] Read from gcloud config

   Failed obtaining project name will fail the tool

1. Get default cluster name if not provided as a parameter
1. Delete cluster if cluster with same name and location already exists in GKE
1. Create a new cluster with the config being provided
1. Write cluster metadata to `${ARTIFACT}/metadata.json`

### Delete

`kntest cluster gke delete` will delete the existing cluster.

The flow is:

1. Acquiring cluster if kubeconfig already points to it
1. If cluster name is defined then getting cluster by its name
1. If no cluster is found from previous step then it fails
1. Delete:
   - [In Prow] Delete cluster asynchronously and release Boskos project
   - [Not in Prow] Delete cluster synchronously

### Get

`kntest cluster gke get` will validate the current kubeconfig context points to
a usable cluster requested by the user.

1. Acquiring cluster if kubeconfig already points to it
1. If cluster name is defined then getting cluster by its name
1. If no cluster is found from previous steps then it fails
