// Copyright 2019 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package valimage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"

	"github.com/knative/test-infra/tools/monitoring/mail"
	"github.com/knative/test-infra/tools/monitoring/subscriber"
)

var (
	// monitoringSubs is a list of subscriptions to subscribe for image vulnerabilities
	monitoringSubs = [...]string{
		"sub-container-analysis-notes-v1beta1",
		"sub-container-analysis-occurrences-v1beta1",
	}
	recipients = []string{"knative-productivity-dev@googlegroups.com"}
)

// Client holds resources for monitoring image vulnerabilities
type Client struct {
	subClients []*subscriber.Client
	mailClient *mail.Config
}

// Setup initialize all the resources for monitoring image vulnerabilities
func Setup(mconfig *mail.Config) (*Client, error) {
	var subClients []*subscriber.Client
	for _, sub := range monitoringSubs {
		sub, err := subscriber.NewSubscriberClient(sub)
		if err != nil {
			return nil, err
		}
		subClients = append(subClients, sub)
	}

	return &Client{
		subClients: subClients,
		mailClient: mconfig,
	}, nil
}

// Run start a background process that listens to the message
func (c *Client) Run() {
	log.Println("Starting image vulnerabilities monitoring")
	for _, sub := range c.subClients {
		go func() {
			err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
				c.sendMessage(msg)
				msg.Ack()
			})
			if err != nil {
				log.Printf("Failed to receive messages due to: %v\n", err)
			}
		}()
	}
}

func (c *Client) sendMessage(msg *pubsub.Message) {
	err := c.mailClient.Send(recipients, "Image Vulnerabilities Detected", toMailContent(msg))
	if err != nil {
		log.Printf("Failed to send alert message %v\n", err)
	}
}

func toMailContent(msg *pubsub.Message) string {
	c := fmt.Sprintf("Message Data: %v\n", string(msg.Data))
	if b, err := json.MarshalIndent(msg, "", "\t"); err == nil {
		c += fmt.Sprintf("\nPubsub Message: %v\n", string(b))
	}
	c += fmt.Sprintf("\nRaw Message: %+v\n", msg)
	log.Printf("Mail Content:\n %v\n", c)
	return c
}
