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

# Runs prow/pj-on-kind.sh with config arguments specific to the prow.knative.dev instance.
# Requries go, docker, and kubectl.

# Documentation: https://github.com/kubernetes/test-infra/blob/master/prow/build_test_update.md#using-pj-on-kindsh
# Example usage:
# ./pj-on-kind.sh pull-knative-test-infra-unit-tests

set -o errexit
set -o nounset
set -o pipefail

PROW_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

# Download prow core config from prow
if [[ $IS_OSX ]]; then
  # On OS X, the file has to be under /private other it cannot be mounted by the container.
  # See https://stackoverflow.com/questions/45122459/docker-mounts-denied-the-paths-are-not-shared-from-os-x-and-are-not-known/45123074
  CONFIG_YAML="/private"$(mktemp)
else
  CONFIG_YAML=$(mktemp)
fi
make -C "${PROW_DIR}/.." get-cluster-credentials
trap "make -C '${PROW_DIR}/..' unset-cluster-credentials" EXIT
kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"

export CONFIG_YAML
export JOB_CONFIG_PATH="${PROW_DIR}/jobs"

bash <(curl -sSfL https://raw.githubusercontent.com/kubernetes/test-infra/master/prow/pj-on-kind.sh) "$@"
