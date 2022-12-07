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
readonly JOBS_CONFIG_DIR="${REPO_ROOT_DIR}/prow/jobs_config/"

set +u
if [ -z "$1" ]; then
	echo "Pass in the release version number you would like to remove ie: '1.5'"
	exit 1
fi

rm "${JOBS_CONFIG_DIR}"/**/*-"$1".yaml -v
set -u
"${REPO_ROOT_DIR}"/hack/generate-configs.sh
git status
