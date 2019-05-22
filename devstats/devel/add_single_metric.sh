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

if [ -z "${PG_DB}" ]
then
  echo "You need to set PG_DB environment variable to run this script"
  exit 1
fi
function finish {
    sync_unlock.sh
}
if [ -z "$TRAP" ]
then
  sync_lock.sh || exit -1
  trap finish EXIT
  export TRAP=1
fi
GHA2DB_CMDDEBUG=1 GHA2DB_RESETTSDB=1 GHA2DB_METRICS_YAML=devel/test_metrics.yaml GHA2DB_TAGS_YAML=devel/test_tags.yaml GHA2DB_LOCAL=1 ./gha2db_sync
