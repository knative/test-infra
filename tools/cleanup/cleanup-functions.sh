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

# Functions for cleaning up GCRs.
# It doesn't do anything when called from command line.

source $(dirname $0)/../../scripts/library.sh

# Delete old images in the given GCR.
# Parameters: $1 - gcr to be cleaned up (e.g. gcr.io/fooProj)
#             $2 - days to keep images
function delete_old_images_from_gcr() {
  [[ -z $1 ]] && abort "missing gcr name"
  [[ -z $2 ]] && abort "missing days to keep images"

  is_protected_gcr $1 && \
    abort "Target GCR set to $1, which is forbidden"

  for image in $(gcloud --format='value(name)' container images list --repository=$1); do
      echo "Checking ${image} for removal"

      delete_old_images_from_gcr ${image} $2

      local target_date=$(date -d "`date`-$2days" +%Y-%m-%d)
      for digest in $(gcloud --format='get(digest)' container images list-tags ${image} \
          --filter="timestamp.datetime<${target_date}" --limit=99999); do
        local full_image="${image}@${digest}"
        echo "Deleting image: ${full_image}"
        if (( DRY_RUN )); then
          echo "[DRY RUN] gcloud container images delete -q --force-delete-tags ${full_image}"
        else
          gcloud container images delete -q --force-delete-tags ${full_image}
        fi
      done
  done
}

# Delete old images in the GCP projects defined in the yaml file provided.
# Parameters: $1 - yaml file path defining projects that will be cleaned up
#             $2 - regex pattern for parsing the project names
#             $3 - days to keep images
function delete_old_gcr_images() {
  [[ -z $1 ]] && abort "missing resource yaml path"
  [[ -z $2 ]] && abort "missing regex pattern for project name"
  [[ -z $3 ]] && abort "missing days to keep images"

  local target_projects # delared here as local + assignment in one line always return 0 exit code
  target_projects="$(grep -Eio "$2" "$1")"
  [[ $? -eq 0 ]] || abort "no project found in $1"

  for project in ${target_projects}; do
    echo "Start deleting images from ${project}"
    delete_old_images_from_gcr "gcr.io/${project}" $3
  done
}
