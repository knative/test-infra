#!/bin/bash

# Copyright 2018 The Knative Authors
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

# This is a helper script for Knative E2E test scripts. To use it:
# 1. Source this script.
# 2. Write the teardown() function, which will tear down your test resources.
# 3. [optional] Write the dump_extra_cluster_state() function. It will be called
#    when a test fails, and can dump extra information about the current state of
#    the cluster (tipically using kubectl).
# 4. Call the initialize() function passing $@ (without quotes).
# 5. Write logic for the end-to-end tests. Run all go tests using report_go_test()
#    and call fail_test() or success() if any of them failed. The envitronment
#    variables DOCKER_REPO_OVERRIDE, K8S_CLUSTER_OVERRIDE and K8S_USER_OVERRIDE
#    will be set accordingly to the test cluster. You can also use the following
#    boolean (0 is false, 1 is true) environment variables for the logic:
#    EMIT_METRICS: true if --emit-metrics is passed.
#    USING_EXISTING_CLUSTER: true if the test cluster is an already existing one,
#                            and not a temporary cluster created by kubetest.
#    All environment variables above are marked read-only.
# Notes:
# 1. Calling your script without arguments will create a new cluster in the GCP
#    project $PROJECT_ID and run the tests against it.
# 2. Calling your script with --run-tests and the variables K8S_CLUSTER_OVERRIDE,
#    K8S_USER_OVERRIDE and DOCKER_REPO_OVERRIDE set will immediately start the
#    tests against the cluster.

# Load github.com/knative/test-infra/images/prow-tests/scripts/library.sh
[ -f /workspace/library.sh ] \
  && source /workspace/library.sh \
  || eval "$(docker run --entrypoint sh gcr.io/knative-tests/test-infra/prow-tests -c 'cat library.sh')"
[ -v KNATIVE_TEST_INFRA ] || exit 1

# Test cluster parameters
readonly E2E_BASE_NAME=$(basename ${REPO_ROOT_DIR})
readonly E2E_CLUSTER_NAME=${E2E_BASE_NAME}-e2e-cls${BUILD_NUMBER}
readonly E2E_NETWORK_NAME=${E2E_BASE_NAME}-e2e-net${BUILD_NUMBER}
readonly E2E_CLUSTER_ZONE=us-central1-a
readonly E2E_CLUSTER_NODES=3
readonly E2E_CLUSTER_MACHINE=n1-standard-4
readonly TEST_RESULT_FILE=/tmp/${E2E_BASE_NAME}-e2e-result

# Tear down the test resources.
function teardown_test_resources() {
  header "Tearing down test environment"
  # Free resources in GCP project.
  if (( ! USING_EXISTING_CLUSTER )); then
    teardown
  fi

  # Delete Knative Serving images when using prow.
  if (( IS_PROW )); then
    echo "Images in ${DOCKER_REPO_OVERRIDE}:"
    gcloud container images list --repository=${DOCKER_REPO_OVERRIDE}
    delete_gcr_images ${DOCKER_REPO_OVERRIDE}
  else
    # Delete the kubernetes source downloaded by kubetest
    rm -fr kubernetes kubernetes.tar.gz
  fi
}

# Exit test if the previous command failed, dumping current state info.
# Parameters: $1 - error message (optional).
function fail_test() {
  [[ $? -eq 0 ]] && return 0
  [[ -n $1 ]] && echo "ERROR: $1"
  dump_k8s_info
  exit 1
}

# Dump info about the test cluster. If dump_extra_cluster_info() is defined, calls it too.
# This is intended to be called when a test fails to provide debugging information.
function dump_cluster_state() {
  echo "***************************************"
  echo "***           TEST FAILED           ***"
  echo "***    Start of information dump    ***"
  echo "***************************************"
  echo ">>> All resources:"
  kubectl get all --all-namespaces
  echo ">>> Services:"
  kubectl get services --all-namespaces
  echo ">>> Events:"
  kubectl get events --all-namespaces
  [[ "$(type -t dump_extra_cluster_state)" == "function" ]] && dump_extra_cluster_state
  echo "***************************************"
  echo "***           TEST FAILED           ***"
  echo "***     End of information dump     ***"
  echo "***************************************"
}

# Create a test cluster with kubetest and call the current script again.
function create_test_cluster() {
  header "Creating test cluster"
  # Smallest cluster required to run the end-to-end-tests
  local CLUSTER_CREATION_ARGS=(
    --gke-create-args="--enable-autoscaling --min-nodes=1 --max-nodes=${E2E_CLUSTER_NODES} --scopes=cloud-platform"
    --gke-shape={\"default\":{\"Nodes\":${E2E_CLUSTER_NODES}\,\"MachineType\":\"${E2E_CLUSTER_MACHINE}\"}}
    --provider=gke
    --deployment=gke
    --cluster="${E2E_CLUSTER_NAME}"
    --gcp-zone="${E2E_CLUSTER_ZONE}"
    --gcp-network="${E2E_NETWORK_NAME}"
    --gke-environment=prod
  )
  if (( ! IS_PROW )); then
    CLUSTER_CREATION_ARGS+=(--gcp-project=${PROJECT_ID:?"PROJECT_ID must be set to the GCP project where the tests are run."})
  fi
  # SSH keys are not used, but kubetest checks for their existence.
  # Touch them so if they don't exist, empty files are create to satisfy the check.
  touch $HOME/.ssh/google_compute_engine.pub
  touch $HOME/.ssh/google_compute_engine
  # Clear user and cluster variables, so they'll be set to the test cluster.
  # DOCKER_REPO_OVERRIDE is not touched because when running locally it must
  # be a writeable docker repo.
  export K8S_USER_OVERRIDE=
  export K8S_CLUSTER_OVERRIDE=
  # Assume test failed (see more details at the end of this script).
  echo -n "1"> ${TEST_RESULT_FILE}
  local test_cmd_args="--run-tests"
  (( EMIT_METRICS )) && test_cmd_args+=" --emit-metrics"
  # Normalize script path; we can't use readlink because it's not available everywhere
  local script=$0
  [[ ${script} =~ ^[\./].* ]] || script="./$0"
  script="$(cd ${script%/*} && echo $PWD/${script##*/})"
  kubetest "${CLUSTER_CREATION_ARGS[@]}" \
    --up \
    --down \
    --extract "gke-${SERVING_GKE_VERSION}" \
    --gcp-node-image ${SERVING_GKE_IMAGE} \
    --test-cmd "${script}" \
    --test-cmd-args "${test_cmd_args}"
  # Delete target pools and health checks that might have leaked.
  # See https://github.com/knative/serving/issues/959 for details.
  # TODO(adrcunha): Remove once the leak issue is resolved.
  local gcp_project=${PROJECT_ID}
  [[ -z ${gcp_project} ]] && gcp_project=$(gcloud config get-value project)
  local http_health_checks="$(gcloud compute target-pools list \
    --project=${gcp_project} --format='value(healthChecks)' --filter="instances~-${E2E_CLUSTER_NAME}-" | \
    grep httpHealthChecks | tr '\n' ' ')"
  local target_pools="$(gcloud compute target-pools list \
    --project=${gcp_project} --format='value(name)' --filter="instances~-${E2E_CLUSTER_NAME}-" | \
    tr '\n' ' ')"
  local region="$(gcloud compute zones list --filter=name=${E2E_CLUSTER_ZONE} --format='value(region)')"
  if [[ -n "${target_pools}" ]]; then
    echo "Found leaked target pools, deleting"
    gcloud compute forwarding-rules delete -q --project=${gcp_project} --region=${region} ${target_pools}
    gcloud compute target-pools delete -q --project=${gcp_project} --region=${region} ${target_pools}
  fi
  if [[ -n "${http_health_checks}" ]]; then
    echo "Found leaked health checks, deleting"
    gcloud compute http-health-checks delete -q --project=${gcp_project} ${http_health_checks}
  fi
  local result="$(cat ${TEST_RESULT_FILE})"
  echo "Test result code is $result"
  exit ${result}
}

# Setup the test cluster for running the tests.
function setup_test_cluster() {
  # Fail fast during setup.
  set -o errexit
  set -o pipefail

  # Set the required variables if necessary.
  if [[ -z ${K8S_USER_OVERRIDE} ]]; then
    export K8S_USER_OVERRIDE=$(gcloud config get-value core/account)
  fi

  if [[ -z ${K8S_CLUSTER_OVERRIDE} ]]; then
    USING_EXISTING_CLUSTER=0
    export K8S_CLUSTER_OVERRIDE=$(kubectl config current-context)
    acquire_cluster_admin_role ${K8S_USER_OVERRIDE} ${E2E_CLUSTER_NAME} ${E2E_CLUSTER_ZONE}
    # Make sure we're in the default namespace. Currently kubetest switches to
    # test-pods namespace when creating the cluster.
    kubectl config set-context $K8S_CLUSTER_OVERRIDE --namespace=default
  fi
  readonly USING_EXISTING_CLUSTER

  if [[ -z ${DOCKER_REPO_OVERRIDE} ]]; then
    export DOCKER_REPO_OVERRIDE=gcr.io/$(gcloud config get-value project)/${E2E_BASE_NAME}e2e-img
  fi

  echo "- Cluster is ${K8S_CLUSTER_OVERRIDE}"
  echo "- User is ${K8S_USER_OVERRIDE}"
  echo "- Docker is ${DOCKER_REPO_OVERRIDE}"

  trap teardown_test_resources EXIT

  if (( USING_EXISTING_CLUSTER )); then
    echo "Deleting any previous Knative Serving instance"
    teardown
  fi

  readonly K8S_CLUSTER_OVERRIDE
  readonly K8S_USER_OVERRIDE
  readonly DOCKER_REPO_OVERRIDE

  # Handle failures ourselves, so we can dump useful info.
  set +o errexit
  set +o pipefail
}

function success() {
  # kubetest teardown might fail and thus incorrectly report failure of the
  # script, even if the tests pass.
  # We store the real test result to return it later, ignoring any teardown
  # failure in kubetest.
  # TODO(adrcunha): Get rid of this workaround.
  echo -n "0"> ${TEST_RESULT_FILE}
  echo "**************************************"
  echo "***        ALL TESTS PASSED        ***"
  echo "**************************************"
  exit 0
}

RUN_TESTS=0
EMIT_METRICS=0
USING_EXISTING_CLUSTER=1

# Parse flags and initialize the test cluster.
function initialize() {
  cd ${REPO_ROOT_DIR}
  for parameter in $@; do
    case $parameter in
      --run-tests) RUN_TESTS=1 ;;
      --emit-metrics) EMIT_METRICS=1 ;;
      *)
        echo "error: unknown option ${parameter}"
        echo "usage: $0 [--run-tests][--emit-metrics]"
        exit 1
        ;;
    esac
    shift
  done
  readonly RUN_TESTS
  readonly EMIT_METRICS

  if (( ! RUN_TESTS )); then
    create_test_cluster
  else
    setup_test_cluster
  fi
}
