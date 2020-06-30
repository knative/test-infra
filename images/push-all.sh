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

# This script pushes all images under current directory

# prow-tests image crashes prow job, don't build it here
EXCLUDED="prow-tests"

set -e

CUR_DIR="$(realpath $(dirname $0))"

for dir in ${CUR_DIR}/*; do
    if [[ -d "${dir}" && -f "${dir}/Makefile" && "${dir}" != "${CUR_DIR}/${EXCLUDED}" ]]; then
        echo "RUNNING: make -C '${dir}' push"
        make -C "${dir}" push
    else
        echo "Skipping '${dir}'"
    fi
done
