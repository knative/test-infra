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

source $(dirname $0)/test-helper.sh
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

echo ">> Testing helper functions"

test_function ${SUCCESS} "0.2" master_version "v0.2.1"
test_function ${SUCCESS} "0.2" master_version "0.2.1"
test_function ${SUCCESS} "1" release_build_number "v0.2.1"
test_function ${SUCCESS} "1" release_build_number "0.2.1"

echo ">> Testing initialization"

test_function ${FAILURE} "error: missing parameter" initialize --version
test_function ${FAILURE} "error: version format" initialize --version a
test_function ${FAILURE} "error: version format" initialize --version 0.0
test_function ${SUCCESS} "" parse_flags --version 1.0.0

test_function ${FAILURE} "error: missing parameter" initialize --branch
test_function ${FAILURE} "error: branch name must be" initialize --branch a
test_function ${FAILURE} "error: branch name must be" initialize --branch 0.0
test_function ${SUCCESS} "" parse_flags --branch release-0.0

test_function ${FAILURE} "error: missing parameter" initialize --release-notes
test_function ${FAILURE} "error: file a doesn't" initialize --release-notes a
test_function ${SUCCESS} "" parse_flags --release-notes $(mktemp)

test_function ${FAILURE} "error: missing parameter" initialize --release-gcs
test_function ${SUCCESS} "" parse_flags --release-gcs a --publish

test_function ${FAILURE} "error: missing parameter" initialize --release-gcr
test_function ${SUCCESS} "" parse_flags --release-gcr a --publish

token_file=$(mktemp)
echo -e "abc " > ${token_file}
test_function ${SUCCESS} ":abc:" call_function_post "echo :\$GITHUB_TOKEN:" initialize --github-token ${token_file}

echo ">> Testing GCR/GCS values"

test_function ${SUCCESS} "GCR flag is ignored" initialize --release-gcr foo
test_function ${SUCCESS} "GCS flag is ignored" initialize --release-gcs foo

test_function ${SUCCESS} "Destination GCR: ko.local" initialize
test_function ${SUCCESS} "::" call_function_post "echo :\$RELEASE_GCS:" initialize

test_function ${SUCCESS} "Destination GCR: gcr.io/knative-nightly" initialize --publish
test_function ${SUCCESS} "published to 'knative-nightly/test-infra'" initialize --publish

test_function ${SUCCESS} "Destination GCR: foo" initialize --release-gcr foo --publish
test_function ${SUCCESS} "published to 'foo'" initialize --release-gcs foo --publish

echo ">> Testing release branching"

test_function ${SUCCESS} "" branch_release
test_function 129 "usage: git tag" call_function_pre BRANCH_RELEASE=1 branch_release
test_function ${FAILURE} "No such file" call_function_pre BRANCH_RELEASE=1 branch_release "K Foo" "a.yaml b.yaml"
test_function ${SUCCESS} "release create" mock_branch_release "K Foo" "$(mktemp) $(mktemp)"

echo ">> Testing validation tests"

test_function ${SUCCESS} "Running release validation" run_validation_tests true
test_function ${SUCCESS} "" call_function_pre SKIP_TESTS=1 run_validation_tests true
test_function ${SUCCESS} "i_passed" run_validation_tests "echo i_passed"
test_function ${FAILURE} "validation tests failed" run_validation_tests false

echo ">> All tests passed"
