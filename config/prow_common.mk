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

SHELL := /bin/bash
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(SELF_DIR)../common.mk

# This file is used by prod and staging Makefiles

# Default settings for the CI/CD system.

CLUSTER       	?= prow
BUILD_CLUSTER 	?= knative-prow-build-cluster
ZONE          	?= us-central1-f
JOB_NAMESPACE 	?= test-pods

SKIP_CONFIG_BACKUP        ?=

# Any changes to file location must be made to staging directory also
# or overridden in the Makefile before this file is included.
PROW_PLUGINS     				?= prod/core/plugins.yaml
PROW_CONFIG      				?= prod/core/config.yaml
PROW_JOB_CONFIG  				?= prod/jobs

PROW_DEPLOYS     				?= prod/cluster
BUILD_CLUSTER_PROW_DEPLOYS    	?= build-cluster/cluster
PROW_GCS         				?= knative-prow
PROW_CONFIG_GCS  				?= gs://$(PROW_GCS)/configs

BOSKOS_RESOURCES 				?= build-cluster/boskos/boskos_resources.yaml

# Useful shortcuts.

SET_CONTEXT   				  := gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"
SET_BUILD_CLUSTER_CONTEXT     := gcloud container clusters get-credentials "$(BUILD_CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"
UNSET_CONTEXT 				  := kubectl config unset current-context

.PHONY: help get-cluster-credentials unset-cluster-credentials
help:
	@echo "Help"
	@echo "'Update' means updating the servers and can only be run by oncall staff."
	@echo "Common usage:"
	@echo " make update-prow-cluster: Update all Prow things on the server to match the current branch. Errors if not master."
	@echo " make update-testgrid-config: Update the Testgrid config"
	@echo " make get-cluster-credentials: Setup kubectl to point to Prow cluster"
	@echo " make unset-cluster-credentials: Clear kubectl context"

# Useful general targets.
get-cluster-credentials:
	$(SET_CONTEXT)

unset-cluster-credentials:
	$(UNSET_CONTEXT)

get-build-cluster-credentials:
	$(SET_BUILD_CLUSTER_CONTEXT)

.PHONY: update-prow-config update-all-boskos-deployments update-boskos-resource update-almost-all-cluster-deployments update-single-cluster-deployment test update-testgrid-config confirm-master

# Update prow config
update-prow-config: confirm-master
	$(SET_CONTEXT)
	python3 <(curl -sSfL https://raw.githubusercontent.com/istio/test-infra/master/prow/recreate_prow_configmaps.py) \
		--prow-config-path=$(realpath $(PROW_CONFIG)) \
		--plugins-config-path=$(realpath $(PROW_PLUGINS)) \
		--job-config-dir=$(realpath $(PROW_JOB_CONFIG)) \
		--wet \
		--silent
	$(UNSET_CONTEXT)

# Update all deployments of boskos
# Boskos is separate because of patching done in staging Makefile
# Double-colon because staging Makefile piggy-backs on this
update-all-boskos-deployments:: confirm-master
	$(SET_BUILD_CLUSTER_CONTEXT)
	@for f in $(wildcard $(BUILD_CLUSTER_PROW_DEPLOYS)/*boskos*.yaml); do kubectl apply -f $${f}; done
	$(UNSET_CONTEXT)

# Update the list of resources for Boskos
update-boskos-resource: confirm-master
	$(SET_BUILD_CLUSTER_CONTEXT)
	kubectl create configmap resources --from-file=config=$(BOSKOS_RESOURCES) --dry-run --save-config -o yaml | kubectl --namespace="$(JOB_NAMESPACE)" apply -f -
	$(UNSET_CONTEXT)

# Update all deployments of cluster except Boskos
# Boskos is separate because of patching done in staging Makefile
# Double-colon because staging Makefile piggy-backs on this
update-almost-all-cluster-deployments:: confirm-master
	$(SET_CONTEXT)
	@for f in $(filter-out $(wildcard $(PROW_DEPLOYS)/*boskos*.yaml),$(wildcard $(PROW_DEPLOYS)/*.yaml)); do kubectl apply -f $${f}; done
	$(UNSET_CONTEXT)
	$(SET_BUILD_CLUSTER_CONTEXT)
	@for f in $(filter-out $(wildcard $(BUILD_CLUSTER_PROW_DEPLOYS)/*boskos*.yaml),$(wildcard $(BUILD_CLUSTER_PROW_DEPLOYS)/*.yaml)); do kubectl apply -f $${f}; done
	$(UNSET_CONTEXT)

# Update single deployment of cluster, expect passing in ${NAME} like `make update-single-cluster-deployment NAME=crier_deployment`
update-single-cluster-deployment: confirm-master
	$(SET_CONTEXT)
	kubectl apply -f $(PROW_DEPLOYS)/$(NAME).yaml
	$(UNSET_CONTEXT)
	$(SET_BUILD_CLUSTER_CONTEXT)
	kubectl apply -f $(BUILD_CLUSTER_PROW_DEPLOYS)/$(NAME).yaml
	$(UNSET_CONTEXT)

# Update all resources on Prow cluster
update-prow-cluster: update-almost-all-cluster-deployments update-all-boskos-deployments update-boskos-resource update-prow-config

# Update TestGrid config.
# Application Default Credentials must be set, otherwise the upload will fail.
# Either export $GOOGLE_APPLICATION_CREDENTIALS pointing to a valid service
# account key, or temporarily use your own credentials by running
# gcloud auth application-default login
update-testgrid-config: confirm-master
	bazel run @k8s//testgrid/cmd/configurator -- \
		--oneshot \
		--output=gs://$(TESTGRID_GCS)/config \
		--yaml=$(realpath $(TESTGRID_CONFIG))

