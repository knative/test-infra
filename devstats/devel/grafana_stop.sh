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
pid=`ps -axu | grep grafana-server | grep $1 | awk '{print $2}'`
echo "stopping $1 grafana server instance"
if [ ! -z "$pid" ]
then
  echo "stopping pid $pid"
  kill $pid
else
  echo "grafana-server $1 not running"
fi
