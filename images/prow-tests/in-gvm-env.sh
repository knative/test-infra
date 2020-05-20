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

source "${HOME}/.gvm/scripts/gvm"

set -e
set -x

if [[ ! -v GO_VERSION ]]; then
  echo "--- FAIL: need GO_VERSION defined for in-gvm-env.sh"
  exit 1
fi
gvm use "${GO_VERSION}"
# Get our original Go directory back into GOPATH
pushd /go
gvm pkgset create --local || echo
gvm pkgset use --local
popd
# At this point, our GOPATH is set to something like:
#  GOPATH=/go:/go/.gvm_local/pkgsets/go1.13.10/local:/root/.gvm/pkgsets/go1.13.10/global
echo "GOPATH is ${GOPATH}"
# We want to install tools to the global bin directory for that version of Go,
#  so set GOBIN
IFS=: read -a arr <<< "${GOPATH}"
export GOBIN="${arr[2]}/bin"

"$@"
