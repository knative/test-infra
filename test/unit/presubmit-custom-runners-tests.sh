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

source $(dirname $0)/presubmit-integration-tests-common.sh

function build_tests() {
  RAN_BUILD_TESTS=1
  return 0
}

function unit_tests() {
  RAN_UNIT_TESTS=1
  return 0
}

function integration_tests() {
  RAN_INTEGRATION_TESTS=1
  return 0
}

RAN_BUILD_TESTS=0
RAN_UNIT_TESTS=0
RAN_INTEGRATION_TESTS=0

trap check_results EXIT

function check_results() {
  (( RAN_BUILD_TESTS )) || test_failed "Build tests did not run"
  (( RAN_UNIT_TESTS )) || test_failed "Unit tests did not run"
  (( RAN_INTEGRATION_TESTS )) || test_failed "Integration tests did not run"
  echo ">> All tests passed"
}

echo ">> Testing custom test runners"

main $@
