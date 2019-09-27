#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

readonly PROJECT=${1:?"First argument must be the boskos project name."}

readonly PROJECT_OWNERS=("prime-engprod-sea@google.com")
readonly PROJECT_GROUPS=("knative-productivity-admins@googlegroups.com")
readonly PROJECT_SAS=(
    "knative-tests@appspot.gserviceaccount.com"
    "prow-job@knative-tests.iam.gserviceaccount.com"
    "prow-job@knative-nightly.iam.gserviceaccount.com"
    "prow-job@knative-releases.iam.gserviceaccount.com")
readonly PROJECT_APIS=(
    "cloudresourcemanager.googleapis.com"
    "compute.googleapis.com"
    "container.googleapis.com")

# Add an owner to the PROJECT
for owner in ${PROJECT_OWNERS[@]}; do
  echo "NOTE: Adding owner ${owner}"
  gcloud projects add-iam-policy-binding ${PROJECT} --member group:${owner} --role roles/owner
done

# Add all GROUPS as editors
for group in ${PROJECT_GROUPS[@]}; do
  echo "NOTE: Adding group ${group}"
  gcloud projects add-iam-policy-binding ${PROJECT} --member group:${group} --role roles/editor
done

for sa in ${PROJECT_SAS[@]}; do
  # Add all service accounts with roles "editor" and "storage admin"
  echo "NOTE: Adding service account ${sa}"
  gcloud projects add-iam-policy-binding ${PROJECT} --member serviceAccount:${sa} --role roles/editor
  gcloud projects add-iam-policy-binding ${PROJECT} --member serviceAccount:${sa} --role roles/storage.admin
  # As required by step 6 in https://github.com/google/knative-gcp/tree/master/docs/storage,
  # grant the GCS service account the permissions to publish to GCP Pub/Sub.
  echo "Activating GCS service account"
  curl -s -X GET -H "Authorization: Bearer $(gcloud auth application-default print-access-token)" "https://www.googleapis.com/storage/v1/projects/${PROJECT}/serviceAccount"
  PROJECT_NUMBER="$(gcloud projects describe ${PROJECT} | grep "^projectNumber" | cut -d':' -f2 | xargs)"
  GCS_SERVICE_ACCOUNT="service-${PROJECT_NUMBER}@gs-project-accounts.iam.gserviceaccount.com"
  echo "GCS service account is ${GCS_SERVICE_ACCOUNT}"
  gcloud projects add-iam-policy-binding ${PROJECT} \
    --member=serviceAccount:${GCS_SERVICE_ACCOUNT} \
    --role roles/pubsub.publisher
done

# Enable APIS
for api in ${PROJECT_APIS[@]}; do
  echo "NOTE: Enabling API ${api}"
  gcloud services enable ${api} --project=${PROJECT}
done
