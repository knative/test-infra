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

set -o errexit
set -o nounset

source $(dirname "$0")/../scripts/library.sh

# Generate Prow configs since we are using generator
readonly CONFIG_GENERATOR_DIR="${REPO_ROOT_DIR}/tools/config-generator"
readonly CONFIG_DIR="${REPO_ROOT_DIR}/config"

# Generate config for production Prow
go run "${CONFIG_GENERATOR_DIR}" \
    --gcs-bucket="knative-prow" \
    --generate-testgrid-config=true \
    --image-docker=gcr.io/knative-tests/test-infra \
    --prow-host=https://prow.knative.dev \
    --testgrid-gcs-bucket="knative-testgrid" \
    --prow-jobs-config-output="${CONFIG_DIR}/prod/prow/jobs/config.yaml" \
    --testgrid-config-output="${CONFIG_DIR}/prod/prow/testgrid/testgrid.yaml" \
    "${CONFIG_DIR}/prod/prow/config_knative.yaml"
