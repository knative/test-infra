# Copyright 2022 The Knative Authors

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT       ?= knative-dns
OPTIONS       ?=

all: default slack blog stats testgrid dispatch

# For testing can deploy with --no-promote to avoid routing traffic to the new version
# i.e. make OPTIONS=--no-promote slack

default slack blog stats testgrid:
	@declare -A REDIRS=( \
	#### Redirection URLs for each subdomain #### \
	["default"]="https://knative.dev" \
	["testgrid"]="https://testgrid.k8s.io/r/knative-own-testgrid" \
	["slack"]="https://join.slack.com/t/knative/shared_invite/zt-fzy43r2f-i4oJraPBHlBFp00v6NeI9w" \
	["blog"]="https://knative.dev/blog" \
	["stats"]="https://knative.teststats.cncf.io/" \
	############################################## \
	) ; sed -e 's/$$SERVICE_NAME/$@/' -e "s*\$$REDIR_TO*$${REDIRS[$@]}*" redir.yaml > app.yaml
	gcloud app deploy --quiet --project="$(PROJECT)" $(OPTIONS) app.yaml
	@rm app.yaml

dispatch:
	gcloud app deploy --quiet --project="$(PROJECT)" $(OPTIONS) dispatch.yaml
