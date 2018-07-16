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

# This script runs the presubmit tests; it is started by prow for each PR.
# For convenience, it can also be executed manually.
# Running the script without parameters, or with the --all-tests
# flag, causes all tests to be executed, in the right order.
# Use the flags --build-tests, --unit-tests and --integration-tests
# to run a specific set of tests.

# Extensions or file patterns that don't require presubmit tests
readonly NO_PRESUBMIT_FILES=(\.md \.png ^OWNERS)

[ -f /workspace/library.sh ] && source /workspace/library.sh || eval "$(docker run --entrypoint sh gcr.io/knative-tests/test-infra/prow-tests -c 'cat library.sh')"
[ -v KNATIVE_TEST_INFRA ] || exit 1

# Helper functions.

function build_tests() {
  header "Running build tests"
  make -C prow test
}

function unit_tests() {
  header "TODO(#7): Running unit tests"
}

function integration_tests() {
  header "TODO(#8): Running integration tests"
}

# Script entry point.

# Parse script arguments:
# --all-tests or no arguments: run all tests
# --build-tests: run only the build tests
# --unit-tests: run only the unit tests
# --integration-tests: run only the integration tests
# --emit-metrics: emit metrics when running the E2E tests
RUN_BUILD_TESTS=0
RUN_UNIT_TESTS=0
RUN_INTEGRATION_TESTS=0
EMIT_METRICS=0

all_parameters=$@
[[ -z $1 ]] && all_parameters="--all-tests"

for parameter in ${all_parameters}; do
  case $parameter in
    --all-tests)
      RUN_BUILD_TESTS=1
      RUN_UNIT_TESTS=1
      RUN_INTEGRATION_TESTS=1
      shift
      ;;
    --build-tests)
      RUN_BUILD_TESTS=1
      shift
      ;;
    --unit-tests)
      RUN_UNIT_TESTS=1
      shift
      ;;
    --integration-tests)
      RUN_INTEGRATION_TESTS=1
      shift
      ;;
    --emit-metrics)
      EMIT_METRICS=1
      shift
      ;;
    *)
      echo "error: unknown option ${parameter}"
      exit 1
      ;;
  esac
done

readonly RUN_BUILD_TESTS
readonly RUN_UNIT_TESTS
readonly RUN_INTEGRATION_TESTS
readonly EMIT_METRICS

cd ${REPO_ROOT_DIR}

# Skip presubmit tests if only markdown files were changed.
if [[ -n "${PULL_PULL_SHA}" ]]; then
  # On a presubmit job
  changes="$(git diff --name-only ${PULL_PULL_SHA} ${PULL_BASE_SHA})"
  no_presubmit_pattern="${NO_PRESUBMIT_FILES[*]}"
  no_presubmit_pattern="\(${no_presubmit_pattern// /\\|}\)$"
  echo -e "Changed files in commit ${PULL_PULL_SHA}:\n${changes}"
  if [[ -z "$(echo "${changes}" | grep -v ${no_presubmit_pattern})" ]]; then
    # Nothing changed other than files that don't require presubmit tests
    header "Commit only contains changes that don't affect tests, skipping"
    exit 0
  fi
fi

# Tests to be performed, in the right order if --all-tests is passed.

result=0
if (( RUN_BUILD_TESTS )); then
  build_tests || result=1
fi
if (( RUN_UNIT_TESTS )); then
  unit_tests || result=1
fi
if (( RUN_INTEGRATION_TESTS )); then
  integration_tests || result=1
fi
exit ${result}
