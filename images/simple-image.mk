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

# Rules for creating simple docker images.

# This should be included by another Makefile
#  and IMAGE_NAME set appropriately. You must also provide a Dockerfile
#  together with a Makefile in each subdirectory.

# Due to the relative path in the docker build commands, all image directories
#  must be a direct child of the images directory. If someone wanted to change
#  this, they'd need to write a little shell to calculate the repo root dir.

REGISTRY ?= gcr.io
PROJECT  ?= knative-tests

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(SELF_DIR)../common.mk

IMG = $(REGISTRY)/$(PROJECT)/test-infra/$(IMAGE_NAME)
TAG := $(shell date +v%Y%m%d)-$(shell git describe --always --dirty --match '^$$')

build:
	docker build --no-cache --pull -t $(IMG):$(TAG) -f Dockerfile ../..

# You can build locally without --no-cache to save time
iterative-build:
	docker build --pull -t $(IMG):local -f Dockerfile ../..

# And get a shell in the container
iterative-shell:
	docker run -it --entrypoint bash $(IMG_PATH):local

push_versioned: confirm-master build
	docker push $(IMG):$(TAG)

push_latest: confirm-master build
	docker tag $(IMG):$(TAG) $(IMG):latest
	docker push $(IMG):latest

push: push_versioned push_latest
