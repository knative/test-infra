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

source $(dirname $0)/../../../../scripts/library.sh

set -e

make -C "${REPO_ROOT_DIR}/config/prod" get-cluster-credentials
trap "make -C '${REPO_ROOT_DIR}/config/prod' unset-cluster-credentials" EXIT

PROW_CONFIG="$(mktemp)"
PROW_JOB_CONFIG="${REPO_ROOT_DIR}/config/prod/prow/jobs/config.yaml"
TESTGRID_YAML="${REPO_ROOT_DIR}/config/prod/prow/k8s-testgrid/k8s-testgrid.yaml"

kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"

docker run -i --rm \
    -v "${PWD}:${PWD}" \
    -v "${PROW_CONFIG}:${PROW_CONFIG}" \
    -v "${PROW_JOB_CONFIG}:${PROW_JOB_CONFIG}" \
    -v "${TESTGRID_YAML}:${TESTGRID_YAML}" \
    -w "${PWD}" \
    gcr.io/k8s-prow/transfigure:v20201110-9e512b5af0 \
    "/etc/github-token/oauth" \
    "${PROW_CONFIG}" \
    "${PROW_JOB_CONFIG}" \
    "${TESTGRID_YAML}" \
    "knative"
