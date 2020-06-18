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

# This script checks for consistency among prow core config, plugins, and job
# configs

source $(dirname $0)/../../../scripts/library.sh

set -e

# Settings for inrepoconfig checking, if only `REPO_NAME` is supplied then:
# `REPO_YAML_PATH` is default to `("/home/prow/go/src/github.com/%s/.prow.yaml",
# o.prowYAMLRepoName)`. See
# https://github.com/kubernetes/test-infra/blob/a3f67762b39ce6eeb221b4d475cf4a1a32ac0c54/prow/cmd/checkconfig/main.go#L144
REPO_NAME_TO_CHECK="${1:-}"
REPO_YAML_PATH_TO_CHECK="${2:-}"

make -C "${REPO_ROOT_DIR}/config/prod" get-cluster-credentials
trap "make -C '${REPO_ROOT_DIR}/config/prod' unset-cluster-credentials" EXIT

CONFIG_YAML="$(mktemp)"
PLUGINS_YAML="$(mktemp)"
JOB_YAML="${REPO_ROOT_DIR}/config/prod/prow/jobs"

kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"
kubectl get configmaps plugins -o "jsonpath={.data['plugins\.yaml']}" > "${PLUGINS_YAML}"
echo "Prow plugins downloaded at ${PLUGINS_YAML}"

REPO_NAME_ARG=""
REPO_YAML_PATH_ARG=""
[[ -n "${REPO_NAME_TO_CHECK}" ]] && REPO_NAME_ARG="--prow-yaml-repo-name=${REPO_NAME_TO_CHECK}"
[[ -n "${REPO_YAML_PATH_TO_CHECK}" ]] && REPO_YAML_PATH_ARG="--prow-yaml-repo-path=${REPO_YAML_PATH_TO_CHECK}"

docker run -i --rm \
    -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${PLUGINS_YAML}:${PLUGINS_YAML}" -v "${JOB_CONFIG_YAML}:${JOB_CONFIG_YAML}" \
    -w "${PWD}" \
    gcr.io/k8s-prow/checkconfig:v20200603-4badfd9f37 \
    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_CONFIG_YAML}" \
    "--plugin-config=${PLUGINS_YAML}" "--strict" "--exclude-warning=mismatched-tide" \
    "${REPO_NAME_ARG}" "${REPO_YAML_PATH_ARG}"
