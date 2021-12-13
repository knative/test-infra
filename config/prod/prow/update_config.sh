#!/usr/bin/env bash

# Copyright 2021 The Knative Authors
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

# This script updates Prow job configs.

source $(dirname $0)/../../../scripts/library.sh

set -e

orig_kubeconfig=${KUBECONFIG}
KUBECONFIG="/tmp/kubeconfig.yaml"
make -C "${REPO_ROOT_DIR}/config/prod" get-cluster-credentials
add_trap "make -C ${REPO_ROOT_DIR}/config/prod unset-cluster-credentials" EXIT
add_trap "rm ${KUBECONFIG}" EXIT
add_trap "KUBECONFIG=${orig_kubeconfig}" EXIT

# Download prow core config from prow
if (( IS_OSX )); then
  # On OS X, the file has to be under /private other it cannot be mounted by the container.
  # See https://stackoverflow.com/questions/45122459/docker-mounts-denied-the-paths-are-not-shared-from-os-x-and-are-not-known/45123074
  CONFIG_YAML="/private$(mktemp)"
  PLUGINS_YAML="/private$(mktemp)"
else
  CONFIG_YAML="$(mktemp)"
  PLUGINS_YAML="$(mktemp)"
fi

JOB_YAML="${REPO_ROOT_DIR}/config/prod/prow/jobs"

kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"
kubectl get configmaps plugins -o "jsonpath={.data['plugins\.yaml']}" > "${PLUGINS_YAML}"
echo "Prow plugins downloaded at ${PLUGINS_YAML}"

docker run -i --rm \
    -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${PLUGINS_YAML}:${PLUGINS_YAML}" -v "${JOB_YAML}:${JOB_YAML}" \
    -v "${KUBECONFIG}:${KUBECONFIG}:ro" \
    -w "${PWD}" \
    gcr.io/k8s-prow/config-bootstrapper:v20211213-d726a04d2c \
    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_YAML}" \
    "--plugin-config=${PLUGINS_YAML}" \
    "--source-path=${REPO_ROOT_DIR}" \
    "--kubeconfig=${KUBECONFIG}" \
    "--dry-run=false"
