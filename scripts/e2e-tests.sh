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

# Test cluster parameters

# Configurable parameters
# export E2E_CLUSTER_REGION as they're used in the cluster setup subprocess
export E2E_CLUSTER_REGION=${E2E_CLUSTER_REGION:-us-central1}

readonly E2E_CLUSTER_MACHINE=${E2E_CLUSTER_MACHINE:-e2-standard-4}
readonly E2E_GKE_ENVIRONMENT=${E2E_GKE_ENVIRONMENT:-prod}
readonly E2E_GKE_COMMAND_GROUP=${E2E_GKE_COMMAND_GROUP:-beta}

# Each knative repository may have a different cluster size requirement here,
# so we allow calling code to set these parameters.  If they are not set we
# use some sane defaults.
readonly E2E_MIN_CLUSTER_NODES=${E2E_MIN_CLUSTER_NODES:-1}
readonly E2E_MAX_CLUSTER_NODES=${E2E_MAX_CLUSTER_NODES:-3}

readonly E2E_BASE_NAME="k${REPO_NAME}"

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
  local test_options=""
  local go_options=""
  [[ ! " $@" == *" -tags="* ]] && go_options="-tags=e2e"
  report_go_test -v -race -count=1 ${go_options} $@ "${test_options}"
}

# Dumps the k8s api server metrics. Spins up a proxy, waits a little bit and
# dumps the metrics to ${ARTIFACTS}/k8s.metrics.txt
function dump_metrics() {
  header ">> Starting kube proxy"
  kubectl proxy --port=8080 &
  local proxy_pid=$!
  sleep 5
  header ">> Grabbing k8s metrics"
  curl -s http://localhost:8080/metrics > "${ARTIFACTS}"/k8s.metrics.txt
  # Clean up proxy so it doesn't interfere with job shutting down
  kill $proxy_pid || true
}

# Dump info about the test cluster. If dump_extra_cluster_info() is defined, calls it too.
# This is intended to be called when a test fails to provide debugging information.
function dump_cluster_state() {
  echo "***************************************"
  echo "***         E2E TEST FAILED         ***"
  echo "***    Start of information dump    ***"
  echo "***************************************"

  local output
  output="${ARTIFACTS}/k8s.dump-$(basename "${E2E_SCRIPT}").txt"
  echo ">>> The dump is located at ${output}"

  for crd in $(kubectl api-resources --verbs=list -o name | sort); do
    local count
    count="$(kubectl get "$crd" --all-namespaces --no-headers 2>/dev/null | wc -l)"
    echo ">>> ${crd} (${count} objects)"
    if [[ "${count}" -gt 0 ]]; then
      {
        echo ">>> ${crd} (${count} objects)"
        echo ">>> Listing"
        kubectl get "${crd}" --all-namespaces
        echo ">>> Details"
      } >> "${output}"
      if [[ "${crd}" == "secrets" ]]; then
        echo "Secrets are ignored for security reasons" >> "${output}"
      elif [[ "${crd}" == "events" ]]; then
        echo "events are ignored as making a lot of noise" >> "${output}"
      else
        kubectl get "${crd}" --all-namespaces -o yaml >> "${output}"
      fi
    fi
  done

  if function_exists dump_extra_cluster_state; then
    echo ">>> Extra dump" >> "${output}"
    dump_extra_cluster_state >> "${output}"
  fi
  echo "***************************************"
  echo "***         E2E TEST FAILED         ***"
  echo "***     End of information dump     ***"
  echo "***************************************"
}

# Read metadata.json and get value for key
# Parameters: $1 - Key for metadata
function get_meta_value() {
  run_kntest metadata get --key "$1"
}

# Override create_test_cluster in scripts/e2e-tests.sh
# Create test cluster with cluster creation lib and write metadata in ${ARTIFACT}/metadata.json
function create_test_cluster() {
  # Fail fast during setup.
  set -o errexit
  set -o pipefail

  if function_exists cluster_setup; then
    cluster_setup || fail_test "cluster setup failed"
  fi

  header "Creating test cluster"
  local creation_args="--save-meta-data"
  creation_args+=" --min-nodes=${E2E_MIN_CLUSTER_NODES} --max-nodes=${E2E_MAX_CLUSTER_NODES} --node-type=${E2E_CLUSTER_MACHINE}"
  creation_args+=" --region=${E2E_CLUSTER_REGION}"
  creation_args+=" --version=${E2E_CLUSTER_VERSION}"
  creation_args+=" --addons HttpLoadBalancing,HorizontalPodAutoscaling"
  (( ! SKIP_ISTIO_ADDON )) && creation_args+=",Istio"
  [[ -n "${GCP_PROJECT}" ]] && creation_args+=" --project ${GCP_PROJECT}"
  creation_args+=" ${EXTRA_CLUSTER_CREATION_FLAGS}"
  # TODO(chizhg): support parameterizing "gke" so that we can create other types of clusters
  run_kntest cluster gke create "${creation_args}" || fail_test "failed creating test cluster"

  # Since calling `create_test_cluster` assumes cluster creation, removing
  # cluster afterwards.
  add_trap "run_kntest cluster gke delete > /dev/null &" EXIT SIGINT

  setup_test_cluster || fail_test "failed setting up test cluster"

  set +o errexit
  set +o pipefail
}

# Override setup_test_cluster, this function is almost copy paste from original
# setup_test_cluster function, other than reading metadata from
# "${ARTIFACTS}/metadata.json" instead of from global vars
function setup_test_cluster() {
  # Fail fast during setup.
  set -o errexit
  set -o pipefail

  header "Setting up test cluster"
  # Run kntest to acquire the existing test cluster, will fail if
  # kubeconfig isn't set or the cluster doesn't exist.
  run_kntest cluster gke get --save-meta-data || fail_test "failed getting test cluster"
  # The step above collects cluster metadata and writes to
  # ${ARTIFACTS}/metadata.json file, use this information.
  echo "Cluster used for running tests: $(cat "${ARTIFACTS}"/metadata.json)"
  local e2e_cluster_name e2e_cluster_region e2e_cluster_zone e2e_project_name
  e2e_cluster_name=$(get_meta_value "E2E:Machine")
  e2e_cluster_region=$(get_meta_value "E2E:Region")
  e2e_cluster_zone=$(get_meta_value "E2E:Zone")
  e2e_project_name=$(get_meta_value "E2E:Project")

  # Set the actual project the test cluster resides in
  # It will be a project assigned by Boskos if test is running on Prow,
  # otherwise will be ${GCP_PROJECT} set up by user.
  export E2E_PROJECT_ID=$e2e_project_name
  readonly E2E_PROJECT_ID

  local k8s_user
  k8s_user=$(gcloud config get-value core/account)
  local k8s_cluster
  # current-context must have been set at this point.
  k8s_cluster=$(kubectl config current-context)
  # If cluster admin role isn't set, this is a brand new cluster.
  # Setup the admin role and also KO_DOCKER_REPO if it is a GKE cluster.
  if [[ -z "$(kubectl get clusterrolebinding cluster-admin-binding 2> /dev/null)" && "${k8s_cluster}" =~ ^gke_.* ]]; then
    acquire_cluster_admin_role "${k8s_user}" "${e2e_cluster_name}" "${e2e_cluster_region}" "${e2e_cluster_zone}"
    # Incorporate an element of randomness to ensure that each run properly publishes images.
    export KO_DOCKER_REPO=gcr.io/${E2E_PROJECT_ID}/${E2E_BASE_NAME}-e2e-img/${RANDOM}
  fi

  # Safety checks
  is_protected_gcr "${KO_DOCKER_REPO}" && \
    abort "\$KO_DOCKER_REPO set to ${KO_DOCKER_REPO}, which is forbidden"

  echo "- gcloud project is ${E2E_PROJECT_ID}"
  echo "- gcloud user is ${k8s_user}"
  echo "- Cluster is ${k8s_cluster}"
  echo "- Docker repository is ${KO_DOCKER_REPO}"

  export KO_DATA_PATH="${REPO_ROOT_DIR}/.git"

  # Do not run teardowns if we explicitly want to skip them.
  (( ! SKIP_TEARDOWNS )) && add_trap teardown_test_resources EXIT SIGINT

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
E2E_CLUSTER_VERSION=""
GKE_ADDONS=""
EXTRA_CLUSTER_CREATION_FLAGS=()
E2E_SCRIPT_CUSTOM_FLAGS=()

# Parse flags and initialize the test cluster.
function initialize() {
  E2E_SCRIPT="$(get_canonical_path "$0")"
  E2E_CLUSTER_VERSION="${SERVING_GKE_VERSION}"

  cd "${REPO_ROOT_DIR}"
  while [[ $# -ne 0 ]]; do
    local parameter=$1
    # Try parsing flag as a custom one.
    if function_exists parse_flags; then
      parse_flags "$@"
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
          --cluster-creation-flag) EXTRA_CLUSTER_CREATION_FLAGS+=($1) ;;
          *) abort "unknown option ${parameter}" ;;
        esac
    esac
    shift
  done

  # Use PROJECT_ID if set, unless --gcp-project was used.
  if [[ -n "${PROJECT_ID:-}" && -z "${GCP_PROJECT}" ]]; then
    echo "\$PROJECT_ID is set to '${PROJECT_ID}', using it to run the tests"
    GCP_PROJECT="${PROJECT_ID}"
  fi
  if (( ! IS_PROW )) && (( ! RUN_TESTS )) && [[ -z "${GCP_PROJECT}" ]]; then
    abort "set \$PROJECT_ID or use --gcp-project to select the GCP project where the tests are run"
  fi

  (( IS_PROW )) && [[ -z "${GCP_PROJECT}" ]] && IS_BOSKOS=1

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
