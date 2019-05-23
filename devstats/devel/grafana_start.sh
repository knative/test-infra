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
  echo "$0: you need to provide grafana name"
  exit 1
fi
./devel/grafana_stop.sh $1 || exit 1
cd /usr/share/grafana.$1
grafana-server -config /etc/grafana.$1/grafana.ini cfg:default.paths.data=/var/lib/grafana.$1 1>/var/log/grafana.$1.log 2>&1 &
