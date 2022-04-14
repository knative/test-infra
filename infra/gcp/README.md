# Knative GCP Infrastructure

We use Terraform to manage our Google Cloud infrastructure.

- [knative-dns](./dns)

  This is the knative-dns project. It holds the `knative.dev` and `knative.team` DNS Zones

- [knative-tests](./tests)

  This is the knative-tests project. It holds the Prow Clusters and associated infrastructure.

- [knative-gsuite](./gsuite)

  This is the knative-gsuite project. It holds a service account that can access Google Workspace Admin API.

- [knative-boskos-01…100](./boskos)

  This the knative-boskos-XX project pools.


# Bootstrapping Terraform - One Time Setup

Terraform needs to be bootstrapped manually before it can be used. This process was done during Knative CNCF Infrastructure Migration. It is noted here for completeness and for potential troubleshooting.

This needs to be ran by a person.

```
# Get the ORG_ID
ORG_ID=$(gcloud organizations describe knative.team --format json | jq .name -r | sed 's:.*/::')

#Create the knative-seed project

gcloud projects create knative-seed --organization $ORG_ID --name "Knative Seed"

#Create the terraform service account

gcloud iam service-accounts create terraform —-display-name “Terraform” --project knative-seed

# Allow Kn Infra Admins to impersonate the service account

gcloud iam service-accounts add-iam-policy-binding terraform@knative-seed.iam.gserviceaccount.com --member='group:kn-infra-gcp-org-admins@knative.team' --role='roles/iam.serviceAccountTokenCreator'

```
