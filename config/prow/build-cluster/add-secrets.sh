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

# Requries gcloud and kubectl. Assume the cluster contexts are already set for
# both service and build clusters

KUBECONFIG_SERVICE_CLUSTER="gke_knative-tests_us-central1-f_prow"
KUBECONFIG_BUILD_CLUSTER="gke_knative-tests_us-central1-f_knative-prow-build-cluster"

secrets=(
    "test-pods;covbot-token;GitHub token for the coverage job."
    "test-pods;flaky-test-reporter-github-token;GitHub token for the flaky test reporter job."
    "test-pods;flaky-test-reporter-slack-token;Slack token the flaky test reporter job."
    "test-pods;housekeeping-github-token;GitHub token the issue tracker job."
    "test-pods;hub-token;GitHub token for the release job."
    "test-pods;nightly-account;Service account for the nightly jobs."
    "test-pods;performance-test;Service account for the performance tests."
    "test-pods;prow-auto-bumper-github-token;GitHub token the Prow updater job."
    "test-pods;prow-updater-robot-ssh-key;SSH key used by the Prow updater job."
    "test-pods;release-account;Service account for the release jobs."
    "test-pods;repoview-token;GitHub token the presubmit jobs (repo view only)."
    "test-pods;test-account;Service account for the tests."
)

for secret in "${secrets[@]}"; do
    IFS=';' # space is set as delimiter
    read -ra PARTS <<< "$secret" # secret is read into an array as tokens separated by IFS
    namespace="${PARTS[0]}"
    secret_name="${PARTS[1]}"
    secret_desc="${PARTS[2]}"

    echo "Copy $secret_name from namespace $namespace"

    kubectl get secret "${secret_name}" -n "${namespace}" --context "${KUBECONFIG_SERVICE_CLUSTER}" --export -o yaml | \
        kubectl apply -n "${namespace}" --context "${KUBECONFIG_BUILD_CLUSTER}" -f -
done
