# Copyright 2020 The Knative Authors
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

# Do not allow server update from wrong branch or dirty working space
# In emergency, could easily edit this file, deleting all these lines
confirm-master:
	@if git diff-index --quiet HEAD; then true; else echo "Git working space is dirty -- will not update server"; false; fi;
# TODO(chizhg): change to `git branch --show-current` after we update the Git version in prow-tests image.
ifneq ("$(shell git rev-parse --abbrev-ref HEAD)","master")
	@echo "Branch is not master -- will not update server"
	@false
endif

