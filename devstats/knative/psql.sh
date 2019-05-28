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

function finish {
    sync_unlock.sh
}
if [ -z "$TRAP" ]
then
  sync_lock.sh || exit -1
  trap finish EXIT
  export TRAP=1
fi
set -o pipefail
> errors.txt
> run.log
GHA2DB_PROJECT=knative PG_DB=knative GHA2DB_LOCAL=1 ./structure 2>>errors.txt | tee -a run.log || exit 1
sudo -u postgres psql knative -c "create extension if not exists pgcrypto" || exit 1
./devel/ro_user_grants.sh knative || exit 2
GHA2DB_PROJECT=knative PG_DB=knative GHA2DB_LOCAL=1 ./gha2db 2019-01-22 0 today now Knative 2>>errors.txt | tee -a run.log || exit 3
# You can get data even starting at 2012-07-01 but special calls are needed before 2015-01-01 - GitHub used different event format then.
# GHA2DB_PROJECT=knative PG_DB=knative GHA2DB_LOCAL=1 GHA2DB_OLDFMT=1 GHA2DB_EXACT=1 ./gha2db 2012-07-01 0 2014-12-31 23 Knative 2>>errors.txt | tee -a run.log || exit 4
GHA2DB_PROJECT=knative PG_DB=knative GHA2DB_LOCAL=1 GHA2DB_MGETC=y GHA2DB_SKIPTABLE=1 GHA2DB_INDEX=1 ./structure 2>>errors.txt | tee -a run.log || exit 5
GHA2DB_PROJECT=knative PG_DB=knative ./shared/setup_repo_groups.sh 2>>errors.txt | tee -a run.log || exit 6
GHA2DB_PROJECT=knative PG_DB=knative ./shared/setup_scripts.sh 2>>errors.txt | tee -a run.log || exit 7
GHA2DB_PROJECT=knative PG_DB=knative ./shared/get_repos.sh 2>>errors.txt | tee -a run.log || exit 8
GHA2DB_PROJECT=knative PG_DB=knative ./shared/import_affs.sh 2>>errors.txt | tee -a run.log || exit 9
GHA2DB_PROJECT=knative PG_DB=knative GHA2DB_LOCAL=1 ./vars || exit 10
