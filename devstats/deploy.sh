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

mkdir /var/www 2>/dev/null
mkdir /var/www/html 2>/dev/null
./copy_devstats_binaries.sh || exit 1
cp cron/net_tcp_config.sh devel/sync_lock.sh devel/sync_unlock.sh $GOPATH/bin/ || exit 2
INIT=1 EXTERNAL=1 GHA2DB_GHAPISKIP=1 SKIPTEMP=1 ./devel/deploy_all.sh || exit 3
echo 'Deploy succeeded'
