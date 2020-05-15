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

FROM golang:1.14.2
LABEL maintainer="Chao Dai <chaodai@google.com>"

# Temporarily add test-infra to the image to build custom tools
COPY . /go/src/knative.dev/test-infra

RUN go install /go/src/knative.dev/test-infra/tools/monitoring/clearalerts
# TODO: jobs should not run '/clearalerts', just 'clearalerts'
RUN cp "$(which clearalerts)" /

RUN go install /go/src/knative.dev/test-infra/tools/monitoring
# TODO: jobs should not run '/monitoring', just 'monitoring'
RUN cp "$(which monitoring)" /

RUN go install /go/src/knative.dev/test-infra/tools/flaky-test-reporter
# TODO: jobs should not run '/flaky-test-reporter', just 'flaky-test-reporter'
RUN cp "$(which flaky-test-reporter)" /
RUN mkdir -p /config
RUN cp /go/src/knative.dev/test-infra/tools/flaky-test-reporter/config/config.yaml /config/config.yaml

RUN go install /go/src/knative.dev/test-infra/tools/flaky-test-retryer
# TODO: jobs should not run '/flaky-test-retryer', just 'flaky-test-retryer'
RUN cp "$(which flaky-test-retryer)" /

# Remove test-infra from the container
RUN rm -fr /go/src/knative.dev/test-infra
