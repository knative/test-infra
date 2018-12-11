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

source $(dirname $0)/helper.sh
source $(dirname $0)/../../scripts/release.sh

set -e

function mock_branch_release() {
  set -e
  BRANCH_RELEASE=1
  TAG=sometag
  function git() {
	echo $@
  }
  function hub() {
	echo $@
  }
  branch_release "$@" 2>&1
}

function call_function_pre() {
  set -e
  local init="$1"
  shift
  eval ${init}
  "$@" 2>&1
}

function call_function_post() {
  set -e
  local post="$1"
  shift
  "$@" 2>&1
  eval ${post}
}

echo ">> Testing helper functions"

test_function 0 "0.2" master_version "v0.2.1"
test_function 0 "0.2" master_version "0.2.1"
test_function 0 "1" release_build_number "v0.2.1"
test_function 0 "1" release_build_number "0.2.1"

echo ">> Testing initialization"

test_function 1 "error: missing version" initialize --version
test_function 1 "error: version format" initialize --version a
test_function 1 "error: version format" initialize --version 0.0
test_function 0 "" parse_flags --version 1.0.0

test_function 1 "error: missing branch" initialize --branch
test_function 1 "error: branch name must be" initialize --branch a
test_function 1 "error: branch name must be" initialize --branch 0.0
test_function 0 "" parse_flags --branch release-0.0

test_function 1 "error: missing release notes" initialize --release-notes
test_function 1 "error: file a doesn't" initialize --release-notes a
test_function 0 "" parse_flags --release-notes $(mktemp)

test_function 1 "error: missing GCS" initialize --release-gcs
test_function 0 "" parse_flags --release-gcs a --publish

test_function 1 "error: missing GCR" initialize --release-gcr
test_function 0 "" parse_flags --release-gcr a --publish

echo ">> Testing GCR/GCS values"

test_function 0 "GCR flag is ignored" initialize --release-gcr foo
test_function 0 "GCS flag is ignored" initialize --release-gcs foo

test_function 0 "Destination GCR: ko.local" initialize
test_function 0 "::" call_function_post "echo :\$RELEASE_GCS:" initialize

test_function 0 "Destination GCR: gcr.io/knative-nightly" initialize --publish
test_function 0 "published to 'knative-nightly/test-infra'" initialize --publish

test_function 0 "Destination GCR: foo" initialize --release-gcr foo --publish
test_function 0 "published to 'foo'" initialize --release-gcs foo --publish

echo ">> Testing release branching"

test_function 0 "" branch_release
test_function 129 "usage: git tag" call_function_pre BRANCH_RELEASE=1 branch_release
test_function 1 "No such file" call_function_pre BRANCH_RELEASE=1 branch_release "K Foo" "a.yaml b.yaml"
test_function 0 "release create" mock_branch_release "K Foo" "$(mktemp) $(mktemp)"

echo ">> Testing validation tests"

test_function 0 "Running release validation" run_validation_tests true
test_function 0 "" call_function_pre SKIP_TESTS=1 run_validation_tests true
test_function 0 "i_passed" run_validation_tests "echo i_passed"
test_function 1 "validation tests failed" run_validation_tests false

echo ">> All tests passed"
