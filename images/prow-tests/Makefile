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

IMAGE_NAME = prow-tests
include ../simple-image.mk

push: confirm-main
	gcloud builds submit --config=./cloudbuild.yaml --substitutions=_TAG=$(TAG),_IMG=$(IMG),_COMMIT_HASH$(COMMIT_HASH) \
		--machine-type=n1-highcpu-8 \
		--project=$(PROJECT) --timeout=1h ./../..

# Build target for validating prow-tests image build in presubmit.
build-ci-test:
	gcloud builds submit --config=./cloudbuild-test.yaml --substitutions=_TAG=$(TAG),_IMG=$(IMG),_COMMIT_HASH=$(COMMIT_HASH)  \
		--machine-type=n1-highcpu-8 \
		--project=$(PROJECT) --timeout=1h ./../..
