#!/usr/bin/env bash

# Copyright 2021 The Knative Authors
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

# This script updates Prow job configs.

if [[ "${1}" != "check" || "${1}" != "update" ]]; then
	echo "Usage is 'config.sh [COMMAND]' where command is 'check' or 'update'."
fi

source $(dirname "${BASH_SOURCE[0]}")/../../../vendor/knative.dev/hack/library.sh

set -e

make -C "${REPO_ROOT_DIR}/config" get-cluster-credentials
add_trap "make -C ${REPO_ROOT_DIR}/config unset-cluster-credentials" EXIT

if [[ "${1}" != "update" ]]; then
	orig_kubeconfig=${KUBECONFIG}
	KUBECONFIG="/tmp/kubeconfig.yaml"
	add_trap "rm ${KUBECONFIG}" EXIT
	add_trap "KUBECONFIG=${orig_kubeconfig}" EXIT
else #check
	# Settings for inrepoconfig checking, if only `REPO_NAME` is supplied then:
	# `REPO_YAML_PATH` is default to `("/home/prow/go/src/github.com/%s/.prow.yaml",
	# o.prowYAMLRepoName)`. See
	# https://github.com/kubernetes/test-infra/blob/a3f67762b39ce6eeb221b4d475cf4a1a32ac0c54/prow/cmd/checkconfig/main.go#L144
	REPO_NAME_TO_CHECK="${2:-}"
	REPO_YAML_PATH_TO_CHECK="${3:-}"
fi

# Download prow core config from prow
if (( IS_OSX )); then
  # On OS X, the file has to be under /private other it cannot be mounted by the container.
  # See https://stackoverflow.com/questions/45122459/docker-mounts-denied-the-paths-are-not-shared-from-os-x-and-are-not-known/45123074
  CONFIG_YAML="/private$(mktemp)"
  PLUGINS_YAML="/private$(mktemp)"
else
  CONFIG_YAML="$(mktemp)"
  PLUGINS_YAML="$(mktemp)"
fi

JOB_YAML="${REPO_ROOT_DIR}/prow/jobs"

kubectl get configmaps config -o "jsonpath={.data['config\.yaml']}" >"${CONFIG_YAML}"
echo "Prow core config downloaded at ${CONFIG_YAML}"
kubectl get configmaps plugins -o "jsonpath={.data['plugins\.yaml']}" > "${PLUGINS_YAML}"
echo "Prow plugins downloaded at ${PLUGINS_YAML}"

if [[ "${1}" != "update" ]]; then
	docker run -i --rm \
	    -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${PLUGINS_YAML}:${PLUGINS_YAML}" -v "${JOB_YAML}:${JOB_YAML}" \
	    -v "${KUBECONFIG}:${KUBECONFIG}:ro" \
	    -w "${PWD}" \
	    gcr.io/k8s-prow/config-bootstrapper:v20220124-9887456efc \
	    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_YAML}" \
	    "--plugin-config=${PLUGINS_YAML}" \
	    "--source-path=${REPO_ROOT_DIR}" \
	    "--kubeconfig=${KUBECONFIG}" \
	    "--dry-run=false"
else #check
	REPO_NAME_ARG=""
	REPO_YAML_PATH_ARG=""
	[[ -n "${REPO_NAME_TO_CHECK}" ]] && REPO_NAME_ARG="--prow-yaml-repo-name=${REPO_NAME_TO_CHECK}"
	[[ -n "${REPO_YAML_PATH_TO_CHECK}" ]] && REPO_YAML_PATH_ARG="--prow-yaml-repo-path=${REPO_YAML_PATH_TO_CHECK}"

	# TODO: Re-enable the mismatched-tide-lenient warning after we're done experimenting with CODEOWNERS in net-contour.
	docker run -i --rm \
	    -v "${PWD}:${PWD}" -v "${CONFIG_YAML}:${CONFIG_YAML}" -v "${PLUGINS_YAML}:${PLUGINS_YAML}" -v "${JOB_YAML}:${JOB_YAML}" \
	    -w "${PWD}" \
	    gcr.io/k8s-prow/checkconfig:v20220124-9887456efc \
	    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_YAML}" \
	    "--plugin-config=${PLUGINS_YAML}" "--strict" "--exclude-warning=mismatched-tide" \
	    "--exclude-warning=long-job-names" \
	    "--exclude-warning=mismatched-tide-lenient" \
	    "${REPO_NAME_ARG}" "${REPO_YAML_PATH_ARG}"
fi
