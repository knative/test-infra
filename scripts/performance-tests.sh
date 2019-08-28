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

# Configurable parameters.
readonly CLUSTER_REGION=${CLUSTER_REGION:-us-central1}
readonly CLUSTER_NODES=${CLUSTER_NODES:-1}
readonly BENCHMARK_ROOT_PATH=${BENCHMARK_ROOT_PATH:-test/performance/benchmarks}

# Setup env vars.
readonly PROJECT_NAME="knative-performance"
readonly USER_NAME="mako-job@knative-performance.iam.gserviceaccount.com"
readonly KO_DOCKER_REPO="gcr.io/${PROJECT_NAME}/${REPO_NAME}"
readonly PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS="/etc/performance-test/service-account.json"
readonly PERF_TEST_GITHUB_TOKEN="/etc/performance-test/github-token"
readonly PERF_TEST_SLACK_TOKEN="/etc/performance-test/slack-token"
readonly CLUSTER_CONFIG_FILE="cluster.properties"
readonly CLUSTER_REGION_CONFIG_NAME="cluster_region"
readonly CLUSTER_NODES_CONFIG_NAME="cluster_nodes"

# Set up the user credentials for cluster operations.
function setup_user() {
  echo ">> Setup User"
  echo "Using gcloud user ${USER_NAME}"
  gcloud config set core/account ${USER_NAME}
  echo "Using secret defined in ${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}"
  gcloud auth activate-service-account ${USER_NAME} --key-file=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  echo "Using gcloud project ${PROJECT_NAME}"
  gcloud config set core/project ${PROJECT_NAME}
}

# Creates a new cluster.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
#             $3 - number of nodes in the cluster
function create_cluster() {
  echo ">> Creating cluster $1 with $3 nodes in $2"
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
  echo ">> Creating service account on cluster $1 in zone $2"
  gcloud container clusters get-credentials $1 --zone=$2 --project=${PROJECT_NAME} || abort "failed to get cluster creds"
  kubectl create secret generic service-account --from-file=robot.json=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  kubectl create secret generic tokens --from-file=github=${PERF_TEST_GITHUB_TOKEN} --from-file=slack=${PERF_TEST_SLACK_TOKEN}
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

  echo ">> Deleting all benchmark jobs to avoid noise in the update process"
  kubectl delete cronjob --all
  kubectl delete job --all
  
  if function_exists update_knative; then
    update_knative || abort "failed to update knative"
  fi
  local benchmark_name=$(get_benchmark_name $1)
  if function_exists update_benchmark; then
    update_benchmark ${BENCHMARK_ROOT_PATH}/${benchmark_name} || abort "failed to update benchmark"
  fi
}

# Create a new cluster and do necessary setups.
# Parameters: $1 - cluster name
#             $2 - cluster zone/region
#             $3 - number of nodes in the cluster
function create_new_cluster() {
  # create a new cluster
  create_cluster $1 $2 $3 || abort "failed to create the new cluster $1"
  
  # create the secrets on the new cluster
  create_secrets $1 $2 || abort "failed to create secrets on the new cluster"

  # update resources on the cluster
  update_cluster $1 $2 || abort "failed to update the cluster"
}

# Create a new cluster for the benchmark and do necessary setups.
# Parameters: $1 - benchmark name
function create_new_benchmark_cluster() {
  local benchmark_path="${BENCHMARK_ROOT_PATH}/$1"
  [ ! -d ${benchmark_path} ] && abort "benchmark $1 does not exist"
  #   
  local cluster_name=$(get_cluster_name $1)
  local cluster_region="${CLUSTER_REGION}"
  local node_count="${CLUSTER_NODES}"
  local config_file_path="${benchmark_path}/${CLUSTER_CONFIG_FILE}"
  if [ ! -f ${config_file_path} ]; then
    echo "cluster.config is not found in benchmark $1, using the default config to create the cluster"
  else
    cluster_region=$(get_config_value ${config_file_path} ${CLUSTER_REGION_CONFIG_NAME} ${CLUSTER_REGION})
    node_count=$(get_config_value ${config_file_path} ${CLUSTER_NODES_CONFIG_NAME} ${CLUSTER_NODES})
  fi

  echo ">> Creating new cluster for benchmark $1 in ${REPO_NAME}"
  create_new_cluster ${cluster_name} ${cluster_region} ${node_count}
}

# Get the value for the given key from the config file, return default value if not found.
# Parameters: $1 - config file path
#             $2 - config name
#             $3 - default value
function get_config_value() {
  local value=$(grep $2 $1 | cut -d'=' -f2)
  echo ${value:-$3}
}

# Get benchmark name from the cluster name.
# Parameters: $1 - cluster name
function get_benchmark_name() {
  # get benchmark_name by removing the prefix from cluster name, e.g. get "load-test" from "serving-load-test"
  echo ${1#$REPO_NAME"-"}
}

# Get cluster name from the benchmark name.
# Parameters: $1 - benchmark name
function get_cluster_name() {
  # cluster_name is [repo_name]-[benchmark_name], e.g. serving-load-test
  echo "${REPO_NAME}-$1"
}

# Delete the old clusters related to the current repo, and recreate them with the same configuration.
function recreate_clusters() {
  header "Recreating all clusters for ${REPO_NAME}"
  local all_clusters=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone,currentNodeCount)")
  echo ">> Listing all clusters:"
  echo "${all_clusters}"
  for cluster in ${all_clusters}; do
    local name=$(echo "${cluster}" | cut -f1 -d",")
    # the cluster name is prefixed with repo name, here we should only handle clusters related to the current repo
    [[ ! ${name} =~ ^${REPO_NAME} ]] && continue
    local zone=$(echo "${cluster}" | cut -f2 -d",")
    local node_count=$(echo "${cluster}" | cut -f3 -d",")
    # we create regional clusters, it will create nodes in all its 3 zones. The node_count we get here is 
    # the total node count, so we'll need to divide with 3 to get the actual regional node count
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
  header "Updating all clusters for ${REPO_NAME}"
  local all_clusters=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone)")
  echo ">> Listing all clusters:"
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

# Try to reset clusters for benchmarks in the current repo.
# There can be three cases:
# 1. If a new benchmark is added, create a new cluster for it;
# 2. If a benchmark is deleted, delete its corresponding cluster;
# 3. If a benchmark is renamed, delete the old cluster and create a new one.
# This function will be run as postsubmit jobs.
function reset_benchmark_clusters() {
  local all_clusters=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="csv[no-heading](name,zone)")
  echo ">> Listing all clusters:"
  echo "${all_clusters}"
  header "Trying to delete unused clusters for ${REPO_NAME}"
  for cluster in ${all_clusters}; do
    local name=$(echo "${cluster}" | cut -f1 -d",")
    # the cluster name is prefixed with repo name, here we should only handle clusters related to the current repo
    [[ ! ${name} =~ ^${REPO_NAME} ]] && continue
    local zone=$(echo "${cluster}" | cut -f2 -d",")
    local cluster_being_used=0
    for benchmark_dir in ${BENCHMARK_ROOT_PATH}/*/; do
      local benchmark_name=$(basename ${benchmark_dir})
      [[ $(get_cluster_name ${benchmark_name}) == ${name} ]] && cluster_being_used=1 && break
    done
    if (( ! cluster_being_used )); then
      gcloud container clusters delete ${name} --zone ${zone} --quiet
    fi
  done
  header "Done deleting unused clusters"

  local all_cluster_names=$(gcloud container clusters list --project="${PROJECT_NAME}" --format="value(name)")
  header "Trying to create new clusters for ${REPO_NAME}"
  for benchmark_dir in ${BENCHMARK_ROOT_PATH}/*/; do
    local benchmark_name=$(basename ${benchmark_dir})
    local cluster_exists=0
    for name in ${all_cluster_names}; do
      [[ $(get_cluster_name ${benchmark_name}) == ${name} ]] && cluster_exists=1 && break
    done
    if (( ! cluster_exists )); then
      create_new_benchmark_cluster ${benchmark_name} || abort "failed to create cluster for the new benchmark ${benchmark_name}"
    fi
  done
  header "Done creating new clusters"
}

# Parse flags and excute the command.
function main() {
  if (( ! IS_PROW )); then
    abort "this script should only be run by Prow since it needs secrets created on Prow cluster"
  fi

  # set up the user credentials for cluster operations
  setup_user || echo "failed to set up user"

  # Try parsing the first flag as a command.  
  case $1 in
    --recreate-clusters) recreate_clusters ;;
    --update-clusters) update_clusters ;;
    --reset-benchmark-clusters) reset_benchmark_clusters ;;
    *) abort "unknown command $1, must be --recreate-clusters, --update-clusters or --reset-benchmark-clusters"
  esac
  shift
}
