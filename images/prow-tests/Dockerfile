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

FROM gcr.io/k8s-testimages/kubekins-e2e:v20200827-e0bb92b-master

# Install extras on top of base image

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN gcloud components update

# Docker
RUN gcloud components install docker-credential-gcr
RUN docker-credential-gcr configure-docker

# Replace kubectl with a working version
# TODO(chizhg): https://github.com/knative/test-infra/issues/1858 was fixed in release-0.16, remove the pin after Knative release-0.19 is cut.
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
# So have created "source-gvm-and-run.sh" to source the above then run gvm
COPY images/prow-tests/source-gvm-and-run.sh /usr/local/bin
# Install our versions of Go
RUN source-gvm-and-run.sh install go1.13.14 --prefer-binary
RUN source-gvm-and-run.sh install go1.14.6 --prefer-binary
RUN source-gvm-and-run.sh install go1.15 --prefer-binary
RUN source-gvm-and-run.sh install master --source=https://github.com/golang/go --name=devel

# Extra tools through go get
# These run using the kubekins version of Go, not any defined by `gvm`
RUN GO111MODULE=on go get github.com/google/ko/cmd/ko@v0.5.1
RUN GO111MODULE=on go get github.com/boz/kail/cmd/kail
RUN go get -u github.com/golang/dep/cmd/dep
# TODO(chizhg): we started using gotestsum from release-0.17, remove go-junit-report after Knative release-0.20 is cut.
RUN go get -u github.com/jstemmer/go-junit-report
RUN GO111MODULE=on go get -u gotest.tools/gotestsum
RUN GO111MODULE=on go get -u github.com/raviqqe/liche@v0.0.0-20200229003944-f57a5d1c5be4  # stable liche version for checking md links
RUN GO111MODULE=on go get -u github.com/google/go-licenses
RUN GO111MODULE=on go get sigs.k8s.io/kind@v0.9.0
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-gke
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-kind
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-tester-exec

# protoc and required golang tooling
ARG PROTOC_VERSION=3.12.3
RUN wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -O protoc.zip \
    && unzip -p protoc.zip bin/protoc > /usr/local/bin/protoc \
    && chmod +x /usr/local/bin/protoc \
    && rm protoc.zip
RUN GO111MODULE=on go get github.com/gogo/protobuf/protoc-gen-gogofaster@v1.3.1

# Extract versions
RUN bazel version > /bazel_version
RUN ko version > /ko_version

# Note, it's pointless to run `source-gvm-and-run.sh use go1.13` in this Dockerfile because the PATH changes it makes won't stay in effect
# We wrap the runner with our own which will run `gvm use $GO_VERSION` for us.
RUN mv "$(which runner.sh)" "$(dirname "$(which runner.sh)")/kubekins-runner.sh"
COPY images/prow-tests/runner.sh /usr/local/bin

# TODO(chizhg): remove it once we replace with kntest in the places where this tool is used
# Install our own tools
RUN GO111MODULE=on go get knative.dev/pkg/testutils/junithelper@release-0.17

# And our scripts
ENV SCRIPTS /scripts
COPY scripts /scripts
COPY images/prow-tests/in-gvm-env.sh /usr/local/bin

# Temporarily add test-infra to the image to build custom tools
COPY . /go/src/knative.dev/test-infra

# Build custom tools in the container
# TODO: githubhelper is obsolete and no longer used in current test-infra (only old release branch scripts),
#       it should be removed after Knative release-0.18 is cut.
RUN go install /go/src/knative.dev/test-infra/tools/githubhelper
# Because githubhelper was previously installed in a stupid way, it was used as `/workspace/githubhelper` so copy it
RUN cp "$(which githubhelper)" .

# If you needed to compile and install different tools for different version of Go,
#  you could do it like:
#  RUN GO_VERSION=go1.13 in-gvm-env.sh go install knative.dev/test-infra/tools/coverage
#  RUN GO_VERSION=go1.14 in-gvm-env.sh go install knative.dev/test-infra/tools/coverage

# Install coverage tool:
RUN go install knative.dev/test-infra/tools/coverage
# This runs "runner.sh coverage" so coverage can be used with any installed Go

# Install kntest
RUN go install knative.dev/test-infra/kntest/cmd/kntest
# Install perf-tests tool
# TODO(chizhg): maybe also move perf-tests tool to be part of kntest?
RUN go install knative.dev/test-infra/pkg/clustermanager/perf-tests

# reset --hard doesn't care about dirty working copy
# This must be the last install
RUN cd /go/src/knative.dev/test-infra && git reset --hard upstream/release-0.14 && go install knative.dev/test-infra/tools/dep-collector

# Cleanup
RUN rm -fr /go/src/knative.dev/test-infra
RUN rm -f /usr/local/bin/source-gvm-and-run.sh
