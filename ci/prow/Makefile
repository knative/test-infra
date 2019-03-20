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

CLUSTER       ?= prow
PROJECT       ?= knative-tests
ZONE          ?= us-central1-f
JOB_NAMESPACE ?= test-pods

PROW_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TESTGRID_DIR := $(PROW_DIR)/../testgrid

get-cluster-credentials:
	gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"ßßß

config:
	go run make_config.go --prow-config-output="$(PROW_DIR)/config.yaml" --testgrid-config-output="$(TESTGRID_DIR)/config.yaml" config_knative.yaml

update-config: get-cluster-credentials
	kubectl create configmap config --from-file=config.yaml=config.yaml --dry-run -o yaml | kubectl replace configmap config -f -

update-plugins: get-cluster-credentials
	kubectl create configmap plugins --from-file=plugins.yaml=plugins.yaml --dry-run -o yaml | kubectl replace configmap plugins -f -

update-boskos: get-cluster-credentials
	kubectl apply -f boskos/config.yaml

update-boskos-config: get-cluster-credentials
	kubectl create configmap resources --from-file=config=boskos/resources.yaml --dry-run -o yaml | kubectl --namespace="$(JOB_NAMESPACE)" replace configmap resources -f -

update-cluster: get-cluster-credentials
	kubectl apply -f cluster.yaml

test:
	@echo "*** Checking config generator for prow and testgrid"
	$(eval TMP_YAML_PROW := $(shell mktemp)) $(eval TMP_YAML_TESTGRID := $(shell mktemp))
	go run make_config.go --prow-config-output="$(TMP_YAML_PROW)" --testgrid-config-output="$(TMP_YAML_TESTGRID)" config_knative.yaml
	diff config.yaml $(TMP_YAML_PROW)
	# TODO(chizhg): this line needs to be commented before we replace the testgrid config file with the auto-generated one
	# diff $(TESTGRID_DIR)/config.yaml $(TMP_YAML_TESTGRID)
	@echo "*** Checking config validity"
	bazel run @k8s//prow/cmd/checkconfig -- --plugin-config=$(PROW_DIR)/plugins.yaml --config-path=$(PROW_DIR)/config.yaml
