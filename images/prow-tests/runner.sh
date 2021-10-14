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

ORIGINAL_GOPATH="${GOPATH}"

source "${HOME}/.gvm/scripts/gvm"

# By default do not switch the Go version and use the default.
version=""

# If GO_VERSION is defined, use it as the version.
# It has a higher priority than go.mod as it's specified more explicitly.
if [[ -v GO_VERSION ]]; then
  echo "GO_VERSION is defined, overwriting Go version to '${GO_VERSION}'"
  version="${GO_VERSION}"
fi

if [[ -n ${version} ]]; then
  echo "Switching Go version to '${version}'"
  gvm use "${version}"
  # Get our original Go directory back into GOPATH
  pushd "${ORIGINAL_GOPATH}" || exit 2
  gvm pkgset create --local || echo
  gvm pkgset use --local
  popd || exit 2
  # At this point, our GOPATH is set to something like:
  #  GOPATH=/go:/go/.gvm_local/pkgsets/go1.13.10/local:/root/.gvm/pkgsets/go1.13.10/global
  # Which is fine for Go, but some scripts assume GOPATH is a single directory :(
  # Lets hope this doesn't blow up in our face someday
  echo "Overriding GOPATH to '${ORIGINAL_GOPATH}'"
  export GOPATH="${ORIGINAL_GOPATH}"
fi

# Work around networking issues, see https://github.com/kubernetes/test-infra/issues/23741
iptables -t mangle -A POSTROUTING -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu

kubekins-runner.sh "$@"
