/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cluster

import (
	"log"

	"knative.dev/test-infra/pkg/clustermanager/e2e-tests/boskos"
	"knative.dev/test-infra/pkg/cmd"
)

type Client struct {
	Boskos *boskos.Client
}

func NewClient() (*Client, error) {
	boskosCli, err := boskos.NewClient("", "", "")
	if err != nil {
		return nil, err
	}

	return &Client{Boskos: boskosCli}, nil
}

func (c *Client) CreateCluster() error {
	log.Printf("creating a cluster...")
	cmd.RunCommand("kubetest2 gke --region=xxx --back-regions=xx --size=xx --up")
	return nil
}
