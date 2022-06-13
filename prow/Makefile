# Copyright 2022 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include Makefile.gcloud.mk

# These are the usual GKE variables.
PROJECT       ?= knative-tests
PROJECT_BUILD ?= knative-tests
REGION        ?= us-central1
CLUSTER       ?= prow
CLUSTER_BUILD ?= prow-build
JOB_NAMESPACE ?= test-pods

.PHONY: deploy
deploy: get-cluster-credentials
# Apply the ProwJob CRD with --server-side=true due to its size.
	kubectl apply --server-side=true -f ./cluster/control-plane/prowjob-crd/
	kubectl apply -f ./cluster/control-plane

.PHONY: deploy-build
deploy-build: get-build-cluster-credentials
	kubectl apply -f ./cluster/build/

# .PHONY: deploy-monitoring
# deploy-monitoring:
# 	$(MAKE) -C ./cluster/monitoring apply
