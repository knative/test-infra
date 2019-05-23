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

if [ -z "$GRAFANA" ]
then
  echo "$0: you need to set GRAFANA env variable (Grafana suffix). For example k8s, all, prometheus etc."
  exit 1
fi
if [ -z "$1" ]
then
  echo "$0: you need to provide at least one json to import"
  exit 2
fi
cp /var/lib/grafana.$GRAFANA/grafana.db ./grafana.$GRAFANA.db || exit 4
./sqlitedb ./grafana.$GRAFANA.db $* || exit 5
./devel/grafana_stop.sh $GRAFANA || exit 6
cp ./grafana.$GRAFANA.db /var/lib/grafana.$GRAFANA/grafana.db || exit 7
./devel/grafana_start.sh $GRAFANA || exit 8
echo "OK, if all is fine delete grafana.$GRAFANA.db.* db backup files and *.was json backup files".
echo "Otherwise use grafana.$GRAFANA.db.* backup file to restore previous Grafana DB."
