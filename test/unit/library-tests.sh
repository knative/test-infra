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

# Fake we're in a Prow job, if running locally.
[[ -z "${PROW_JOB_ID:-}" ]] && PROW_JOB_ID=123

source $(dirname $0)/../../scripts/library.sh

set -e

function test_report() {
  local REPORT="$(mktemp)"
  function grepit() {
    if ! grep "$1" ${REPORT} > /dev/null; then
      echo "**** '$1' not found"
      exit 1
    fi
  }
  ARTIFACTS=/tmp
  report_go_test -tags=library -run $1 ./test > ${REPORT} || true
  cat ${REPORT}
  grepit "$2"
  grepit "XML report written"
}

# Cleanup bazel stuff to avoid confusing Prow
function cleanup_bazel() {
  bazel clean
}

trap cleanup_bazel EXIT

echo "*** Tests for report_go_test() ***"

echo ">> Test that test passes"
test_report TestSucceeds "^--- PASS: TestSucceeds"

echo ">> Testing that test fails with fatal"
test_report TestFailsWithFatal "^fatal\s\+TestFailsWithFatal"

echo ">> Testing that test fails with panic"
test_report TestFailsWithPanic "^panic: test timed out"

echo ">> Testing that test fails with SIGQUIT"
test_report TestFailsWithSigQuit "^SIGQUIT: quit"

echo ">> All tests passed"

