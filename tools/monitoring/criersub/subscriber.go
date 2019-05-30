/*
Copyright 2019 The Knative Authors

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

package criersub

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

// ReportMessage is a message structure used to pass a prowjob status to Pub/Sub topic.
// TODO(yt3liu): Use the ReportMessage structure in "k8s.io/test-infra/prow/pubsub/reporter"
type ReportMessage struct {
	Project string `json:"project"`
	Topic   string `json:"topic"`
	RunID   string `json:"runid"`
	Status  string `json:"status"`
	URL     string `json:"url"`
	GCSPath string `json:"gcs_path"`
}

// SubscriberClient provide methods to run SubscriptionOperation
// It implements all methods in SubscriberOperation
type SubscriberClient struct {
	SubscriberOperation
}

// SubscriberOperation defines a list of methods for subscribing messages
type SubscriberOperation interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
	String() string
}

// NewSubscriberClient returns a new SubscriberClient used to read crier pubsub messages
func NewSubscriberClient(ctx context.Context, projectID string, subName string) (*SubscriberClient, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &SubscriberClient{c.Subscription(subName)}, err
}

// ReceiveMessageAckAll acknowledges all incoming pusub messages and convert the pubsub message to crier ReportMessage.
// It executes `f` only if the pubsub message can be converted to ReportMessage. Otherwise, ignore the message.
func (c *SubscriberClient) ReceiveMessageAckAll(ctx context.Context, f func(*ReportMessage)) error {
	return c.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if rmsg := toReportMessage(msg); rmsg != nil {
			f(rmsg)
		}
		msg.Ack()
	})
}

func toReportMessage(msg *pubsub.Message) *ReportMessage {
	rmsg := &ReportMessage{}
	if err := json.Unmarshal(msg.Data, rmsg); err != nil {
		log.Printf("Failed to convert message (%v) to ReportMessage\nError %v", msg, err)
		return nil
	}
	return rmsg
}
