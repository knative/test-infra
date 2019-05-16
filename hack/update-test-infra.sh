#!/usr/bin/env bash

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

set -o errexit
set -o nounset
set -o pipefail

echo
echo "This script updates the vendored test-infra in all Knative repos,"
echo "creating commits for each one of them. The PRs must still be created"
echo "through GitHub UI (just open the link given at the end of the process)."
echo "This script expects the Knative repositories to be located under"
echo "\$GOPATH/src/github.com/knative (as instructed by the development docs)."

cd ${GOPATH}
cd src/github.com/knative

for repo in *; do
  [[ "${repo}" == "test-infra" ]] && continue
  cd ${repo}
  echo -e "\n\n**** Updating test-infra in knative/${repo} ***\n\n"
  branch="update-test-infra-$(basename $(mktemp))"
  git checkout master
  git remote update -p
  git pull
  git checkout -b ${branch} upstream/master
  if [[ -f "Gopkg.lock" ]]; then
    dep ensure -update github.com/knative/test-infra
    ./hack/update-deps.sh
    [[ -z "$(git diff)" ]] && continue
    git commit -a -m "Update test-infra to the latest version"
    git push -u origin ${branch}
    echo -e "\nCheck the PR created above, and make changes if necessary"
  else
    echo -e "\nGopkg.lock not found, skip updating"
  fi
  echo -n "Hit [ENTER] to continue..."
  read
  cd ..
done
