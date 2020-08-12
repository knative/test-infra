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

# This is a collection of functions for infra related setups, mainly
# cluster provisioning. It doesn't do anything when called from command line.

source $(dirname "${BASH_SOURCE[0]}")/library.sh

# Test cluster parameters

# Configurable parameters
# export E2E_GKE_CLUSTER_REGION and E2E_GKE_CLUSTER_ZONE as they're used in the cluster setup subprocess
export E2E_GKE_CLUSTER_REGION=${E2E_GKE_CLUSTER_REGION:-us-central1}
# By default we use regional clusters.
export E2E_GKE_CLUSTER_ZONE=${E2E_GKE_CLUSTER_ZONE:-}

# Default backup regions in case of stockouts; by default we don't fall back to a different zone in the same region
readonly E2E_GKE_CLUSTER_BACKUP_REGIONS=${E2E_GKE_CLUSTER_BACKUP_REGIONS:-us-west1 us-east1}
readonly E2E_GKE_CLUSTER_BACKUP_ZONES=${E2E_GKE_CLUSTER_BACKUP_ZONES:-}

readonly E2E_GKE_ENVIRONMENT=${E2E_GKE_ENVIRONMENT:-prod}
readonly E2E_GKE_COMMAND_GROUP=${E2E_GKE_COMMAND_GROUP:-beta}
readonly E2E_GKE_CLUSTER_MACHINE=${E2E_GKE_CLUSTER_MACHINE:-e2-standard-4}
readonly E2E_GKE_ADDONS=${E2E_GKE_ADDONS:-}

# Each knative repository may have a different cluster size requirement here,
# so we allow calling code to set these parameters.  If they are not set we
# use some sane defaults.
readonly E2E_MIN_CLUSTER_NODES=${E2E_MIN_CLUSTER_NODES:-1}
readonly E2E_MAX_CLUSTER_NODES=${E2E_MAX_CLUSTER_NODES:-3}

readonly E2E_GKE_NETWORK_NAME=$(build_resource_name e2e-net)
readonly E2E_CLUSTER_NAME=$(build_resource_name e2e-cls)

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
    if [[ "${count}" -gt "0" ]]; then
      {
        echo ">>> ${crd} (${count} objects)"

        echo ">>> Listing"
        kubectl get "${crd}" --all-namespaces

        echo ">>> Details"
        if [[ "${crd}" == "secrets" ]]; then
          echo "Secrets are ignored for security reasons"
        elif [[ "${crd}" == "events" ]]; then
          echo "events are ignored as making a lot of noise"
        else
          kubectl get "${crd}" --all-namespaces -o yaml
        fi
      } >> "${output}"
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

# On a Prow job, save some metadata about the test for Testgrid.
function save_metadata() {
  (( ! IS_PROW )) && return
  local geo_key="Region"
  local geo_value="${E2E_GKE_CLUSTER_REGION}"
  if [[ -n "${E2E_GKE_CLUSTER_ZONE}" ]]; then
    geo_key="Zone"
    geo_value="${E2E_GKE_CLUSTER_REGION}-${E2E_GKE_CLUSTER_ZONE}"
  fi
  local cluster_version
  cluster_version="$(gcloud container clusters list --project="${E2E_PROJECT_ID}" --format='value(currentMasterVersion)')"
  run_kntest metadata set --key="E2E:${geo_key}" --value="${geo_value}"
  run_kntest metadata set --key="E2E:Machine" --value="${E2E_GKE_CLUSTER_MACHINE}"
  run_kntest metadata set --key="E2E:Version" --value="${cluster_version}"
  run_kntest metadata set --key="E2E:MinNodes" --value="${E2E_MIN_CLUSTER_NODES}"
  run_kntest metadata set --key="E2E:MaxNodes" --value="${E2E_MAX_CLUSTER_NODES}"
}

# Create a GKE test cluster with kubetest2 and call the current script again.
# Parameters: $1 - extra cluster creation flags
#             $2 - extra kubetest2 flags
#             $3 - test command to run by the kubetest2 tester
function create_gke_test_cluster() {
  # Fail fast during setup.
  set -o errexit
  set -o pipefail

  local extra_cluster_creation_flags
  IFS=' ' read -r -a extra_cluster_creation_flags <<< "$1"

  if function_exists cluster_setup; then
    cluster_setup || fail_test "cluster setup failed"
  fi

  echo "Cluster will have a minimum of ${E2E_MIN_CLUSTER_NODES} and a maximum of ${E2E_MAX_CLUSTER_NODES} nodes."
  local CLUSTER_CREATION_ARGS=(
    "gke"
    "--create-command=${E2E_GKE_COMMAND_GROUP} container clusters create --quiet --enable-autoscaling
      --min-nodes=${E2E_MIN_CLUSTER_NODES} --max-nodes=${E2E_MAX_CLUSTER_NODES}
      --cluster-version=${E2E_CLUSTER_VERSION}
      --scopes=cloud-platform --enable-basic-auth --no-issue-client-certificate
      --no-enable-ip-alias --no-enable-autoupgrade
      --addons=${GKE_ADDONS} ${extra_cluster_creation_flags[@]}"
    "--environment=${E2E_GKE_ENVIRONMENT}"
    "--cluster-name=${E2E_CLUSTER_NAME}"
    "--num-nodes=${E2E_MIN_CLUSTER_NODES}"
    "--machine-type=${E2E_CLUSTER_MACHINE}"
    "--network=${E2E_GKE_NETWORK_NAME}"
    "--require-gcp-ssh-key=false"
    --up
  )
  if (( ! IS_BOSKOS )); then
    CLUSTER_CREATION_ARGS+=("--gcp-project=${GCP_PROJECT}")
  fi

  local gcloud_project="${GCP_PROJECT}"
  [[ -z "${gcloud_project}" ]] && gcloud_project="$(gcloud config get-value project)"
  echo "gcloud project is ${gcloud_project}"
  echo "gcloud user is $(gcloud config get-value core/account)"
  (( IS_BOSKOS )) && echo "Using boskos for the test cluster"
  [[ -n "${GCP_PROJECT}" ]] && echo "GCP project for test cluster is ${GCP_PROJECT}"
  echo "Test script is ${E2E_SCRIPT}"
  # Set arguments for this script again
  local test_cmd_args="--run-tests"
  (( SKIP_KNATIVE_SETUP )) && test_cmd_args+=" --skip-knative-setup"
  [[ -n "${GCP_PROJECT}" ]] && test_cmd_args+=" --project ${GCP_PROJECT}"
  [[ -n "${E2E_SCRIPT_CUSTOM_FLAGS[*]}" ]] && test_cmd_args+=" ${E2E_SCRIPT_CUSTOM_FLAGS[*]}"
  local extra_flags=()
  if (( IS_BOSKOS )); then
    # Add arbitrary duration, wait for Boskos projects acquisition before error out
    extra_flags+=("--boskos-acquire-timeout-seconds=1200")
  elif (( ! SKIP_TEARDOWNS )); then
    # Only let kubetest2 tear down the cluster if not using Boskos and teardowns are not expected to be skipped,
    # it's done by Janitor if using Boskos
    extra_flags+=("--down")
  fi

  # Create cluster and run the tests
  create_gke_test_cluster_with_retries "${E2E_SCRIPT} ${test_cmd_args}" \
    "${CLUSTER_CREATION_ARGS[@]}" "${extra_flags[@]}"
  local result="$?"
  # Ignore any errors below, this is a best-effort cleanup and shouldn't affect the test result.
  set +o errexit
  function_exists cluster_teardown && cluster_teardown
  echo "Artifacts were written to ${ARTIFACTS}"
  echo "Test result code is ${result}"
  exit "${result}"
}

# Retry backup regions/zones if cluster creations failed due to stockout.
# Parameters: $1     - test command to run by the kubetest2 tester
#             $2..$n - any other kubetest2 flags other than geo and cluster version flag.
function create_gke_test_cluster_with_retries() {
  local tester_command
  IFS=' ' read -r -a tester_command <<< "$1"
  local kubetest2_flags=( "${@:2}" )
  local cluster_creation_log=/tmp/${E2E_BASE_NAME}-cluster_creation-log
  # zone_not_provided is a placeholder for e2e_cluster_zone to make for loop below work
  local zone_not_provided="zone_not_provided"

  local e2e_cluster_regions=("${E2E_GKE_CLUSTER_REGION}")
  local e2e_cluster_zones=("${E2E_GKE_CLUSTER_ZONE}")

  if [[ -n "${E2E_GKE_CLUSTER_BACKUP_ZONES}" ]]; then
    e2e_cluster_zones+=("${E2E_GKE_CLUSTER_BACKUP_ZONES}")
  elif [[ -n "${E2E_GKE_CLUSTER_BACKUP_REGIONS}" ]]; then
    e2e_cluster_regions+=("${E2E_GKE_CLUSTER_BACKUP_REGIONS}")
    e2e_cluster_zones=("${zone_not_provided}")
  else
    echo "No backup region/zone set, cluster creation will fail in case of stockout"
  fi

  for e2e_cluster_region in "${e2e_cluster_regions[@]}"; do
    for e2e_cluster_zone in "${e2e_cluster_zones[@]}"; do
      E2E_GKE_CLUSTER_REGION=${e2e_cluster_region}
      E2E_GKE_CLUSTER_ZONE=${e2e_cluster_zone}
      [[ "${E2E_GKE_CLUSTER_ZONE}" == "${zone_not_provided}" ]] && E2E_GKE_CLUSTER_ZONE=""
      local cluster_creation_zone="${E2E_GKE_CLUSTER_REGION}"
      [[ -n "${E2E_GKE_CLUSTER_ZONE}" ]] && cluster_creation_zone="${E2E_GKE_CLUSTER_REGION}-${E2E_GKE_CLUSTER_ZONE}"

      header "Creating test cluster ${E2E_CLUSTER_VERSION} in ${cluster_creation_zone}"
      if run_go_tool k8s-sigs.io/kubetest2 \
        kubetest2 "${kubetest2_flags[@]}" --region="${cluster_creation_zone}" \
          --test=exec -- "${tester_command[@]}" 2>&1 \
         | tee "${cluster_creation_log}"; then
        return 0
      fi
      # Retry if cluster creation failed because of:
      # - stockout (https://github.com/knative/test-infra/issues/592)
      # - latest GKE not available in this region/zone yet (https://github.com/knative/test-infra/issues/694)
      [[ -z "$(grep -Fo 'does not have enough resources available to fulfill' "${cluster_creation_log}")" \
          && -z "$(grep -Fo 'ResponseError: code=400, message=No valid versions with the prefix' "${cluster_creation_log}")" \
          && -z "$(grep -Po 'ResponseError: code=400, message=Master version "[0-9a-z\-\.]+" is unsupported' "${cluster_creation_log}")"  \
          && -z "$(grep -Po 'only \d+ nodes out of \d+ have registered; this is likely due to Nodes failing to start correctly' "${cluster_creation_log}")" ]] \
          && return 1
    done
  done
  echo "No more region/zones to try, quitting"
  return 1
}
