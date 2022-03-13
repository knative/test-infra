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

set -Eeuo pipefail

REPO_ROOT_DIR="$(dirname "$(dirname "$(realpath "${BASH_SOURCE[0]}")")")"

# Generate Prow configs since we are using generator
readonly CONFIG_GENERATOR_DIR="${REPO_ROOT_DIR}/tools/configgen"

# Clean up existing generated config files.
rm -rf "${REPO_ROOT_DIR}/prow/jobs/generated/*"

# Generate config for Prow jobs and TestGrid
cd "${CONFIG_GENERATOR_DIR}" && go run . \
    --prow-jobs-config-input="${REPO_ROOT_DIR}/prow/jobs_config" \
    --prow-jobs-config-output="${REPO_ROOT_DIR}/prow/jobs/generated" \
    --all-prow-jobs-config="${REPO_ROOT_DIR}/prow/jobs" \
    --testgrid-config-output="${REPO_ROOT_DIR}/config/prow/k8s-testgrid/k8s-testgrid.yaml"
