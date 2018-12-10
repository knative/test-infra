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

source $(dirname $0)/helper.sh

readonly _CLEANUP_SCRIPT="tools/cleanup/cleanup.sh"
source "$(dirname ${BASH_SOURCE})/../../${_CLEANUP_SCRIPT}"

readonly _SUCCESS=0
readonly _FAILURE=1

readonly _FAKE_NIGHTLY_PROJECT_NAME="gcr.io/knative-nightly"
readonly _FAKE_BOSKOS_PROJECT_NAME="gcr.io/fake-boskos-project"
readonly _PROJECT_RESOURCE_YAML="ci/prow/boskos/resources.yaml"
readonly _RE_PROJECT_NAME="knative-boskos-[a-zA-Z0-9]+"

set -e
cd ${REPO_ROOT_DIR}

echo ">> Testing directly invoking cleanup script"

test_function ${_FAILURE} "error: unknown option" ${_CLEANUP_SCRIPT} "--option-not-exist"
test_function ${_FAILURE} "error: missing project" ${_CLEANUP_SCRIPT} "delete-old-gcr-images-from-project"
test_function ${_FAILURE} "error: missing resource" ${_CLEANUP_SCRIPT} "delete-old-gcr-images"

test_function ${_FAILURE} "error: expecting value following" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --project-resource-yaml --dry-run
test_function ${_FAILURE} "error: expecting value following" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --re-project-name --dry-run
test_function ${_FAILURE} "error: expecting value following" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --project-to-cleanup --dry-run
test_function ${_FAILURE} "error: expecting value following" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --days-to-keep --dry-run

test_function ${_FAILURE} "error: days to keep" ${_CLEANUP_SCRIPT} "delete-old-gcr-images-from-project" --days-to-keep "a" --dry-run
test_function ${_FAILURE} "error: expecting value following" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --days-to-keep --dry-run
test_function ${_FAILURE} "error: days to keep" ${_CLEANUP_SCRIPT} "delete-old-gcr-images" --days-to-keep "a" --dry-run

# Test individual functions
echo ">> Testing deleting images from single project"

test_function ${_FAILURE} "error: missing" delete_old_gcr_images_from_project 
test_function ${_SUCCESS} "" mock_gcloud_function delete_old_gcr_images_from_project ${_FAKE_BOSKOS_PROJECT_NAME}
test_function ${_SUCCESS} "" mock_gcloud_function delete_old_gcr_images_from_project ${_FAKE_BOSKOS_PROJECT_NAME} 1

echo ">> Testing deleting images from multiple projects"

test_function ${_FAILURE} "error: missing" delete_old_gcr_images
test_function ${_FAILURE} "error: no project" delete_old_gcr_images "${_PROJECT_RESOURCE_YAML}file_not_exist" ${_RE_PROJECT_NAME}
test_function ${_FAILURE} "error: no project" delete_old_gcr_images ${_PROJECT_RESOURCE_YAML} "${_RE_PROJECT_NAME}project-not-exist"
test_function ${_SUCCESS} "Start" mock_gcloud_function delete_old_gcr_images ${_PROJECT_RESOURCE_YAML}
test_function ${_SUCCESS} "Start" mock_gcloud_function delete_old_gcr_images ${_PROJECT_RESOURCE_YAML} ${_RE_PROJECT_NAME}

echo ">> All tests passed"
