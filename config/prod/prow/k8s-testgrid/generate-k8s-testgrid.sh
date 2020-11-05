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

# This script fetches config.yaml from cluster.

# config.yaml can be viewed here:
# https://github.com/GoogleCloudPlatform/oss-test-infra/blob/master/prow/prowjobs/GoogleCloudPlatform/oss-test-infra/gcp-oss-test-infra-config.yaml

# Base this script loosely on:
# https://github.com/knative/test-infra/blob/master/config/prod/prow/check_config.sh


CONFIG_YAML="$(mktemp)" # randomize the name of the local config file we will write

kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"

# Docker's -i/--interactive allows you to send commands to the container via standard input ("STDIN"), which means you can "interactively" type commands to the pseudo-tty/terminal created by the -t switch.
# if you’d like Docker to automatically clean up the container and remove the file system when the container exits, you can add the --rm flag
docker run -i --rm \
    -v "${PWD}:${PWD}" \
    -v "${CONFIG_YAML}:${CONFIG_YAML}" \
    -w "${PWD}" \
    gcr.io/k8s-prow/transfigure:v20201110-9e512b5af0 \
    "/etc/github-token/oauth" \
    "${CONFIG_YAML}" \
    "config/prod/prow/jobs/config.yaml" \
    "config/prod/prow/k8s-testgrid/k8s-testgrid.yaml" \
    "knative" 