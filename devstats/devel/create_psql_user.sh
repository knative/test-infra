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
  echo "$0: You need to set PG_PASS environment variable to run this script"
  exit 1
fi

if [ -z "$1" ]
then
  echo "$0: user name required"
  exit 1
fi

if [ -z "$ONLY" ]
then
  all="knative devstats"
else
  all=$ONLY
fi

cp ./util_sql/drop_psql_user.sql /tmp/drop_user.sql || exit 1
FROM="{{user}}" TO="$1" MODE=ss ./replacer /tmp/drop_user.sql || exit 1

if [ ! -z "$DROP" ]
then
  echo "Drop from public"
  sudo -u postgres psql < /tmp/drop_user.sql || exit 1
  for proj in $all
  do
    echo "Drop from $proj"
    sudo -u postgres psql "$proj" < /tmp/drop_user.sql || exit 1
  done
fi

if [ ! -z "$NOCREATE" ]
then
  echo "Skipping create"
  exit 0
fi

echo "Create role"
sudo -u postgres psql -c "create user \"$1\" with password '$PG_PASS'" || exit 1

for proj in $all
do
  echo "Grants $proj"
  ./devel/psql_user_grants.sh "$1" "$proj" || exit 2
done
echo 'OK'
