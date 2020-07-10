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

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(SELF_DIR)../../../common.mk

# Hardcode the cluster information instead of making them parameter, as there
# can only be a single instance of this tool, otherwise there will be N times
# comments left on PRs
CLUSTER := monitoring
PROJECT := knative-tests
ZONE := us-central1-a

DEPLOYMENT_YAML := retryer_service.yaml

.PHONY: deploy
deploy: confirm-master
ifdef GOOGLE_APPLICATION_CREDENTIALS
	gcloud auth activate-service-account --key-file="$(GOOGLE_APPLICATION_CREDENTIALS)"
endif
	gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"
	kubectl apply -f "$(DEPLOYMENT_YAML)"
	kubectl config unset current-context
