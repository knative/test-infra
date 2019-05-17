# Copyright 2019 The Knative Authors
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

FROM golang:1.10.2
RUN apt-get update

ENV TEMP_TOOLS_REPO_DIR /go/src/github.com/knative/test-infra/tools
ENV TOOL_NAME monitoring

# Temporarily add test-infra to the image to build custom tools
ADD ./tools $TEMP_TOOLS_REPO_DIR

# These go dependencies need to be included before cloud mail files to avoid stale checkout due to
# no .git file in cloud mail
RUN go get -u cloud.google.com/go
RUN go get -u google.golang.org/genproto/googleapis/api/annotations
RUN go get -u google.golang.org/genproto/googleapis/iam/v1
RUN go get -u google.golang.org/genproto/protobuf/field_mask

#TODO(yt3liu): Use the alpha released after the go library builds correctly with the latest release
#RUN wget https://storage.googleapis.com/cloud-mail-client-libraries/cloudmail-v1alpha3-go.zip
RUN wget https://storage.googleapis.com/knative-monitoring/cloudmail-v1alpha3-go.zip

RUN apt-get install zip unzip
RUN unzip cloudmail-v1alpha3-go.zip -d cloudmail-v1alpha3-go
RUN apt-get -qq install rsync
RUN cp -r cloudmail-v1alpha3-go/google.golang.org/genproto/googleapis/cloud/mail /go/src/google.golang.org/genproto/googleapis/cloud/mail
RUN cp -r cloudmail-v1alpha3-go/cloud.google.com/go/mail /go/src/cloud.google.com/go/mail

RUN make -C $TEMP_TOOLS_REPO_DIR/$TOOL_NAME/
RUN cp $TEMP_TOOLS_REPO_DIR/$TOOL_NAME/$TOOL_NAME /$TOOL_NAME

# Remove test-infra from the container
RUN rm -fr $TEMP_TOOLS_REPO_DIR

ENTRYPOINT ["/monitoring"]
