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

source $(dirname $0)/../scripts/library.sh

set -e

trap 'echo "Please rerun \`make -C config/prow config\`."' ERR
header "Checking generated config for production prow and testgrid"
export PROW_CONFIG=$(mktemp)
export PROW_JOB_CONFIG=$(mktemp)
export PROW_PLUGINS=$(mktemp)
export TESTGRID_CONFIG=$(mktemp)
subheader "Regenerating config for production prow and testgrid"
make -C config/prow config
subheader "Comparing the generated config files with the existing config files"
diff --ignore-matching-lines="^# Copyright " "${PROW_CONFIG}" "config/prow/core/config.yaml"
diff --ignore-matching-lines="^# Copyright " "${PROW_JOB_CONFIG}" "config/prow/jobs/config.yaml"
diff --ignore-matching-lines="^# Copyright " "${PROW_PLUGINS}" "config/prow/core/plugins.yaml"
diff --ignore-matching-lines="^# Copyright " "${TESTGRID_CONFIG}" "config/prow/testgrid/testgrid.yaml"

trap 'echo "Please rerun \`make -C config/prow-staging config\`."' ERR
header "Checking generated config for staging prow"
subheader "Regenerating config for staging prow"
make -C config/prow-staging config
subheader "Comparing the generated config files with the existing config files"
diff --ignore-matching-lines="^# Copyright " "${PROW_CONFIG}" "config/prow-staging/core/config.yaml"
diff --ignore-matching-lines="^# Copyright " "${PROW_JOB_CONFIG}" "config/prow-staging/jobs/config.yaml"
diff --ignore-matching-lines="^# Copyright " "${PROW_PLUGINS}" "config/prow-staging/core/plugins.yaml"

trap 'echo "Prow config files have errors, please check."' ERR
header "Validating production Prow config files"
bazel run @k8s//prow/cmd/checkconfig -- \
  --config-path="$(realpath "config/prow/core/config.yaml")" \
  --job-config-path="$(realpath "config/prow/jobs/config.yaml")" \
  --plugin-config="$(realpath "config/prow/core/plugins.yaml")"

header "Validating staging Prow config files"
bazel run @k8s//prow/cmd/checkconfig -- \
  --config-path="$(realpath "config/prow-staging/core/config.yaml")" \
  --job-config-path="$(realpath "config/prow-staging/jobs/config.yaml")" \
  --plugin-config="$(realpath "config/prow-staging/core/plugins.yaml")"

trap 'echo "Testgrid config file has errors, please check."' ERR
header "Validating Testgrid config file"
bazel run @k8s//testgrid/cmd/configurator -- \
  --validate-config-file \
  --yaml="$(realpath "config/prow/testgrid/testgrid.yaml")"
# Ensure TestGrid config can be converted to binary form.
bazel run @k8s//testgrid/cmd/configurator -- \
  --oneshot \
  --output=/dev/null \
  --yaml="$(realpath "config/prow/testgrid/testgrid.yaml")"
