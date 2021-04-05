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

FROM gcr.io/k8s-testimages/kubekins-e2e:v20210113-cc576af-master AS stable-stuff

ENV DEBIAN_FRONTEND noninteractive
ARG GCLOUD_VERSION=317.0.0
ARG KUBECTL_VERSION=v1.18.7

# If we don't run update, the apt-get install below don't work
RUN apt-get update

# gcloud
RUN gcloud components update --version=${GCLOUD_VERSION}

# kubectl
RUN rm -f "$(which kubectl)" && \
    wget "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -O /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

# Docker
RUN gcloud components install docker-credential-gcr
RUN docker-credential-gcr configure-docker

# Extra tools through apt-get
RUN apt-get install -y uuid-runtime   # for uuidgen
RUN apt-get install -y rubygems       # for mdl
RUN gem install chef-utils -v 16.6.14 # for mdl, this is required for running the `gem install mdl` command below
RUN apt-get install -y shellcheck     # for shellcheck presubmit

# Extra tools through gem
# TODO: this tool isn't used anymore, remove when release-0.22 is cut; and don't install rubygems above either
RUN gem install mdl

# Install go 1.13, 1.14, and 1.15 using https://github.com/moovweb/gvm
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
RUN source-gvm-and-run.sh install go1.13.15 --prefer-binary
RUN source-gvm-and-run.sh install go1.14.12 --prefer-binary
RUN source-gvm-and-run.sh install go1.15.5 --prefer-binary

# protoc and required golang tooling
ARG PROTOC_VERSION=3.12.3
RUN wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -O protoc.zip \
    && unzip -p protoc.zip bin/protoc > /usr/local/bin/protoc \
    && chmod +x /usr/local/bin/protoc \
    && rm protoc.zip
# protoc-gen-gogofaster is installed in below section

# Note, it's pointless to run `source-gvm-and-run.sh use go1.13` in this Dockerfile because the PATH changes it makes won't stay in effect
# We wrap the runner with our own which will run `gvm use $GO_VERSION` for us.
RUN mv "$(which runner.sh)" "$(dirname "$(which runner.sh)")/kubekins-runner.sh"
COPY images/prow-tests/runner.sh /usr/local/bin

# Used if you want to install different programs for different versions of Go
COPY images/prow-tests/in-gvm-env.sh /usr/local/bin
# If you needed to compile and install different tools for different version of Go,
#  you could do it like:
#  RUN GO_VERSION=go1.13 in-gvm-env.sh go install knative.dev/test-infra/tools/coverage
#  RUN GO_VERSION=go1.14 in-gvm-env.sh go install knative.dev/test-infra/tools/coverage
# But it must be done in stable-stuff or the last FROM, because it does not install to /go/bin

############################################################
FROM stable-stuff AS external-go-gets

ARG KUBETEST2_VERSION=43699a7ba223ddce766f62ecd7504c7233605710

# Extra tools through go get
# These run using the kubekins version of Go, not any defined by `gvm`
RUN GO111MODULE=on go get github.com/google/ko/cmd/ko@v0.6.0
RUN GO111MODULE=on go get github.com/boz/kail/cmd/kail
RUN go get -u github.com/golang/dep/cmd/dep
# TODO(chizhg): we started using gotestsum from release-0.17, remove go-junit-report after Knative release-0.20 is cut.
RUN go get -u github.com/jstemmer/go-junit-report
RUN GO111MODULE=on go get -u gotest.tools/gotestsum
RUN GO111MODULE=on go get -u github.com/raviqqe/liche@v0.0.0-20200229003944-f57a5d1c5be4  # stable liche version for checking md links
RUN GO111MODULE=on go get sigs.k8s.io/kind@v0.9.0
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2@${KUBETEST2_VERSION}
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-gke@${KUBETEST2_VERSION}
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-kind@${KUBETEST2_VERSION}
RUN GO111MODULE=on go get -u sigs.k8s.io/kubetest2/kubetest2-tester-exec@${KUBETEST2_VERSION}

RUN GO111MODULE=on go get github.com/gogo/protobuf/protoc-gen-gogofaster@v1.3.1

# TODO(chizhg): remove it once we replace with kntest in the places where this tool is used
# Install our own tools
RUN GO111MODULE=on go get knative.dev/pkg/testutils/junithelper@release-0.17

############################################################
FROM stable-stuff AS test-infra-builds

# We do this instead of `go get` so the tools are locked to the Dockerfile's commit
COPY . /go/src/knative.dev/test-infra

# Build custom tools in the container

RUN go install knative.dev/test-infra/tools/coverage
RUN go install knative.dev/test-infra/kntest/cmd/kntest

# TODO(chizhg): maybe also move perf-tests tool to be part of kntest?
RUN go install knative.dev/test-infra/pkg/clustermanager/perf-tests

############################################################
FROM stable-stuff

COPY --from=external-go-gets /go/bin/* /go/bin/
COPY --from=test-infra-builds /go/bin/* /go/bin/

# Only needed in this Dockerfile
RUN rm -f /usr/local/bin/source-gvm-and-run.sh

# This tool is poorly written and requires its mod cache to be present to function!
# (The purpose of the multiple stages is to keep the mod cache out of the final image)
RUN GO111MODULE=on go get -u github.com/google/go-licenses

# Extract versions
RUN bazel version > /bazel_version
RUN ko version > /ko_version

