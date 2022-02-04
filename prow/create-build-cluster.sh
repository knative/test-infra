#!/usr/bin/env bash

# Copyright 2020 The Knative Authors
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

# This script is used to create a new build cluster for use with knative-prow. The cluster will have a 
# single pd-ssd nodepool that will have autoupgrade and autorepair enabled.
#
# Usage: populate the parameters by setting them below or specifying environment variables then run
# the script and follow the prompts. You'll be prompted to share some credentials and commands
# with the current oncall.
# Requires gcloud and kubectl.

set -o errexit
set -o nounset
set -o pipefail

# Knative specific variables
export TEAM="knative"
export PROJECT="${PROJECT:-knative-tests}"
export ZONE="us-central1-f"
export CLUSTER="knative-prow-build-cluster"
export MACHINE="n1-standard-16"
export GCSBUCKET="${GCSBUCKET:-knative-prow}"
export NODECOUNT=4
export OUT_FILE="build-cluster-kubeconfig.yaml"

# Specific to Prow instance, don't change these.
export PROW_INSTANCE_NAME="${PROW_INSTANCE_NAME:-knative-prow}"
export GCS_BUCKET="${GCS_BUCKET:-${PROW_INSTANCE_NAME}}"
export ADMIN_IAM_MEMBER="${ADMIN_IAM_MEMBER:-group:mdb.cloud-kubernetes-engprod-oncall@google.com}"

# Specific to the build cluster
export TEAM="${TEAM:-}"
export PROJECT="${PROJECT:-${PROW_INSTANCE_NAME}-build-${TEAM}}"
export ZONE="${ZONE:-us-west1-b}"
export CLUSTER="${CLUSTER:-${PROJECT}}"

# Only needed for creating cluster
export MACHINE="${MACHINE:-n1-standard-8}"
export NODECOUNT="${NODECOUNT:-5}"
export DISKSIZE="${DISKSIZE:-100GB}"

# Only needed for creating project
export FOLDER_ID="${FOLDER_ID:-0123}"
export BILLING_ACCOUNT_ID="${BILLING_ACCOUNT_ID:-0123}"  # Find the billing account ID in the cloud console.

bash <(curl -sSfL https://raw.githubusercontent.com/kubernetes/test-infra/master/prow/create-build-cluster.sh) "$@"
