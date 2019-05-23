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
  echo "$0: missing first parameter: lower case project name"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 1
fi

if [ -z "$2" ]
then
  echo "$0: missing second parameter: 'Project Name'"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 2
fi

if [ -z "$3" ]
then
  echo "$0: missing third parameter: 'project_org/main_repo'"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 3
fi

if [ -z "$4" ]
then
  echo "$0: missing fourth parameter: start-date in YYYY-MM-DD format (>= 2015-01-01)"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 4
fi

./set_project_name.sh homebrew Homebrew 'Homebrew/brew' '2015-01-01' "$1" "$2" "$3" "$4"
