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

# Required rules for building kubernetes/test-infra
# These all come from http://github.com/kubernetes/test-infra/blob/master/WORKSPACE

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains(go_version = "1.11.4")

git_repository(
    name = "io_bazel_rules_k8s",
    commit = "9d2f6e8e21f1b5e58e721fc29b806957d9931930",
    remote = "https://github.com/bazelbuild/rules_k8s.git",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "5235045774d2f40f37331636378f21fe11f69906c0386a790c5987a09211c3c4",
    strip_prefix = "rules_docker-8010a50ef03d1e13f1bebabfc625478da075fa60",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/8010a50ef03d1e13f1bebabfc625478da075fa60.tar.gz"],
)

# External repositories

git_repository(
    name = "k8s",
    remote = "http://github.com/kubernetes/test-infra.git",
    commit = "c4a1fe42ebf91a06a81b189223d19f7e4332634b",  # HEAD as of 12/15/2018
)

