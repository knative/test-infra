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

if [[ "${1}" != "check" && "${1}" != "update" ]]; then
	echo "Usage is 'config.sh [COMMAND]' where command is 'check' or 'update'."
	exit 1
fi

source $(dirname "${BASH_SOURCE[0]}")/../vendor/knative.dev/hack/library.sh

set -ex

orig_kubeconfig=${KUBECONFIG}
export KUBECONFIG
KUBECONFIG="$(mktemp -d)/kubeconfig.yaml"
touch "${KUBECONFIG}"
# add_trap "rm ${KUBECONFIG}" EXIT
add_trap "export KUBECONFIG=${orig_kubeconfig}" EXIT

make -C "${REPO_ROOT_DIR}/config" get-cluster-credentials
add_trap "make -C ${REPO_ROOT_DIR}/config unset-cluster-credentials" EXIT

if [[ "${1}" == "check" ]]; then
	# Settings for inrepoconfig checking, if only `REPO_NAME` is supplied then:
	# `REPO_YAML_PATH` is default to `("/home/prow/go/src/github.com/%s/.prow.yaml",
	# o.prowYAMLRepoName)`. See
	# https://github.com/kubernetes/test-infra/blob/a3f67762b39ce6eeb221b4d475cf4a1a32ac0c54/prow/cmd/checkconfig/main.go#L144
	REPO_NAME_TO_CHECK="${2:-}"
	REPO_YAML_PATH_TO_CHECK="${3:-}"
fi

JOB_YAML="${REPO_ROOT_DIR}/prow/jobs"
CONFIG_YAML="${REPO_ROOT_DIR}/prow/config.yaml"
PLUGINS_YAML="${REPO_ROOT_DIR}/prow/plugins.yaml"

if [[ "${1}" == "update" ]]; then
	if [[ ! -x "$(command -v config-bootstrapper)" ]]; then
		echo "--- FAIL: config-bootstrapper not installed, please install it from https://github.com/kubernetes/test-infra/tree/master/prow/cmd/config-bootstrapper"; exit 1;
	fi
	# FIXME: this command is supposed to work but seems not, need to figure out why...
	config-bootstrapper \
	    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_YAML}" \
	    "--plugin-config=${PLUGINS_YAML}" \
	    "--source-path=." \
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
	    gcr.io/k8s-prow/checkconfig:v20220215-03bad893f1 \
	    "--config-path=${CONFIG_YAML}" "--job-config-path=${JOB_YAML}" \
	    "--plugin-config=${PLUGINS_YAML}" "--strict" "--exclude-warning=mismatched-tide" \
	    "--exclude-warning=long-job-names" \
	    "--exclude-warning=mismatched-tide-lenient" \
	    "${REPO_NAME_ARG}" "${REPO_YAML_PATH_ARG}"
fi
