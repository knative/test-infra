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
  echo "Arguments required: path sha, none given"
  exit 1
fi
if [ -z "$2" ]
then
  echo "Arguments required: path sha, only path given"
  exit 2
fi

cd "$1" || exit 3
git show -s --format=%ct "$2" || exit 4
#files=`git diff-tree --no-commit-id --name-only -M8 -m -r "$2"` || exit 5
#files=`git diff-tree --no-commit-id --name-only -r "$2"` || exit 5
files=`git diff-tree --no-commit-id --name-only -M7 -r "$2"` || exit 5
for file in $files
do
    file_and_size=`git ls-tree -r -l "$2" "$file" | awk '{print $5 "♂♀" $4}'`
    if [ -z "$file_and_size" ]
    then
      echo "$file♂♀-1"
    else
      echo "$file_and_size"
    fi
done
