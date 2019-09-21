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

readonly CUSTOM_ROLE_FILE="custom_role.yaml"
readonly CUSTOM_ROLE_NAME="KnativeIntegrationTestsRunner"

project_count=$(grep -n "knative-boskos-" resources.yaml | wc -l)
for i in $(seq -f "%02g" 1 ${project_count})
do
  PROJECT="knative-boskos-$i"
  while true; do
    # Create the role.
    O=$(gcloud iam roles create ${CUSTOM_ROLE_NAME} -q --project ${PROJECT} --file ${CUSTOM_ROLE_FILE} 2>&1)
    E=$(echo $O | grep "ERROR: (gcloud.iam.roles.create)" | grep "already exists.")
    echo $O
    # If role already exists, update it.
    if [ ! -z "$E" ]; then
      gcloud iam roles describe ${CUSTOM_ROLE_NAME} --project ${PROJECT} | grep "^etag: " >> ${CUSTOM_ROLE_FILE}
      O=$(gcloud iam roles update ${CUSTOM_ROLE_NAME} -q --project ${PROJECT} --file ${CUSTOM_ROLE_FILE} 2>&1)
    fi
    # If permission is invalid, remove it.
    E=$(echo $O | grep -o "INVALID_ARGUMENT: Permission [a-zA-Z.]\+ is " | cut -f3 -d' ')
    if [ -z "$E" ]; then
      echo $O
      break
    fi
    echo "Removing $E"
    grep -v "^- ${E}$" ${CUSTOM_ROLE_FILE} > ${CUSTOM_ROLE_FILE}1
    mv ${CUSTOM_ROLE_FILE}1 ${CUSTOM_ROLE_FILE}
  done
done
