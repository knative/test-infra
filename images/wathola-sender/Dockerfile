# Copyright 2020 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM docker.io/library/golang:1.13-alpine as builder

RUN apk add --update git alpine-sdk
RUN mkdir -p /code
ADD . /code/
WORKDIR /code/tools/wathola
RUN make sender

FROM scratch
COPY --from=builder /code/tools/wathola/build/wathola-sender /wathola-sender
WORKDIR /
CMD ["/wathola-sender"]
