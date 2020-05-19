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

set -e

# This script updates test-infra scripts in-repo.
# Run it to update (usually from hack/update-deps.sh) the current scripts.
# Scripts are installed to REPO_ROOT/scripts/test-infra

# The following arguments are accepted:
# TODO: --verify
#  Verify the contents of scripts/test-infra match the contents from commit sha in scripts/test-infra/sha
# --branch X
#  Defines which branch of test-infra to get scripts from; defaults to master
# --first-time
#  Run this script from your repo root directory to install scripts for the first time
#  TODO: also sed -i the scripts in the current repo to point to new file

declare -i FIRST_TIME_SETUP=0
declare SCRIPTS_BRANCH=master

while [[ $# -ne 0 ]]; do
  parameter="$1"
  case ${parameter} in
    --branch)
      shift
      SCRIPTS_BRANCH="$1"
      ;;
    --first-time)
      FIRST_TIME_SETUP=1
      ;;
    *) abort "unknown option ${parameter}" ;;
  esac
  shift
done

function do_read_tree() {
    mkdir -p scripts/test-infra
    git read-tree --prefix=scripts/test-infra -u "test-infra/${SCRIPTS_BRANCH}:scripts"
    git show-ref -s -- "refs/remotes/test-infra/${SCRIPTS_BRANCH}" > scripts/test-infra/sha
    git add scripts/test-infra/sha
    echo "test-infra scripts installed to scripts/test-infra from branch ${SCRIPTS_BRANCH}"
}

function run() {
  if (( FIRST_TIME_SETUP )); then
    if [[ ! -d .git ]]; then
      echo "I don't believe you are in a repo root; exiting"
      exit 5
    fi
    git remote add test-infra https://github.com/knative/test-infra.git || true
    git fetch test-infra "${SCRIPTS_BRANCH}"
    do_read_tree
  else
    pushd "$(dirname "${BASH_SOURCE[0]}")/../.."
    trap popd EXIT

    git remote add test-infra https://github.com/knative/test-infra.git || true
    git fetch test-infra "${SCRIPTS_BRANCH}"
    git rm -fr scripts/test-infra
    rm -fR scripts/test-infra
    do_read_tree
  fi
}

run
