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

project=$1
owners=("prime-engprod-sea@google.com")
groups=("knative-productivity-admins@googlegroups.com")
sas=("knative-tests@appspot.gserviceaccount.com" "prow-job@knative-tests.iam.gserviceaccount.com" "prow-job@knative-nightly.iam.gserviceaccount.com" "prow-job@knative-releases.iam.gserviceaccount.com") 
apis=("compute.googleapis.com" "container.googleapis.com")

# Add an owner to the project
for owner in ${owners[@]}; do
    gcloud projects add-iam-policy-binding $project --member group:$owner --role roles/OWNER
done

# Add all groups as editors
for group in ${groups[@]}; do
    gcloud projects add-iam-policy-binding $project --member group:$group --role roles/EDITOR
done

# Add all service accounts as editors
for sa in ${sas[@]}; do
    gcloud projects add-iam-policy-binding $project --member serviceAccount:$sa --role roles/EDITOR
done

# Enable apis
for api in ${apis[@]}; do
    gcloud services enable $api --project=$1
done
