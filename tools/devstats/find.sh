#!/bin/bash

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

if [ -z "$1" ]
then
  echo "You need to provide path as first arument"
  exit 1
fi
if [ -z "$2" ]
then
  echo "You need to provide file name pattern as a first argument"
  exit 1
fi
if [ -z "$3" ]
then
  echo "You need to provide regexp pattern to search for as a second argument"
  exit 1
fi
find "$1" -type f -iname "$2" -not -name "out" -not -path '*.git/*' -exec grep -EHIn "$3" "{}" \; | tee -a out
