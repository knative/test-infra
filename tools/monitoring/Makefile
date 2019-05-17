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

all:
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/ptypes/empty
	go get -u github.com/golang/protobuf/ptypes/timestamp
	go get -u github.com/googleapis/gax-go/v2

	go get -u google.golang.org/api/iterator
	go get -u google.golang.org/api/option
	go get -u google.golang.org/api/transport
	go get -u google.golang.org/grpc
	go get -u google.golang.org/grpc/codes
	go get -u google.golang.org/grpc/metadata
	go get -u google.golang.org/grpc/status

	go get -u gopkg.in/yaml.v2

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
