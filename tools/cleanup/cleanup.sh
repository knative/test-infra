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

source $(dirname $0)/../../scripts/library.sh

# Global variables
DAYS_TO_KEEP_IMAGES=365 # Keep images up to 1 year by default
RE_PROJECT_NAME="knative-boskos-[a-zA-Z0-9]+"
PROJECT_RESOURCE_YAML=""
PROJECT_TO_CLEANUP=""
DRY_RUN=0

FUNCTION_TO_RUN=""


function parse_args() {
  while [[ $# -ne 0 ]]; do
    local parameter=$1
    case ${parameter} in
      --project-resource-yaml)
        [[ -z $2 || $2 =~ ^-- ]] && abort "expecting value following $1"
        shift
        PROJECT_RESOURCE_YAML=$1
        ;;
      --re-project-name)
        [[ -z $2 || $2 =~ ^-- ]] && abort "expecting value following $1"
        shift
        RE_PROJECT_NAME=$1
        ;;
     --project-to-cleanup)
        [[ -z $2 || $2 =~ ^-- ]] && abort "expecting value following $1"
        shift
        PROJECT_TO_CLEANUP=$1
        ;;
     --days-to-keep)
        [[ -z $2 || $2 =~ ^-- ]] && abort "expecting value following $1"
        shift
        DAYS_TO_KEEP_IMAGES=$1
        ;;
      --dry-run)
        DRY_RUN=1
        ;;
      *) abort "unknown option ${parameter}" ;;
    esac
    shift
  done

  readonly DAYS_TO_KEEP_IMAGES
  readonly PROJECT_RESOURCE_YAML
  readonly RE_PROJECT_NAME
  readonly PROJECT_TO_CLEANUP
  readonly DRY_RUN


  is_int $DAYS_TO_KEEP_IMAGES || abort "days to keep has to be integer"

  (( DRY_RUN )) && echo "-- Running in dry-run mode, no image deletion --"
  echo "Removing images with following rules:"
  case ${FUNCTION_TO_RUN} in
    delete-old-gcr-images)
      echo "- from projects defined in $PROJECT_RESOURCE_YAML, matching $RE_PROJECT_NAME"
      ;;
    delete-old-gcr-images-from-project)
      echo "- from project $PROJECT_TO_CLEANUP"
      ;;
    *) ;;
  esac
  echo "- older than $DAYS_TO_KEEP_IMAGES days"
}

# Delete old images in given GCR project
# Parameters: $1 - gcr to be cleaned up, i.e. gcr.io/fooProj
# Parameters: $2 - days to keep images
function delete_old_gcr_images_from_project() {
  local project_to_cleanup_override=$PROJECT_TO_CLEANUP
  [[ $# -ge 1 ]] && project_to_cleanup_override=$1
  
  [[ -z ${project_to_cleanup_override} ]] && abort "missing project name"
  is_protected_gcr ${project_to_cleanup_override} && \
    abort "\$project_to_cleanup_override set to ${project_to_cleanup_override}, which is forbidden"

  for image in $(gcloud --format='value(name)' container images list --repository=${project_to_cleanup_override}); do
      echo "Checking ${image} for removal"

      delete_old_gcr_images_from_project ${image} ${DAYS_TO_KEEP_IMAGES}

      local target_date=$(date -d "`date`-${DAYS_TO_KEEP_IMAGES}days" +%Y-%m-%d)
      for digest in $(gcloud --format='get(digest)' container images list-tags ${image} \
          --filter="timestamp.datetime<${target_date}" --limit=99999); do
        local full_image="${image}@${digest}"
        if (( $DRY_RUN )); then
          echo "DRYRUN - Deleting image: $full_image"
        else
          echo "Deleting image: $full_image"
          gcloud container images delete -q --force-delete-tags ${full_image}
        fi
      done
  done
}

# Delete old images in GCR projects defined in yaml file provided
# Parameters: $1 - yaml file path defining projects to be cleaned up
# Parameters: $2 - regex pattern for parsing projects' names
# Parameters: $3 - days to keep images
function delete_old_gcr_images() {
  local project_resource_yaml_override=$PROJECT_RESOURCE_YAML
  local re_project_name_override=$RE_PROJECT_NAME
  [[ $# -ge 1 ]] && project_resource_yaml_override=$1
  [[ $# -ge 2 ]] && re_project_name_override=$2

  [[ -z ${project_resource_yaml_override} ]] && abort "missing resource yaml path"
  [[ -z ${re_project_name_override} ]] && abort "missing regex pattern after resource yaml path"

  local target_projects # delared here as local + assignment in one line always return 0 exit code  
  target_projects=$(grep -Eio "${re_project_name_override}" "${project_resource_yaml_override}")
  [[ $? -eq 0 ]] || abort "no project found in ${project_resource_yaml_override}"

  for project in ${target_projects}; do
    echo "Start deleting images from ${project}"
    delete_old_gcr_images_from_project "gcr.io/${project}" $DAYS_TO_KEEP_IMAGES
  done
}


cd ${REPO_ROOT_DIR}

# Entry point
# Dispatching commands to different functions
if [[ $# -ne 0 ]]; then
  case $1 in
    delete-old-gcr-images)
      FUNCTION_TO_RUN=$1
      shift
      parse_args $@
      delete_old_gcr_images
      ;;
    delete-old-gcr-images-from-project)
      FUNCTION_TO_RUN=$1
      shift
      parse_args $@
      delete_old_gcr_images_from_project
      ;;
    *) abort "unknown option ${parameter}" ;;
  esac
else
  warning "no params passed in, no ops. do not source this script unless you know what you are doing"
fi
