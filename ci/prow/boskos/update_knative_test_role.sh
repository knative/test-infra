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

set -e

readonly CUSTOM_ROLE_NAME="KnativeTests"

readonly LAST=${1:?"First argument is the last number of the project(s) to update."}
readonly OPERATION=${2:?"Second argument is the operation on this role, can be --add-permissions or --remove-permissions"}
readonly PERMISSIONS=${3:?"Third argument is a list of permissions separated with comma"}

for i in $(seq -f "%02g" 1 ${LAST})
do
  PROJECT="knative-boskos-$i"
  gcloud iam roles update ${CUSTOM_ROLE_NAME} --project ${PROJECT} ${OPERATION} ${PERMISSIONS}
done
