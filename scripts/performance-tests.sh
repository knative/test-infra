#!/bin/bash

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

# This is a helper script for Knative performance test scripts.
# See README.md for instructions on how to use it.

source $(dirname ${BASH_SOURCE})/library.sh

# Setup env vars.
readonly PROJECT_NAME="knative-performance"
readonly USER_NAME="mako-job@knative-performance.iam.gserviceaccount.com"
readonly KO_DOCKER_REPO="gcr.io/${PROJECT_NAME}/${REPO_NAME}"
readonly PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS="/etc/performance-test/service-account.json"
readonly PERF_TEST_GITHUB_TOKEN="/etc/performance-test/github-token"
readonly PERF_TEST_SLACK_TOKEN="/etc/performance-test/slack-token"
export TEST_ROOT_PATH="${GOPATH}/src/knative.dev/${REPO_NAME}/test/performance"

# Set up the user credentials for cluster operations.
function setup_user() {
  header "Setup User"

  echo "Using gcloud user ${USER_NAME}"
  gcloud config set core/account ${USER_NAME}
  echo "Using secret defined in ${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}"
  gcloud auth activate-service-account ${USER_NAME} --key-file=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  gcloud config set core/project ${PROJECT_NAME}
}

# Creates a new cluster.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
#             $3 - cluster node num
function create_cluster() {
  header "Creating cluster $1 with $3 nodes in $2"
  gcloud beta container clusters create $1 \
    --addons=HorizontalPodAutoscaling,HttpLoadBalancing \
    --cluster-version=latest \
    --enable-autorepair \
    --enable-ip-alias \
    --enable-stackdriver-kubernetes \
    --machine-type=n1-standard-4 \
    --num-nodes=$3 \
    --region=$2 \
    --scopes cloud-platform
}

# Create service account, github & slack token secrets on the cluster.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
function create_secrets() {
  echo "Creating service account on cluster $1 in zone $2"
  gcloud container clusters get-credentials $1 --zone=$2 --project=${PROJECT_NAME} || abort "failed to get cluster creds"
  kubectl create secret generic service-account --from-file=robot.json=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  kubectl create secret generic tokens --from-file=github-token=${PERF_TEST_GITHUB_TOKEN} --from-file=slack-token=${PERF_TEST_SLACK_TOKEN}
}

# Update resources installed on the cluster.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
function update_cluster() {
  gcloud container clusters get-credentials $1 --zone=$2 --project=${PROJECT_NAME} || abort "failed to get cluster creds"
  echo ">> Setting up 'prod' config-mako"
  cat | kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-mako
data:
  # This should only be used by our performance automation.
  environment: prod
EOF

  echo ">> Deleting all benchmark jobs to avoid noise"
  kubectl delete cronjob --all
  kubectl delete job --all
  
  if function_exists update_knative; then
    update_knative || abort "failed to update knative"
  fi
  # get benchmark_name by removing the prefix from cluster name, e.g. get "load-test" from "serving-load-test"
  local benchmark_name=${1#$REPO_NAME"-"}
  if function_exists update_benchmark; then
    update_benchmark ${TEST_ROOT_PATH}/${benchmark_name} || abort "failed to update benchmark"
  fi
}

# Create a new cluster and create required secrets on it.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
#             $3 - cluster node num
function create_new_cluster() {
  # create a new cluster
  create_cluster $1 $2 $3 || abort "failed to create the new cluster $1"
  
  # create the secrets on the new cluster
  create_secrets $1 $2 || abort "failed to create secrets on the new cluster"

  # update resources on the cluster
  update_cluster $1 $2 || abort "failed to update the cluster"
}

# Delete the old clusters related to the current repo, and recreate them with the same configuration.
function recreate_clusters() {
  header "Recreating all clusters"
  local all_clusters=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone,currentNodeCount)")
  echo "${all_clusters}"
  for cluster in ${all_clusters}; do
    local name=$(echo "${cluster}" | cut -f1 -d",")
    # the cluster name is prefixed with repo name, here we should only handle clusters related to the current repo
    [[ ! ${name} =~ ^${REPO_NAME} ]] && continue
    local zone=$(echo "${cluster}" | cut -f2 -d",")
    local node_count=$(echo "${cluster}" | cut -f3 -d",")
    (( node_count=node_count/3 ))

    # delete the old cluster
    gcloud container clusters delete ${name} --zone ${zone} --quiet

    # create a new cluster and install required resources
    create_new_cluster ${name} ${zone} ${node_count}
  done

  header "Done recreating all clusters"
}

# Update the clusters related to the current repo.
function update_clusters() {
  header "Update all clusters"
  local all_clusters=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone)")
  echo "${all_clusters}"
  for cluster in ${all_clusters}; do
    local name=$(echo "${cluster}" | cut -f1 -d",")
    # the cluster name is prefixed with repo name, here we should only handle clusters related to the current repo
    [[ ! ${name} =~ ^${REPO_NAME} ]] && continue
    local zone=$(echo "${cluster}" | cut -f2 -d",")

    # update all resources installed on the cluster
    update_cluster ${name} ${zone}
  done

  header "Done updating all clusters"
}

CREATE_CLUSTER_BENCHMARK=0

# Parse flags and excute the command.
function main() {
  # set up the user credentials for cluster operations
  setup_user || echo "failed to set up user"

  local create_cluster_benchmark=0

  local command=$1
  # Try parsing the first flag as a command.  
  case ${command} in
    --recreate_clusters) recreate_clusters ;;
    --update_clusters) update_clusters ;;
    --create_cluster_benchmark) create_cluster_benchmark=1 ;;
    *) abort "unknown command ${command}, must be --recreate_clusters / --update_clusters / --create_cluster_benchmark"
  esac
  shift

  local num_nodes=1
  local cluster_name="" 
  local cluster_region="us-central1"
  while [[ $# -ne 0 ]]; do
    [[ $# -ge 2 ]] || abort "missing parameter after $1"
    local parameter=$1
    # Try parsing the options.
    case ${parameter} in
      --name) cluster_name=$2 ;;
      --region) cluster_region=$2 ;;
      --num_nodes) num_nodes=$2 ;;
      *) abort "unknown option ${parameter}" ;;
    esac
    shift
    shift
  done

  if (( create_cluster_benchmark )); then
    [[ ! -z "$cluster_name" ]] || abort "cluster name must be set when creating a new cluster"
    [[ ! -z "$cluster_region" ]] || abort "cluster region must be set when creating a new cluster"
    [[ ! -z "$num_nodes" ]] || abort "number of nodes must be set when creating a new cluster"

    create_new_cluster "${REPO_NAME}-${cluster_name}" "${cluster_region}" "${num_nodes}"
  fi
}
