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

# Setup env vars
export PROJECT_NAME="knative-performance"
export USER_NAME="mako-job@knative-performance.iam.gserviceaccount.com"
export TEST_ROOT_PATH="$GOPATH/src/knative.dev/${REPO_NAME}/test/performance"
export KO_DOCKER_REPO="gcr.io/knative-performance"

# Creates a new cluster.
# $1 -> name, $2 -> zone/region, $3 -> num_nodes
function create_cluster() {
  header "Creating cluster $1 with $3 nodes in $2"
  gcloud beta container clusters create ${1} \
    --addons=HorizontalPodAutoscaling,HttpLoadBalancing \
    --machine-type=n1-standard-4 \
    --cluster-version=latest --region=${2} \
    --enable-stackdriver-kubernetes --enable-ip-alias \
    --num-nodes=${3} \
    --enable-autorepair \
    --scopes cloud-platform
}

# Create serice account secret on the cluster.
# $1 -> cluster_name, $2 -> cluster_zone
function create_secret() {
  echo "Create service account on cluster $1 in zone $2"
  gcloud container clusters get-credentials $1 --zone=$2 --project=${PROJECT_NAME} || abort "Failed to get cluster creds"
  kubectl create secret generic service-account --from-file=robot.json=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
}

# Set up the user credentials for cluster operations.
function setup_user() {
  header "Setup User"

  gcloud config set core/account ${USER_NAME}
  gcloud auth activate-service-account ${USER_NAME} --key-file=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  gcloud config set core/project ${PROJECT_NAME}

  echo "gcloud user is $(gcloud config get-value core/account)"
  echo "Using secret defined in ${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}"
}

# Update resources installed on the cluster.
# $1 -> cluster_name, $2 -> cluster_zone
function update_cluster() {
  local cluster_name=$1
  local cluster_zone=$2
  gcloud container clusters get-credentials $cluster_name --zone=$cluster_zone --project=${PROJECT_NAME} || abort "Failed to get cluster creds"
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

  if function_exists update_knative; then
    update_knative || fail_test "failed to update knative"
  fi
  local benchmark_name=${cluster_name#$REPO_NAME"-"}
  if function_exists update_benchmark; then
    update_benchmark ${benchmark_name} || fail_test "failed to update benchmark"
  fi
}

# Create a new cluster and install serving components and apply benchmark yamls.
# $1 -> cluster_name, $2 -> cluster_zone, $3 -> node_count
function create_new_cluster() {
  # create a new cluster
  create_cluster $1 $2 $3 || abort "Failed to create the new cluster $1"
  
  # create the secret on the new cluster
  create_secret $1 $2 || abort "Failed to create secrets on the new cluster"

  # update components on the cluster, e.g. serving and istio
  update_cluster $1 $2 || abort "Failed to update the cluster"
}

function recreate_clusters() {
  # set up the user credentials for cluster operations
  setup_user

  header "Recreating all clusters"
  for cluster in $(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone,currentNodeCount)"); do
    [[ ! ${PROJECT_NAME} =~ ^${REPO_NAME} ]] && continue
    name=$(echo $cluster | cut -f1 -d",")
    zone=$(echo $cluster | cut -f2 -d",")
    node_count=$(echo $cluster | cut -f3 -d",")
    (( node_count=node_count/3 ))

    # delete the old cluster
    gcloud container clusters delete ${name} --zone ${zone} --quiet

    # create a new cluster and update all the components
    create_new_cluster ${name} ${zone} ${node_count}
  done

  header "Done recreating all clusters"
}

function update_clusters() {
  # set up the credential for cluster operations
  setup_user

  # Get all clusters to update and ko apply config. Use newline to split
  header "Update all clusters"
  IFS=$'\n'
  for cluster in $(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone)"); do
    [[ ! ${PROJECT_NAME} =~ ^${REPO_NAME} ]] && continue
    name=$(echo $cluster | cut -f1 -d",")
    zone=$(echo $cluster | cut -f2 -d",")

    update_cluster ${name} ${zone}
  done

  header "Done updating all clusters"
}

NUM_NODES=1
CLUSTER_NAME=""
CLUSTER_REGION="us-central1"
RECREATE_CLUSTERS=0
UPDATE_CLUSTERS=0
CREATE_CLUSTER_BENCHMARK=0

# Parse flags and excute the command.
function main() {
  if (( ! IS_PROW )); then
    abort "this script should only be run by Prow since it needs secrets created on Prow cluster"
  fi

  local command=$1
  # Try parsing the first flag as a command.  
  case ${command} in
    --recreate_clusters) RECREATE_CLUSTERS=1 ;;
    --update_clusters) UPDATE_CLUSTERS=1 ;;
    --create_cluster_benchmark) CREATE_CLUSTER_BENCHMARK=1 ;;
    *) abort "unknown command ${command}, must be --recreate_clusters / --update_clusters / --create_cluster_benchmark"
  esac
  shift

  while [[ $# -ne 0 ]]; do
    [[ $# -ge 2 ]] || abort "missing parameter after $1"
    local parameter=$1
    # Try parsing the options.
    case ${parameter} in
      --num_nodes) NUM_NODES=$2 ;;
      --name) CLUSTER_NAME=$2 ;;
      --region) CLUSTER_REGION=$2 ;;
      *) abort "unknown option ${parameter}" ;;
    esac
    shift
    shift
  done

  readonly NUM_NODES
  readonly CLUSTER_NAME
  readonly CLUSTER_REGION
  readonly RECREATE_CLUSTERS
  readonly UPDATE_CLUSTERS
  readonly CREATE_CLUSTER_BENCHMARK

  if (( RECREATE_CLUSTERS )); then
    create_test_cluster
  elif (( UPDATE_CLUSTERS )); then
    setup_test_cluster
  else
    create_cluster_benchmark
  fi
}

main $@
