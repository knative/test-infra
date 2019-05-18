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

if [ -z "${PG_PASS}" ]
then
  echo "You need to set PG_PASS environment variable to run this script"
  exit 1
fi
if [ -z "$ONLY" ]
then
  all="homebrew"
else
  all=$ONLY
fi
for proj in $all
do
    db=$proj
    if [ "$proj" = "kubernetes" ]
    then
      db="gha"
    elif [ "$proj" = "all" ]
    then
      db="allprj"
    fi
    echo "Project: $proj, PDB: $db"
    sudo -u postgres psql "$db" -c "delete from gha_vars" || exit 1
    GHA2DB_LOCAL=1 GHA2DB_PROJECT=$proj PG_DB=$db ./vars || exit 2
done
echo 'OK'
