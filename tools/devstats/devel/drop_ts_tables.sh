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
  echo "$0: need database name argument"
  exit 1
fi
proj=$1
tables=`sudo -u postgres psql $proj -qAntc '\dt' | cut -d\| -f2`
for table in $tables
do
  base=${table:0:1}
  if ( [ "$base" = "t" ] || [ "$base" = "s" ] )
  then
    sudo -u postgres psql $proj -c "drop table $table" || exit 1
    echo "dropped $table"
  fi
done
