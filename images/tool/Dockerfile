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

FROM golang:1.20.2

# Required input: -e TOOL_NAME=name, i.e. TOOL_NAME=flaky-test-reporter
ARG TOOL_NAME
# Optional input: -e TOOLS_SUBDIR=name, i.e. TOOLS_SUBDIR=monitoring/
ARG TOOLS_SUBDIR

# TODO: jobs should not run `/command`, just `command`

RUN GOBIN=/tools go install "knative.dev/toolbox/${TOOL_NAME}@latest"

ENTRYPOINT ["/tools/${TOOL_NAME}"]
