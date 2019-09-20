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

CUSTOM_ROLE="custom_role.yaml"
CUSTOM_ROLE_PERMS=/tmp/custom_role_permissions.yaml

rm -f ${CUSTOM_ROLE_PERMS}
for r in editor storage.admin; do
  gcloud iam roles describe roles/${r} | grep "^- " >> ${CUSTOM_ROLE_PERMS}
done
# If you want to add a new permission, add one line like below and uncomment
# echo "- container.secrets.create" >> ${CUSTOM_ROLE_PERMS}
echo 'title: "Knative Integration Tests Runner"' > ${CUSTOM_ROLE}
echo 'description: "The custom role used for Knative integration tests"' >> ${CUSTOM_ROLE}
echo 'stage: "GA"' >> ${CUSTOM_ROLE}
echo 'includedPermissions:' >> ${CUSTOM_ROLE}
cat ${CUSTOM_ROLE_PERMS} | sort | uniq >> ${CUSTOM_ROLE}
rm -f ${CUSTOM_ROLE_PERMS}
