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

# ARTWORK
# GET=1 (attempt to fetch Postgres database and Grafana database from the test server)
# INIT=1 (needs PG_PASS_RO, PG_PASS_TEAM, initialize from no postgres database state, creates postgres logs database and users)
# SKIPVARS=1 (if set it will skip final Postgres vars regeneration)
set -o pipefail
exec > >(tee run.log)
exec 2> >(tee errors.txt)
if [ -z "$PG_PASS" ]
then
  echo "$0: You need to set PG_PASS environment variable to run this script"
  exit 1
fi
if ( [ ! -z "$INIT" ] && ( [ -z "$PG_PASS_RO" ] || [ -z "$PG_PASS_TEAM" ] ) )
then
  echo "$0: You need to set PG_PASS_RO, PG_PASS_TEAM when using INIT"
  exit 1
fi

GRAF_USRSHARE="/usr/share/grafana"
GRAF_VARLIB="/var/lib/grafana"
GRAF_ETC="/etc/grafana"
export GRAF_USRSHARE
export GRAF_VARLIB
export GRAF_ETC

host=`hostname`
function finish {
    sync_unlock.sh
    rm -f /tmp/deploy.wip 2>/dev/null
}
if [ -z "$TRAP" ]
then
  sync_lock.sh || exit -1
  trap finish EXIT
  export TRAP=1
  > /tmp/deploy.wip
fi

if [ ! -z "$INIT" ]
then
  ./devel/init_database.sh || exit 1
fi

PROJ=knative PROJDB=knative PROJREPO="Knative/serving" ORGNAME="Knative" PORT=3001 ICON="-" GRAFSUFF=knative GA="-" ./devel/deploy_proj.sh || exit 2

if [ -z "$SKIPVARS" ]
then
  ./devel/vars_all.sh || exit 3
fi
echo "$0: All deployments finished"
