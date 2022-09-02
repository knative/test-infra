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

# Default settings for the CI/CD system.

PROJECT       ?= knative-tests
PROW_GCS      ?= knative-prow
PROW_HOST     ?= https://prow.knative.dev
TESTGRID_GCS  ?= knative-testgrid

CLUSTER           ?= prow
BUILD_CLUSTER     ?= knative-prow-build-cluster
REGION            ?= us-central1
JOB_NAMESPACE     ?= test-pods

SKIP_CONFIG_BACKUP        ?=

# Any changes to file location must be made to staging directory also
# or overridden in the Makefile before this file is included.
PROW_PLUGINS                    ?= prow/core/plugins.yaml
PROW_CONFIG                     ?= prow/core/config.yaml
PROW_JOB_CONFIG                 ?= prow/jobs

PROW_GCS                        ?= knative-prow
PROW_CONFIG_GCS                 ?= gs://$(PROW_GCS)/configs

# Useful shortcuts.

SET_CONTEXT                     := gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --region="$(REGION)"
SET_BUILD_CLUSTER_CONTEXT       := gcloud container clusters get-credentials "$(BUILD_CLUSTER)" --project="$(PROJECT)" --region="$(REGION)"
UNSET_CONTEXT                   := kubectl config unset current-context

.PHONY: help activate-serviceaccount get-cluster-credentials unset-cluster-credentials
help:
	@echo "Help"
	@echo "'Update' means updating the servers and can only be run by oncall staff."
	@echo "Common usage:"
	@echo " make update-testgrid-config: Update the Testgrid config"
	@echo " make get-cluster-credentials: Setup kubectl to point to Prow cluster"
	@echo " make unset-cluster-credentials: Clear kubectl context"

# Useful general targets.
activate-serviceaccount:
ifdef GOOGLE_APPLICATION_CREDENTIALS
	gcloud auth activate-service-account --key-file="$(GOOGLE_APPLICATION_CREDENTIALS)"
endif

get-cluster-credentials: activate-serviceaccount
	$(SET_CONTEXT)

unset-cluster-credentials:
	$(UNSET_CONTEXT)

get-build-cluster-credentials: activate-serviceaccount
	$(SET_BUILD_CLUSTER_CONTEXT)

.PHONY: update-testgrid-config confirm-main

# Update TestGrid config.
# Application Default Credentials must be set, otherwise the upload will fail.
# Either export $GOOGLE_APPLICATION_CREDENTIALS pointing to a valid service
# account key, or temporarily use your own credentials by running
# gcloud auth application-default login
update-testgrid-config: confirm-main
	docker run -i --rm \
		-v "$(PWD):$(PWD)" \
		-v "$(realpath $(TESTGRID_CONFIG)):$(realpath $(TESTGRID_CONFIG))" \
		-v "$(GOOGLE_APPLICATION_CREDENTIALS):$(GOOGLE_APPLICATION_CREDENTIALS)" \
		-e "GOOGLE_APPLICATION_CREDENTIALS" \
		-w "$(PWD)" \
		gcr.io/k8s-prow/configurator:v20220124-9887456efc \
		"--oneshot" \
		"--output=gs://$(TESTGRID_GCS)/config" \
		"--yaml=$(realpath $(TESTGRID_CONFIG))"
