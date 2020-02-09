# Copyright 2019 The Knative Authors
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

# This is the canonical path of this Makefile.
STAGING_DIR := $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

# Staging cluster configuration

PROJECT                   := knative-tests-staging
PROW_GCS                  := knative-prow-staging
PROW_HOST                 := https://prow-staging.knative.dev
TESTGRID_GCS              := knative-testgrid-staging
GENERATE_MAINTENANCE_JOBS := false

KNATIVE_CONFIG   := $(STAGING_DIR)/config_staging.yaml
PROW_PLUGINS     := $(STAGING_DIR)/plugins.yaml
PROW_CONFIG      := $(STAGING_DIR)/config.yaml
TESTGRID_CONFIG  := $(STAGING_DIR)/testgrid.yaml
BOSKOS_RESOURCES := $(STAGING_DIR)/boskos_resources.yaml

include $(STAGING_DIR)/../Makefile

# Append to update-all-boskos-deployments: only 3 replicas of janitor
update-all-boskos-deployments::
	$(SET_CONTEXT)
	kubectl patch deployment boskos-janitor -n test-pods -p '{"spec":{"replicas":3}}'
	$(UNSET_CONTEXT)

update-all-cluster-deployments::
	$(SET_CONTEXT)
	$(eval host := $(subst https://,,$(PROW_HOST)))
	kubectl get ingress -oyaml | sed "s/prow.knative.dev/$(host)/" | kubectl apply -f -
	$(UNSET_CONTEXT)
