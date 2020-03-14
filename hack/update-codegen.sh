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

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../scripts/library.sh

CONFIG_GENERATOR_DIR="${REPO_ROOT_DIR}/tools/config-generator"
CONFIG_DIR="${REPO_ROOT_DIR}/config/prow"
PROW_GCS="knative-prow"
GENERATE_TESTGRID_CONFIG="true"
GENERATE_MAINTENANCE_JOBS="true"
PROJECT="knative-tests"
PROW_HOST="https://prow.knative.dev"
TESTGRID_GCS="knative-testgrid"
KNATIVE_CONFIG="config_knative.yaml"

function generate_config() {
    go run "${CONFIG_GENERATOR_DIR}" \
		--gcs-bucket=${PROW_GCS} \
		--generate-testgrid-config=${GENERATE_TESTGRID_CONFIG} \
		--generate-maintenance-jobs=${GENERATE_MAINTENANCE_JOBS} \
		--image-docker=gcr.io/${PROJECT}/test-infra \
		--plugins-config-output="${CONFIG_DIR}/core/plugins.yaml" \
		--prow-config-output="${CONFIG_DIR}/core/config.yaml" \
		--prow-jobs-config-output="${CONFIG_DIR}/jobs/config.yaml" \
		--prow-host=${PROW_HOST} \
		--testgrid-config-output="${CONFIG_DIR}/testgrid/testgrid.yaml" \
		--testgrid-gcs-bucket=${TESTGRID_GCS} \
		"${CONFIG_DIR}/${KNATIVE_CONFIG}"
}


# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh

# Generate Prow configs since we are using generator

# Generate config for production Prow
generate_config

# Generate config for staging Prow
CONFIG_DIR="${REPO_ROOT_DIR}/config/prow-staging"
PROW_GCS="knative-prow-staging"
GENERATE_TESTGRID_CONFIG="false"
GENERATE_MAINTENANCE_JOBS="false"
PROJECT="knative-tests-staging"
PROW_HOST="https://prow-staging.knative.dev"
TESTGRID_GCS="knative-testgrid-staging"
KNATIVE_CONFIG="config_staging.yaml"

generate_config
