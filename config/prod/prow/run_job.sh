#!/usr/bin/env bash

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

# Simple script to start a Prow job.

source $(dirname $0)/../../../scripts/library.sh

# First parameter is expected to be the prow job name
JOB_NAME="$1"
# Second parameter is optional, need to be supplied for running pull- and post-
# jobs. Create a token following
# https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line,
# with `repo` permission, save it to a file, and supply it's path as parameter
GITHUB_TOKEN_PATH="$2"

[[ -z "${JOB_NAME}" ]] && abort "pass the name of the job to start as argument"

set -e

JOB_YAML=$(mktemp)
CONFIG_YAML=${REPO_ROOT_DIR}/config/prod/prow/core/config.yaml
JOB_CONFIG_YAML=${REPO_ROOT_DIR}/config/prod/prow/jobs

if [[ -n "${GITHUB_TOKEN_PATH}" ]]; then
    docker run -i --rm \
        -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${JOB_CONFIG_YAML}:${JOB_CONFIG_YAML}" \
        -v "${GITHUB_TOKEN_PATH}:${GITHUB_TOKEN_PATH}" \
        -w "${PWD}" \
        gcr.io/k8s-prow/mkpj:v20200427-84e5e2b2c \
        "--job=${JOB_NAME}" "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_CONFIG_YAML}" \
        "--github-token-path=${GITHUB_TOKEN_PATH}" \
        > ${JOB_YAML}
else
    failed=0
    docker run -i --rm \
        -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${JOB_CONFIG_YAML}:${JOB_CONFIG_YAML}" \
        -w "${PWD}" \
        gcr.io/k8s-prow/mkpj:v20200427-84e5e2b2c \
        "--job=${JOB_NAME}" "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_CONFIG_YAML}" \
        > ${JOB_YAML} || failed=1

    if (( failed )); then
        echo "ERROR: failed generating job config using mkpj"
        echo "It's possible that GitHub token is missing, please try passing github token path as second parameter"
        exit 1
    fi
fi

echo "Job YAML file saved to ${JOB_YAML}"

make -C "${REPO_ROOT_DIR}/config/prod" get-cluster-credentials
kubectl apply -f ${JOB_YAML}
make -C "${REPO_ROOT_DIR}/config/prod" unset-cluster-credentials
