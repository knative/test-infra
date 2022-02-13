## kntest cluster gke

`kntest cluster gke` command is used for creating, deleting or getting a GKE
cluster

## Usage

This tool can be invoked from command line. The following parameters are common
for all subcommands:

- `--gcp-credential-file`: the GCP credential file that will be used in the
  cluster operations \
  Will fall back to `GOOGLE_APPLICATION_CREDENTIALS` if it's not set.
- `--project`: GCP project, default empty
- `--name`: cluster name, default empty
- `--region`: GKE region, default "us-central1". \
  Can be more than one to set as backup regions.
- `--resource-type`: Boskos resource type, default "gke-project"
- `--save-meta-data`: whether or not save the meta data for the current cluster
  into `metadata.json`, default to be false.

## Subcommands

### Create

`kntest cluster gke create` will always create a new cluster. It accepts the
following extra parameters:

- `--min-nodes`: minimum number of nodes, default 1
- `--max-nodes`: maximum number of nodes, default 3
- `--node-type`: GCE node type, default "e2-standard-4"
- `--release-channel`: GKE release channel, default empty
- `--version`: GKE version, default "latest"
- `--addons`: GKE addons, comma separated list, default empty

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
