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

# Requires gcloud and kubectl. Assume the cluster contexts are already set for
# both service and build clusters

export KUBECONFIG_SERVICE_CLUSTER="gke_knative-tests-staging_us-central1-f_prow"
export KUBECONFIG_BUILD_CLUSTER="gke_knative-tests-staging_us-central1-f_knative-prow-build-cluster"

../../prod/build-cluster/add-secrets.sh
