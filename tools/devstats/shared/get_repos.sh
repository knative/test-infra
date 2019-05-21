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

if ( [ -z "$PG_PASS" ] || [ -z "$PG_DB" ] || [ -z "$GHA2DB_PROJECT" ] )
then
  echo "$0: you need to set GHA2DB_PROJECT, PG_DB and PG_PASS env variables to use this script"
  exit 1
fi
GHA2DB_PROJECTS_OVERRIDE="+$GHA2DB_PROJECT" GHA2DB_LOCAL=1 GHA2DB_PROCESS_COMMITS=1 GHA2DB_PROCESS_REPOS=1 GHA2DB_EXTERNAL_INFO=1 GHA2DB_PROJECTS_COMMITS="$GHA2DB_PROJECT" ./get_repos
