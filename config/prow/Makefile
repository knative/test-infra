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

# Default settings for the CI/CD system.

CLUSTER       ?= prow
PROJECT       ?= knative-tests
ZONE          ?= us-central1-f
JOB_NAMESPACE ?= test-pods
PROW_HOST     ?= https://prow.knative.dev

SKIP_CONFIG_BACKUP        ?=
GENERATE_MAINTENANCE_JOBS ?= true
GENERATE_TESTGRID_CONFIG  ?= true

CONFIG_GENERATOR_DIR ?= ../../tools/config-generator

KNATIVE_CONFIG   ?= config_knative.yaml
PROW_PLUGINS     ?= core/plugins.yaml
PROW_CONFIG      ?= core/config.yaml
PROW_JOB_CONFIG  ?= jobs/config.yaml
PROW_DEPLOYS     ?= cluster
PROW_GCS         ?= knative-prow
PROW_CONFIG_GCS  ?= gs://$(PROW_GCS)/configs

TESTGRID_GCS    ?= knative-testgrid
TESTGRID_CONFIG ?= testgrid/testgrid.yaml

BOSKOS_RESOURCES ?= boskos/boskos_resources.yaml

# Useful shortcuts.

SET_CONTEXT   := gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"
UNSET_CONTEXT := kubectl config unset current-context
MAKE_CONFIG   := go run "${CONFIG_GENERATOR_DIR}"

.PHONY: help get-cluster-credentials unset-cluster-credentials
help:
	@echo "Help"
	@echo "'Update' means updating the servers and can only be run by oncall staff."
	@echo "Common usage:"
	@echo " make update-prow: Update all Prow things  except Boskos resources for some reason"
	@echo " make update-prow-job-config: Update just the Prow job config"
	@echo " make jobs/config.yaml: Generate Prow job config.yaml"
	@echo " make update-all-cluster-deployments: Update all Prow deployments except Boskos for some reason"
	@echo " make update-all-boskos-deployments: Update Boskos deployments"
	@echo " make update-boskos-resource: Update the Boskos resources"
	@echo " make update-testgrid-config: Update the Testgrid config"

# Useful general targets.
get-cluster-credentials:
	$(SET_CONTEXT)

unset-cluster-credentials:
	$(UNSET_CONTEXT)

# Generate the Prow config file
jobs/config.yaml:
	$(MAKE_CONFIG) \
		--gcs-bucket=$(PROW_GCS) \
		--generate-testgrid-config=$(GENERATE_TESTGRID_CONFIG) \
		--generate-maintenance-jobs=$(GENERATE_MAINTENANCE_JOBS) \
		--image-docker=gcr.io/$(PROJECT)/test-infra \
		--plugins-config-output=$(PROW_PLUGINS) \
		--prow-config-output=$(PROW_CONFIG) \
		--prow-jobs-config-output=$(PROW_JOB_CONFIG) \
		--prow-host=$(PROW_HOST) \
		--testgrid-config-output=$(TESTGRID_CONFIG) \
		--testgrid-gcs-bucket=$(TESTGRID_GCS) \
		$(KNATIVE_CONFIG)

.PHONY: update-prow-config update-prow-job-config update-prow-plugins update-all-boskos-deployments update-boskos-resource update-all-cluster-deployments update-single-cluster-deployment update-prow test update-testgrid-config

# Update prow config
update-prow-config:
	$(SET_CONTEXT)
	kubectl create configmap config --from-file=config.yaml=$(PROW_CONFIG) --dry-run --save-config -o yaml | kubectl apply -f -
	$(UNSET_CONTEXT)

# Update all prow job configs
update-prow-job-config: jobs/config.yaml
	$(SET_CONTEXT)
ifndef SKIP_CONFIG_BACKUP
	$(eval OLD_YAML_CONFIG := $(shell mktemp))
	$(eval NEW_YAML_CONFIG := $(shell mktemp))
	$(eval GCS_DEST := $(PROW_CONFIG_GCS)/config-$(shell date '+%Y_%m_%d_%H:%M:%S').yaml)
	@kubectl get configmap job-config -o jsonpath="{.data['config\.yaml']}" 2>/dev/null > "${OLD_YAML_CONFIG}"
	@gsutil cp "${OLD_YAML_CONFIG}" "${GCS_DEST}" > /dev/null
	@sed -e "s/((COMMIT_HASH_TOKEN))/$(shell git rev-parse HEAD)/" $(PROW_JOB_CONFIG) > "${NEW_YAML_CONFIG}"
else
	$(eval NEW_YAML_CONFIG := $(PROW_JOB_CONFIG))
endif
# We'll have to use `kubectl replace` here because `kubectl apply` will add the whole original configmap to the
# `last-applied-configuration` annotation, but k8s has a maxmium 256kb limit for it, see https://github.com/kubernetes/kubectl/issues/712
	kubectl create configmap job-config --from-file=config.yaml=$(NEW_YAML_CONFIG) --dry-run -o yaml | kubectl replace configmap job-config -f -
	$(UNSET_CONTEXT)
ifndef SKIP_CONFIG_BACKUP
	diff "${OLD_YAML_CONFIG}" "${NEW_YAML_CONFIG}" --color=auto || true
	@echo "Inspect uploaded config file at: ${NEW_YAML_CONFIG}"
	@echo "Old config file saved at: ${GCS_DEST}"
endif

# Update prow plugins
update-prow-plugins:
	$(SET_CONTEXT)
	kubectl create configmap plugins --from-file=plugins.yaml=$(PROW_PLUGINS) --dry-run --save-config -o yaml | kubectl apply -f -
	$(UNSET_CONTEXT)

# Update all deployments of boskos
update-all-boskos-deployments:
	$(SET_CONTEXT)
	@for f in $(wildcard $(PROW_DEPLOYS)/*boskos*.yaml); do kubectl apply -f $${f}; done
	$(UNSET_CONTEXT)

# Update the list of resources for Boskos
update-boskos-resource:
	$(SET_CONTEXT)
	kubectl create configmap resources --from-file=config=$(BOSKOS_RESOURCES) --dry-run --save-config -o yaml | kubectl --namespace="$(JOB_NAMESPACE)" apply -f -
	$(UNSET_CONTEXT)

# Update all deployments of cluster except Boskos for some reason
update-all-cluster-deployments:
	$(SET_CONTEXT)
	@for f in $(filter-out $(wildcard $(PROW_DEPLOYS)/*boskos*.yaml),$(wildcard $(PROW_DEPLOYS)/*.yaml)); do kubectl apply -f $${f}; done
	$(UNSET_CONTEXT)

# Update single deployment of cluster, expect passing in ${NAME} like `make update-single-cluster-deployment NAME=crier_deployment`
update-single-cluster-deployment:
	$(SET_CONTEXT)
	kubectl apply -f $(PROW_DEPLOYS)/$(NAME).yaml
	$(UNSET_CONTEXT)

# Update all resources on Prow cluster except boskos resource
update-prow: update-all-cluster-deployments update-all-boskos-deployments update-prow-plugins update-prow-config update-prow-job-config

# Ensure that configs are valid and up-to-date
test:
	@echo "*** Checking config generator for prow and testgrid"
	$(eval TMP_YAML_PROW := $(shell mktemp))
	$(eval TMP_YAML_PROW_JOBS := $(shell mktemp))
	$(eval TMP_YAML_PLUGINS := $(shell mktemp))
	$(eval TMP_YAML_TESTGRID := $(shell mktemp))
	$(MAKE_CONFIG) \
		--gcs-bucket=$(PROW_GCS) \
		--generate-maintenance-jobs=$(GENERATE_MAINTENANCE_JOBS) \
		--image-docker=gcr.io/$(PROJECT)/test-infra \
		--plugins-config-output=$(TMP_YAML_PLUGINS) \
		--prow-config-output=$(TMP_YAML_PROW) \
		--prow-jobs-config-output=$(TMP_YAML_PROW_JOBS) \
		--prow-host=$(PROW_HOST) \
		--testgrid-config-output=$(TMP_YAML_TESTGRID) \
		--testgrid-gcs-bucket=$(TESTGRID_GCS) \
		$(KNATIVE_CONFIG)
	diff --ignore-matching-lines="^# Copyright " $(PROW_CONFIG) $(TMP_YAML_PROW)
	diff --ignore-matching-lines="^# Copyright " $(PROW_JOB_CONFIG) $(TMP_YAML_PROW_JOBS)
	diff --ignore-matching-lines="^# Copyright " $(PROW_PLUGINS) $(TMP_YAML_PLUGINS)
	diff --ignore-matching-lines="^# Copyright " $(TESTGRID_CONFIG) $(TMP_YAML_TESTGRID)
	@echo "*** Checking configs validity"
	bazel run @k8s//prow/cmd/checkconfig -- \
		--config-path=$(realpath $(PROW_CONFIG)) \
		--job-config-path=$(realpath $(PROW_JOB_CONFIG)) \
		--plugin-config=$(realpath $(PROW_PLUGINS))
	bazel run @k8s//testgrid/cmd/configurator -- \
		--validate-config-file \
		--yaml=$(realpath $(TESTGRID_CONFIG)) \
# Ensure TestGrid config can be converted to binary form.
	bazel run @k8s//testgrid/cmd/configurator -- \
		--oneshot \
		--output=/dev/null \
		--yaml=$(realpath $(TESTGRID_CONFIG))

# Update TestGrid config.
# Application Default Credentials must be set, otherwise the upload will fail.
# Either export $GOOGLE_APPLICATION_CREDENTIALS pointing to a valid service
# account key, or temporarily use your own credentials by running
# gcloud auth application-default login
update-testgrid-config:
	bazel run @k8s//testgrid/cmd/configurator -- \
		--oneshot \
		--output=gs://$(TESTGRID_GCS)/config \
		--yaml=$(realpath $(TESTGRID_CONFIG))
