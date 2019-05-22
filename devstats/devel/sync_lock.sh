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

if [ -f "/tmp/deploy.wip" ]
then
  echo "another deploy process is running, exiting"
  exit 1
fi
wait_for_command.sh devstats 600 || exit 2
cronctl.sh devstats off || exit 3
cronctl.sh devstats.sh off || exit 4
if [ -z "$FROM_WEBHOOK" ]
then
  wait_for_command.sh webhook 600 || exit 5
  cronctl.sh webhook off || exit 6
  killall webhook
fi
echo 'All sync and deploy jobs stopped and disabled'
