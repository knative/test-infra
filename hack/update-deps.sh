#!/usr/bin/env bash

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

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../scripts/library.sh

cd ${REPO_ROOT_DIR}

# Ensure we have everything we need under vendor/
dep ensure

rm -rf $(find vendor/ -name 'OWNERS')
rm -rf $(find vendor/ -name '*_test.go')

# TODO(yt3liu): Remove the mail dependencies from here once it is available through `go get`
# https://github.com/knative/test-infra/issues/841
curl https://storage.googleapis.com/knative-monitoring/cloudmail-v1alpha3-go.zip  \
  --output /tmp/cloudmail-v1alpha3-go.zip
unzip /tmp/cloudmail-v1alpha3-go.zip -d /tmp/cloudmail-v1alpha3-go
mkdir -p vendor/cloud.google.com/go/mail/apiv1alpha3/
mkdir -p vendor/google.golang.org/genproto/googleapis/cloud/mail/v1alpha3
cp -r /tmp/cloudmail-v1alpha3-go/cloud.google.com/go/mail vendor/cloud.google.com/go/
cp -r /tmp/cloudmail-v1alpha3-go/google.golang.org/genproto/googleapis/cloud/mail \
  vendor/google.golang.org/genproto/googleapis/cloud
rm -rf /tmp/cloudmail-v1alpha3-go.zip /tmp/cloudmail-v1alpha3-go
