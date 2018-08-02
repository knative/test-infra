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

# Load github.com/knative/test-infra/images/prow-tests/scripts/presubmit-tests.sh
[ -f /workspace/presubmit-tests.sh ] \
  && source /workspace/presubmit-tests.sh \
  || eval "$(docker run --entrypoint sh gcr.io/knative-tests/test-infra/prow-tests -c 'cat presubmit-tests.sh')"
[ -v KNATIVE_TEST_INFRA ] || exit 1

function build_tests() {
  header "Running build tests"
  make -C ci/prow test
  make -C ci/testgrid test
}

function unit_tests() {
  header "TODO(#7): Running unit tests"
}

function integration_tests() {
  header "TODO(#8): Running integration tests"
}

main $@
