#!/usr/bin/env bash

# Copyright 2021 The Knative Authors
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

readonly JAVA_VERSION=16-open
export SDKMAN_DIR=/usr/local/sdkman

# Setup sdkman
curl -sL https://get.sdkman.io | bash

# Setup sdkman configs
echo sdkman_auto_answer=true > $SDKMAN_DIR/etc/config
echo sdkman_auto_selfupdate=true >> $SDKMAN_DIR/etc/config

# Load sdk command
source "$SDKMAN_DIR/bin/sdkman-init.sh"

# Install java and maven
sdk version
sdk update
sdk install java $JAVA_VERSION
sdk install maven

# Print versions
java -version
javac -version
mvn -v