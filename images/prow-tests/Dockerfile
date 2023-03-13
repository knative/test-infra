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

FROM debian:bullseye-20230227-slim AS base

# Pinned versions of stuff we pull in
ARG CLOUD_SDK_VERSION=405.0.0
ARG KUBECTL_VERSION=v1.24.4
ARG DOCKER_VERSION=5:20.10.17~3-0~debian-bullseye
ARG MAVEN_VERSION=3.8.4
ARG JAVA_VERSION=16
ARG PROTOC_VERSION=3.17.0
ARG TFENV_VERSION=v2.2.3
ARG KOPS_VERSION=v1.25.0
ARG COMMIT_HASH

RUN echo "${COMMIT_HASH}" > /commit_hash

WORKDIR /workspace
RUN mkdir -p /workspace
ENV WORKSPACE=/workspace \
    TERM=xterm
ENV DEBIAN_FRONTEND noninteractive

#
# BEGIN: GCLOUD SETUP
#

ENV PATH=/google-cloud-sdk/bin:/workspace:${PATH} \
    CLOUDSDK_CORE_DISABLE_PROMPTS=1

# net-tools is used by serving tests
RUN apt-get update -qqy && apt-get install -qqy \
        curl \
        gcc \
        python3-dev \
        python3-pip \
        apt-transport-https \
        lsb-release \
        openssh-client \
        ca-certificates \
        git \
        software-properties-common \
        bison \
        uuid-runtime \
        shellcheck \
        unzip \
        zip \
        wget \
        gnupg \
        jq \
        procps \
        net-tools \
        gnuplot

RUN pip3 install -U crcmod==1.7
RUN curl -fsSLO https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
RUN tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz -C /
RUN rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
RUN gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud components install alpha beta gke-gcloud-auth-plugin && \
    gcloud --version

#
# END: GCLOUD SETUP
#

# kubectl
RUN curl -fsSL "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

# kops
RUN curl -fsSL https://github.com/kubernetes/kops/releases/download/${KOPS_VERSION}/kops-linux-amd64 -o /usr/local/bin/kops && \
    chmod +x /usr/local/bin/kops

# tfenv
RUN git clone -b ${TFENV_VERSION} https://github.com/tfutils/tfenv.git ~/.tfenv && \
    ln -s ~/.tfenv/bin/* /usr/local/bin

#
# BEGIN: DOCKER IN DOCKER SETUP
#

# Add the Docker apt-repository
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
    $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
# TODO: the `sed` is a bit of a hack, look into alternatives.
# Why this exists: `docker service start` on debian runs a `cgroupfs_mount` method,
# We're already inside docker though so we can be sure these are already mounted.
# Trying to remount these makes for a very noisy error block in the beginning of
# the pod logs, so we just comment out the call to it... :shrug:
RUN apt-get update -qqy && \
    apt-get install -qqy --no-install-recommends docker-ce="${DOCKER_VERSION}" && \
    sed -i 's/cgroupfs_mount$/#cgroupfs_mount\n/' /etc/init.d/docker \
    && update-alternatives --set iptables /usr/sbin/iptables-legacy \
    && update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

# Move Docker's storage location
RUN echo 'DOCKER_OPTS="${DOCKER_OPTS} --data-root=/docker-graph"' | \
    tee --append /etc/default/docker
# NOTE this should be mounted and persisted as a volume ideally (!)
# We will make a fallback one now just in case
RUN mkdir /docker-graph

#
# END: DOCKER IN DOCKER SETUP
#

#
# BEGIN: JAVA SETUP (This is for kafka broker)
#

RUN curl -fsSL https://adoptopenjdk.jfrog.io/adoptopenjdk/api/gpg/key/public | gpg --dearmor -o /usr/share/keyrings/adoptopenjdk-archive-keyring.gpg && \
    echo \
    "deb [signed-by=/usr/share/keyrings/adoptopenjdk-archive-keyring.gpg] https://adoptopenjdk.jfrog.io/adoptopenjdk/deb \
    $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/adoptopenjdk.list
RUN apt-get update -qqy && \
    apt-get install -qqy adoptopenjdk-${JAVA_VERSION}-hotspot && \
    rm -rf /var/lib/apt/lists/*
ENV JAVA_HOME=/usr/lib/jvm/adoptopenjdk-${JAVA_VERSION}-hotspot-amd64

ENV MAVEN_HOME=/usr/local/maven
ENV M2_HOME=$MAVEN_HOME
ENV PATH=${M2_HOME}/bin:${PATH}

RUN curl -fsSL https://archive.apache.org/dist/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz -o /tmp/apache-maven-${MAVEN_VERSION}-bin.tar.gz && \
    tar xf /tmp/apache-maven-${MAVEN_VERSION}-bin.tar.gz -C /tmp && \
    mv /tmp/apache-maven-${MAVEN_VERSION} $MAVEN_HOME

RUN java -version && \
    mvn -v

#
# END: JAVA SETUP
#

# protoc and required golang tooling
RUN curl -fsSL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -o protoc.zip \
    && unzip -p protoc.zip bin/protoc > /usr/local/bin/protoc \
    && chmod +x /usr/local/bin/protoc \
    && rm protoc.zip
# protoc-gen-gogofaster is installed below

############################################################
FROM golang:1.20.2 AS external-go-gets
# Capture the version that depndabot bumps so that we can install it into the base image
RUN go version | cut -d' ' -f3 > /golang_version.latest

ARG KUBETEST2_VERSION=d306c412d528f0a5dda4b6bd1c7550a905f5edb9
ARG KUBETEST2_KOPS_VERSION=9b957a1d6189ee612506c9b6dc5b8b1a19e9c2aa # v1.25.0 release
ARG KIND_VERSION=v0.17.0
ARG KO_VERSION=v0.12.0
ARG PROTOC_GEN_GO_VERSION=v1.28.0
ARG PROTOC_GEN_GOGOFASTER_VERSION=v1.3.2
ARG GOTESTSUM_VERSION=v1.8.2
ARG KAIL_VERSION=v0.15.0
ARG GO_LICENSES_VERSION=v1.2.1
ARG DOCKER_CREDENTIAL_GCR_VERSION=v2.0.5
ARG JSONNET_VERSION=v0.18.0
ARG COSIGN_VERSION=v1.13.0

# Extra tools through go install
# These run using the golang image version of Go, not any defined by `gvm`
RUN go install github.com/google/ko@${KO_VERSION}
RUN go install github.com/boz/kail/cmd/kail@${KAIL_VERSION}
RUN go install gotest.tools/gotestsum@${GOTESTSUM_VERSION}
RUN go install sigs.k8s.io/kind@${KIND_VERSION}
RUN go install sigs.k8s.io/kubetest2@${KUBETEST2_VERSION}
RUN go install sigs.k8s.io/kubetest2/kubetest2-gke@${KUBETEST2_VERSION}
RUN go install sigs.k8s.io/kubetest2/kubetest2-kind@${KUBETEST2_VERSION}
RUN go install sigs.k8s.io/kubetest2/kubetest2-tester-exec@${KUBETEST2_VERSION}
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}
RUN go install github.com/google/go-licenses@${GO_LICENSES_VERSION}
RUN go install sigs.k8s.io/kind@${KIND_VERSION}
RUN go install github.com/google/go-jsonnet/cmd/jsonnet@${JSONNET_VERSION}
RUN go install github.com/sigstore/cosign/cmd/cosign@${COSIGN_VERSION}

# According to https://github.com/knative/test-infra/pull/2762 protoc-gen-gogofaster is unsupported (and shouldn't be used?)
# Not sure when it can be removed, maybe knative-0.26?
RUN go install github.com/gogo/protobuf/protoc-gen-gogofaster@${PROTOC_GEN_GOGOFASTER_VERSION}

# Docker
RUN go install github.com/GoogleCloudPlatform/docker-credential-gcr@${DOCKER_CREDENTIAL_GCR_VERSION}

# We do this instead of `go get` so the tools are locked to the Dockerfile's commit
COPY . /go/src/knative.dev/test-infra

# Build custom tools in the container
RUN cd /go/src/knative.dev/test-infra && go install ./tools/kntest/cmd/kntest
RUN cd /go/src/knative.dev/test-infra/tools/provenance-generator && go install .

# TODO(chizhg): maybe also move perf-tests tool to be part of kntest?
RUN cd /go/src/knative.dev/test-infra && go install ./pkg/clustermanager/perf-tests

# Install kubetest2-kops deployer
RUN git clone https://github.com/kubernetes/kops /tmp/kops && cd /tmp/kops/tests/e2e/kubetest2-kops && git checkout ${KUBETEST2_KOPS_VERSION} && \
    go install .

############################################################
FROM rust:1.67 AS external-rust-gets

ARG CODESIGN_VERSION=cb1e430a96f7d3d38af663a553b0e3e15a813d45
# change this when a new version with https://github.com/indygreg/apple-platform-rs/pull/20 is cut

RUN cargo install --git https://github.com/indygreg/apple-platform-rs --rev ${CODESIGN_VERSION} --bin rcodesign apple-codesign

############################################################
FROM base

COPY --from=external-rust-gets /usr/local/cargo/bin/rcodesign /usr/local/bin/rcodesign
COPY --from=external-go-gets /go/bin/* /go/bin/
COPY --from=external-go-gets /golang_version.latest /golang_version.latest

# Install go using https://github.com/moovweb/gvm
# GVM_NO_UPDATE_PROFILE=true means do not alter /root/.bashrc to automatically source gvm config, so when not using runner.sh, image works normally
# Install the tool:
RUN curl -fsSL https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | GVM_NO_UPDATE_PROFILE=true bash
# gvm requires one to "source /root/.gvm/scripts/gvm" after installing
#  but in Dockerfile, each RUN is its own shell, so source'd in-RUN env vars are not propagated
# So have created "source-gvm-and-run.sh" to source the above then run gvm
COPY images/prow-tests/source-gvm-and-run.sh /usr/local/bin
# Install our version of Go.

RUN source-gvm-and-run.sh install `cat /golang_version.latest` --prefer-binary
# The below will need to be updated to whatever the latest is when the project is ready.
RUN source-gvm-and-run.sh install go1.19.6 --prefer-binary
RUN source-gvm-and-run.sh use go1.19.6 --default

# Note, it's pointless to run `source-gvm-and-run.sh use goX.XX` in this Dockerfile because the PATH changes it makes won't stay in effect
# We wrap the runner with our own which will run `gvm use $GO_VERSION` for us.
COPY images/prow-tests/runner.sh /usr/local/bin

# Used if you want to install different programs for different versions of Go
COPY images/prow-tests/in-gvm-env.sh /usr/local/bin
# If you needed to compile and install different tools for different version of Go,
#  you could do it like:
#  RUN GO_VERSION=goX.XX in-gvm-env.sh go install knative.dev/test-infra/tools/cleanup
#  RUN GO_VERSION=goX.XX+1 in-gvm-env.sh go install knative.dev/test-infra/tools/cleanup
# But it must be done in base or the last FROM, because it does not install to /go/bin

# Not needed outside of this Dockerfile
RUN rm -f /usr/local/bin/source-gvm-and-run.sh

ENV PATH /go/bin:$PATH
ENV USE_GKE_GCLOUD_AUTH_PLUGIN True

# Extract versions
RUN ko version > /ko_version

# Ensure docker config is in the final image
RUN docker-credential-gcr configure-docker --registries=gcr.io,us-docker.pkg.dev
