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

readonly FIRST=${1:?"First argument is the first number of the new project(s)."}
readonly NUMBER=${2:?"Second argument is the number of new projects."}
readonly BILLING_ACCOUNT=${3:?"Third argument must be the billing account."}
readonly OUTPUT_FILE=${4:?"Fourth argument should be a file name all project names will be appended to in a resources.yaml format."}

readonly CUSTOM_ROLE_NAME="KnativeIntegrationTestsRunner"
readonly CUSTOM_ROLE_FILE="custom_role.yaml"

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

for (( i=0; i<${NUMBER}; i++ )); do
  PROJECT="knative-boskos-$(( i + ${FIRST} ))"
  # This Folder ID is google.com/google-default
  # If this needs to be changed for any reason, GCP project settings must be updated.
  # Details are available in Google's internal issue 137963841.
  gcloud projects create ${PROJECT} --folder=396521612403
  gcloud beta billing projects link ${PROJECT} --billing-account=${BILLING_ACCOUNT}

  # Set permissions for users on this new project
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

  # Create the custom role in this new project
  gcloud iam roles create ${CUSTOM_ROLE_NAME} -q --project ${PROJECT} --file ${CUSTOM_ROLE_FILE}
  for sa in ${PROJECT_SAS[@]}; do
    echo "NOTE: Adding service account ${sa}"
    # Bind the custom role to the SA
    gcloud projects add-iam-policy-binding ${PROJECT} --member serviceAccount:${sa} --role projects/${PROJECT}/roles/${CUSTOM_ROLE_NAME}
  done

  # Enable APIS
  for api in ${PROJECT_APIS[@]}; do
    echo "NOTE: Enabling API ${api}"
    gcloud services enable ${api} --project=${PROJECT}
  done
  last_project_line=$(grep -n "knative-boskos-" resources.yaml | tail -n 1 | cut -d: -f1)
  ((last_project_line++))
  sed -e "${last_project_line}i\ \ - ${PROJECT}" -i ${OUTPUT_FILE}
done
