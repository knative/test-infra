# Copyright 2018 The Knative Authors
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

GUBERNATOR_SRC := test-infra/gubernator
APP_YAML := $(GUBERNATOR_SRC)/app.yaml

deploy:
	# Fetch latest source, patch for our instance
	rm -fr test-infra
	git clone http://github.com/kubernetes/test-infra.git
	cp config.yaml $(GUBERNATOR_SRC)
	sed -i -e '/^runtime: .*/a service: gubernator' $(APP_YAML)
	sed -i -e "/^handlers:/a\- url: /timeline\n  script: redir_github.app\n" $(APP_YAML)
	sed -i -e 's/user:kubernetes/user:knative/' $(GUBERNATOR_SRC)/view_pr.py
	sed -i -e 's/Kubernetes/Knative/' $(GUBERNATOR_SRC)/templates/index.html
	cp redir_github.py $(GUBERNATOR_SRC)
	# Deploy
	PROJECT=knative-tests make -C $(GUBERNATOR_SRC) deploy
	# Cleanup
	rm -fr test-infra
