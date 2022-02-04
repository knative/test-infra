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

export KUBECONFIG

# https://github.com/istio/test-infra/issues/1636
# This prevents the Kube current-context in the execution environment from being
# overwritten unless the intention is made explicit w/ the `save` parameter.
# e.g.
#               make get-cluster-credentials save=true

.PHONY: save-kubeconfig
save-kubeconfig:
ifndef save
        $(eval KUBECONFIG=$(shell mktemp))
endif

.PHONY: activate-serviceaccount
activate-serviceaccount:
ifdef GOOGLE_APPLICATION_CREDENTIALS
        gcloud auth activate-service-account --key-file="$(GOOGLE_APPLICATION_CREDENTIALS)"
endif

.PHONY: get-cluster-credentials
get-cluster-credentials: save-kubeconfig activate-serviceaccount
        gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"

.PHONY: get-build-cluster-credentials
get-build-cluster-credentials: save-kubeconfig activate-serviceaccount
        gcloud container clusters get-credentials "$(CLUSTER_BUILD)" --project="$(PROJECT_BUILD)" --zone="$(ZONE)"

# These are the usual GKE variables.
PROJECT       ?= knative-tests
PROJECT_BUILD ?= knative-tests
ZONE          ?= us-central1-f
CLUSTER       ?= prow
CLUSTER_BUILD ?= knative-prow-build-cluster
JOB_NAMESPACE ?= test-pods

.PHONY: deploy
deploy: get-cluster-credentials
# Apply the ProwJob CRD with --server-side=true due to its size.
	kubectl apply --server-side=true -f ./cluster/prowjob-crd/
	kubectl apply -f ./cluster/

.PHONY: deploy-build
deploy-build: get-build-cluster-credentials
	kubectl apply -f ./cluster/build/

.PHONY: deploy-monitoring
deploy-monitoring:
	$(MAKE) -C ./cluster/monitoring apply
