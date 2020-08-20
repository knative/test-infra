package cluster

import (
	"log"

	"knative.dev/test-infra/pkg/clustermanager/boskos"
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
