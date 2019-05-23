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

if ( [ -z "$GHA2DB_PROJECT" ] || [ -z "$PG_DB" ] || [ -z "$PG_PASS" ] )
then
  echo "$0: you need to set GHA2DB_PROJECT, PG_DB and PG_PASS env variables to use this script"
  exit 1
fi
proj=$GHA2DB_PROJECT
echo "Setting up $proj repository groups sync script"
sudo -u postgres psql $PG_DB -c "insert into gha_postprocess_scripts(ord, path) select 0, 'scripts/$proj/repo_groups.sql' on conflict do nothing"
echo "Setting $proj up default postprocess scripts"
./runq util_sql/default_postprocess_scripts.sql
echo "Setting $proj up repository groups postprocess script"
./runq util_sql/repo_groups_postprocess_script_from_repos.sql
