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

set -o errexit
set -o nounset
set -o pipefail

tmp_dir="$(mktemp -d)"
echo "The scripts will be downloaded at ${tmp_dir}"

curl https://raw.githubusercontent.com/kubernetes/test-infra/master/rbe/install.sh -o "${tmp_dir}/install.sh"
curl https://raw.githubusercontent.com/kubernetes/test-infra/master/rbe/configure.sh -o "${tmp_dir}/configure.sh"

proj=knative-tests
pool=prow-pool
workers=5
disk=200
machine=n2-standard-2
bot=prow-job@knative-tests.iam.gserviceaccount.com

"${tmp_dir}/install.sh" "$proj" "$pool" "$workers" "$disk" "$machine" "$bot"