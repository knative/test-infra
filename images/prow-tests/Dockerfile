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

# This is the base Dockerfile for knative-tests images, please remember to also
# do `make push` in `../prow-tests-*` directories to make sure these images are in
# sync

# TODO(chizhg): probably build the image from scratch?
FROM gcr.io/k8s-testimages/kubekins-e2e:v20200326-e946722-master

# Install extras on top of base image

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN gcloud components update

# Docker
RUN gcloud components install docker-credential-gcr
RUN docker-credential-gcr configure-docker

# Replace kubectl with a working version
# TODO(chizhg): remove once https://github.com/knative/test-infra/issues/1858 is fixed
ARG KUBECTL_VERSION=v1.17.4
RUN rm -f "$(which kubectl)" && \
    wget "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -O /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

# Extra tools through apt-get
RUN apt-get install -y uuid-runtime  # for uuidgen
RUN apt-get install -y rubygems      # for mdl
RUN apt-get install -y shellcheck    # for presubmit

# Extra tools through gem
RUN gem install mixlib-config -v 2.2.4  # required because ruby is 2.1
RUN gem install mdl -v 0.5  # required because later version of mdl requires ruby >= 2.4

# Install go 1.12, 1.13, and 1.14 using https://github.com/moovweb/gvm
# To use this, in Prow jobs set env var GO_VERSION="go1.14" (or go1.13, go1.12, etc)
# You can use e.g. "go1.14" even if we install a specific version here like "go1.14.2"
# (Obviously the version of Go must be installed on the image to be available -- we don't want to be downloading and installing it on every test iteration)
# Missing prerequisite:
RUN apt-get install -y bison
# GVM_NO_UPDATE_PROFILE=true means do not alter /root/.bashrc to automatically source gvm config, so when not using runner.sh, image works normally
# Install the tool:
RUN curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | GVM_NO_UPDATE_PROFILE=true bash
# gvm requires one to "source /root/.gvm/scripts/gvm" after installing
#  but in Dockerfile, each RUN is its own shell, so source'd in-RUN env vars are not propagated
# So have created "my-gvm.sh" to source the above then run gvm
COPY images/prow-tests/my-gvm.sh /usr/local/bin
# Install our versions of Go
RUN my-gvm.sh install go1.12.17 --prefer-binary
RUN my-gvm.sh install go1.13.10 --prefer-binary
RUN my-gvm.sh install go1.14.2 --prefer-binary

# Extra tools through go get
# These run using the kubekins version of Go, not any defined by `gvm`
RUN GO111MODULE=on go get github.com/google/ko/cmd/ko@v0.5.1
RUN GO111MODULE=on go get github.com/boz/kail/cmd/kail
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/jstemmer/go-junit-report
RUN GO111MODULE=on go get -u github.com/raviqqe/liche@v0.0.0-20200229003944-f57a5d1c5be4  # stable liche version for checking md links
RUN GO111MODULE=on go get -u github.com/google/go-licenses

# Extract versions
RUN bazel version > /bazel_version
RUN ko version > /ko_version

# Set default Go version
ENV GO_VERSION go1.13
# Note, it's pointless to run `my-gvm.sh use go1.13` in this Dockerfile because the PATH changes it makes won't stay in effect
# We wrap the runner with our own which will run `gvm use $GO_VERSION` for us.
RUN mv "$(which runner.sh)" "$(dirname "$(which runner.sh)")/kubekins-runner.sh"
COPY images/prow-tests/runner.sh /usr/local/bin

# Install our own tools
RUN go get knative.dev/pkg/testutils/junithelper

# And our scripts
ENV SCRIPTS /scripts
COPY scripts /scripts

# Temporarily add test-infra to the image to build custom tools
COPY . /go/src/knative.dev/test-infra

# Build custom tools in the container
# TODO: githubhelper is obsolete and no longer used in current test-infra (only old release branch scripts)
#       it should be deleted from the repo when those releases are no longer supported
RUN go install /go/src/knative.dev/test-infra/tools/githubhelper
# Because githubhelper was previously installed in a stupid way, it was used as `/workspace/githubhelper` so copy it
RUN cp "$(which githubhelper)" .

# RUN go install knative.dev/test-infra/kntest/cmd/kntest # Uncomment when you wish to have kntest installed

# reset --hard doesn't care about dirty working copy
RUN cd /go/src/knative.dev/test-infra && git reset --hard upstream/release-0.14 && go install knative.dev/test-infra/tools/dep-collector

# Cleanup
RUN rm -fr /go/src/knative.dev/test-infra
RUN rm -f /usr/local/bin/my-gvm.sh
