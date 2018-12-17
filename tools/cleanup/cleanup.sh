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

# This is a script to clean up stale resources

source $(dirname $0)/cleanup-functions.sh

# Global variables
DAYS_TO_KEEP_IMAGES=365 # Keep images up to 1 year by default
RE_PROJECT_NAME="knative-boskos-[a-zA-Z0-9]+"
PROJECT_RESOURCE_YAML=""
GCR_TO_CLEANUP=""
ARTIFACTS_DIR=""
DRY_RUN=0


function parse_args() {
  while [[ $# -ne 0 ]]; do
    local parameter=$1
    case ${parameter} in
      --dry-run) DRY_RUN=1 ;;
      *)
        [[ -z $2 || $2 =~ ^-- ]] && abort "expecting value following $1"
        shift
        case ${parameter} in
          --project-resource-yaml) PROJECT_RESOURCE_YAML=$1 ;;
          --re-project-name) RE_PROJECT_NAME=$1 ;;
          --gcr-to-cleanup) GCR_TO_CLEANUP=$1 ;;
          --days-to-keep) DAYS_TO_KEEP_IMAGES=$1 ;;
          --artifacts) ARTIFACTS_DIR=$1 ;;
          --service-account)
            gcloud auth activate-service-account --key-file=$1 || exit 1
            ;;
          *) abort "unknown option ${parameter}" ;;
        esac
    esac
    shift
  done

  is_int ${DAYS_TO_KEEP_IMAGES} || abort "days to keep has to be integer"

  readonly DAYS_TO_KEEP_IMAGES
  readonly PROJECT_RESOURCE_YAML
  readonly RE_PROJECT_NAME
  readonly GCR_TO_CLEANUP
  readonly ARTIFACTS_DIR
  readonly DRY_RUN
}

# Script entry point

cd ${REPO_ROOT_DIR}

if [[ -z $1 ]]; then
  abort "missing parameters to the tool"
fi

FUNCTION_TO_RUN=$1
shift
parse_args $@

(( DRY_RUN )) && echo "-- Running in dry-run mode, no image deletion --"

echo "Removing images with following rules:"
echo "- older than ${DAYS_TO_KEEP_IMAGES} days"
case ${FUNCTION_TO_RUN} in
  delete-old-gcr-images)
    echo "- from projects defined in '${PROJECT_RESOURCE_YAML}', matching '${RE_PROJECT_NAME}"
    delete_old_gcr_images "${PROJECT_RESOURCE_YAML}" "${RE_PROJECT_NAME}" "${DAYS_TO_KEEP_IMAGES}"
    ;;
  delete-old-images-from-gcr)
    echo "- from gcr '${GCR_TO_CLEANUP}'"
    delete_old_images_from_gcr "${GCR_TO_CLEANUP}" "${DAYS_TO_KEEP_IMAGES}"
    ;;
  *) abort "unknown option '${FUNCTION_TO_RUN}'" ;;
esac

# Gubernator considers job failure if "junit_*.xml" not found under artifact,
#   create a placeholder file to make this job succeed
if [[ ! -z ${ARTIFACTS_DIR} ]]; then
  echo "<testsuite time='0'/>" > "${ARTIFACTS_DIR}/junit_knative.xml"
fi
