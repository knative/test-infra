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

if [ -z "$ONLY" ]
then
  all="knative"
else
  all=$ONLY
fi
for proj in $all
do
    suff=$proj
    NOCOPY=1 GRAFANA=$suff devel/import_jsons_to_sqlite.sh grafana/dashboards/$proj/*.json || exit 2
done
echo 'OK'
