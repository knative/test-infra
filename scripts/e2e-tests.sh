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

# This is a helper script for Knative E2E test scripts.
# See README.md for instructions on how to use it.

source $(dirname "${BASH_SOURCE[0]}")/library.sh

# Build a resource name based on $E2E_BASE_NAME, a suffix and $BUILD_NUMBER.
# Restricts the name length to 40 chars (the limit for resource names in GCP).
# Name will have the form $E2E_BASE_NAME-<PREFIX>$BUILD_NUMBER.
# Parameters: $1 - name suffix
function build_resource_name() {
  local prefix=${E2E_BASE_NAME}-$1
  local suffix=${BUILD_NUMBER}
  # Restrict suffix length to 20 chars
  if [[ -n "${suffix}" ]]; then
    suffix=${suffix:${#suffix}<20?0:-20}
  fi
  local name="${prefix:0:20}${suffix}"
  # Ensure name doesn't end with "-"
  echo "${name%-}"
}

# Test cluster parameters

# Configurable parameters
# export E2E_CLUSTER_REGION and E2E_CLUSTER_ZONE as they're used in the cluster setup subprocess
export E2E_CLUSTER_REGION=${E2E_CLUSTER_REGION:-us-central1}
# By default we use regional clusters.
export E2E_CLUSTER_ZONE=${E2E_CLUSTER_ZONE:-}

# Default backup regions in case of stockouts; by default we don't fall back to a different zone in the same region
readonly E2E_CLUSTER_BACKUP_REGIONS=${E2E_CLUSTER_BACKUP_REGIONS:-us-west1 us-east1}
readonly E2E_CLUSTER_BACKUP_ZONES=${E2E_CLUSTER_BACKUP_ZONES:-}

readonly E2E_GKE_ENVIRONMENT=${E2E_GKE_ENVIRONMENT:-prod}
readonly E2E_GKE_COMMAND_GROUP=${E2E_GKE_COMMAND_GROUP:-beta}

# Each knative repository may have a different cluster size requirement here,
# so we allow calling code to set these parameters.  If they are not set we
# use some sane defaults.
readonly E2E_MIN_CLUSTER_NODES=${E2E_MIN_CLUSTER_NODES:-1}
readonly E2E_MAX_CLUSTER_NODES=${E2E_MAX_CLUSTER_NODES:-3}
readonly E2E_CLUSTER_MACHINE=${E2E_CLUSTER_MACHINE:-e2-standard-4}

readonly E2E_BASE_NAME="k${REPO_NAME}"
readonly E2E_CLUSTER_NAME=$(build_resource_name e2e-cls)
readonly E2E_NETWORK_NAME=$(build_resource_name e2e-net)
readonly TEST_RESULT_FILE=/tmp/${E2E_BASE_NAME}-e2e-result

# Flag whether test is using a boskos GCP project
IS_BOSKOS=0

# Tear down the test resources.
function teardown_test_resources() {
  # On boskos, save time and don't teardown as the cluster will be destroyed anyway.
  (( IS_BOSKOS )) && return
  header "Tearing down test environment"
  function_exists test_teardown && test_teardown
  (( ! SKIP_KNATIVE_SETUP )) && function_exists knative_teardown && knative_teardown
}

# Run the given E2E tests. Assume tests are tagged e2e, unless `-tags=XXX` is passed.
# Parameters: $1..$n - any go test flags, then directories containing the tests to run.
function go_test_e2e() {
  local go_test_args=()
  # Remove empty args as `go test` will consider it as running tests for the current directory, which is not expected.
  for arg in "$@"; do
    [[ -n "$arg" ]] && go_test_args+=("$arg")
  done
  [[ ! " $*" == *" -tags="* ]] && go_test_args+=("-tags=e2e")
  report_go_test -race -count=1 "${go_test_args[@]}"
}

# Setup the test cluster for running the tests.
function setup_test_cluster() {
  # Fail fast during setup.
  set -o errexit
  set -o pipefail

  header "Test cluster setup"
  kubectl get nodes

  header "Setting up test cluster"

  # Set the actual project the test cluster resides in
  # It will be a project assigned by Boskos if test is running on Prow,
  # otherwise will be ${GCP_PROJECT} set up by user.
  E2E_PROJECT_ID="$(gcloud config get-value project)"
  export E2E_PROJECT_ID
  readonly E2E_PROJECT_ID

  # Save some metadata about cluster creation for using in prow and testgrid
  save_metadata

  local k8s_user
  k8s_user=$(gcloud config get-value core/account)
  local k8s_cluster
  k8s_cluster=$(kubectl config current-context)

  is_protected_cluster "${k8s_cluster}" && \
    abort "kubeconfig context set to ${k8s_cluster}, which is forbidden"

  # If cluster admin role isn't set, this is a brand new cluster
  # Setup the admin role and also KO_DOCKER_REPO if it is a GKE cluster
  if [[ -z "$(kubectl get clusterrolebinding cluster-admin-binding 2> /dev/null)" && "${k8s_cluster}" =~ ^gke_.* ]]; then
    acquire_cluster_admin_role "${k8s_user}" "${E2E_CLUSTER_NAME}" "${E2E_CLUSTER_REGION}" "${E2E_CLUSTER_ZONE}"
    # Incorporate an element of randomness to ensure that each run properly publishes images.
    export KO_DOCKER_REPO=gcr.io/${E2E_PROJECT_ID}/${E2E_BASE_NAME}-e2e-img/${RANDOM}
  fi

  # Safety checks
  is_protected_gcr "${KO_DOCKER_REPO}" && \
    abort "\$KO_DOCKER_REPO set to ${KO_DOCKER_REPO}, which is forbidden"

  # Use default namespace for all subsequent kubectl commands in this context
  kubectl config set-context "${k8s_cluster}" --namespace=default

  echo "- gcloud project is ${E2E_PROJECT_ID}"
  echo "- gcloud user is ${k8s_user}"
  echo "- Cluster is ${k8s_cluster}"
  echo "- Docker is ${KO_DOCKER_REPO}"

  export KO_DATA_PATH="${REPO_ROOT_DIR}/.git"

  # Do not run teardowns if we explicitly want to skip them.
  (( ! SKIP_TEARDOWNS )) && trap teardown_test_resources EXIT

  # Handle failures ourselves, so we can dump useful info.
  set +o errexit
  set +o pipefail

  if (( ! SKIP_KNATIVE_SETUP )) && function_exists knative_setup; then
    # Wait for Istio installation to complete, if necessary, before calling knative_setup.
    (( ! SKIP_ISTIO_ADDON )) && (wait_until_batch_job_complete istio-system || return 1)
    knative_setup || fail_test "Knative setup failed"
  fi
  if function_exists test_setup; then
    test_setup || fail_test "test setup failed"
  fi
}

# Signal (as return code and in the logs) that all E2E tests passed.
function success() {
  echo "**************************************"
  echo "***        E2E TESTS PASSED        ***"
  echo "**************************************"
  dump_metrics
  exit 0
}

# Exit test, dumping current state info.
# Parameters: $1 - error message (optional).
function fail_test() {
  [[ -n $1 ]] && echo "ERROR: $1"
  dump_cluster_state
  dump_metrics
  exit 1
}

RUN_TESTS=0
SKIP_KNATIVE_SETUP=0
SKIP_ISTIO_ADDON=0
SKIP_TEARDOWNS=0
GCP_PROJECT=""
E2E_SCRIPT=""
E2E_CLUSTER_VERSION="latest"
GKE_ADDONS=""
EXTRA_CLUSTER_CREATION_FLAGS=()
E2E_SCRIPT_CUSTOM_FLAGS=()

# Parse flags and initialize the test cluster.
function initialize() {
  E2E_SCRIPT="$(get_canonical_path "$0")"

  cd "${REPO_ROOT_DIR}"
  while [[ $# -ne 0 ]]; do
    local parameter=$1
    # Try parsing flag as a custom one.
    if function_exists parse_flags; then
      parse_flags $@
      local skip=$?
      if [[ ${skip} -ne 0 ]]; then
        # Skip parsed flag (and possibly argument) and continue
        # Also save it to it's passed through to the test script
        for ((i=1;i<=skip;i++)); do
          E2E_SCRIPT_CUSTOM_FLAGS+=("$1")
          shift
        done
        continue
      fi
    fi
    # Try parsing flag as a standard one.
    case ${parameter} in
      --run-tests) RUN_TESTS=1 ;;
      --skip-knative-setup) SKIP_KNATIVE_SETUP=1 ;;
      --skip-teardowns) SKIP_TEARDOWNS=1 ;;
      --skip-istio-addon) SKIP_ISTIO_ADDON=1 ;;
      *)
        [[ $# -ge 2 ]] || abort "missing parameter after $1"
        shift
        case ${parameter} in
          --gcp-project) GCP_PROJECT=$1 ;;
          --cluster-version) E2E_CLUSTER_VERSION=$1 ;;
          --cluster-creation-flag) EXTRA_CLUSTER_CREATION_FLAGS+=("$1") ;;
          *) abort "unknown option ${parameter}" ;;
        esac
    esac
    shift
  done

  (( IS_PROW )) && [[ -z "${GCP_PROJECT_ID:-}" ]] && IS_BOSKOS=1

  (( SKIP_ISTIO_ADDON )) || GKE_ADDONS="--addons=Istio"

  readonly RUN_TESTS
  readonly GCP_PROJECT
  readonly IS_BOSKOS
  readonly EXTRA_CLUSTER_CREATION_FLAGS
  readonly SKIP_KNATIVE_SETUP
  readonly SKIP_TEARDOWNS
  readonly GKE_ADDONS

  if (( ! RUN_TESTS )); then
    create_test_cluster
  else
    setup_test_cluster
  fi
}
