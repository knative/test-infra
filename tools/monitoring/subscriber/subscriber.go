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

package subscriber

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

// Refs describes how the repo was constructed.
// TODO(yt3liu): Use the ReportMessage structure in "k8s.io/test-infra/prow/pubsub/reporter"
// https://github.com/knative/test-infra/issues/912
type Refs struct {
	// Org is something like kubernetes or k8s.io
	Org string `json:"org"`
	// Repo is something like test-infra
	Repo string `json:"repo"`
	// RepoLink links to the source for Repo.
	RepoLink string `json:"repo_link,omitempty"`

	BaseRef string `json:"base_ref,omitempty"`
	BaseSHA string `json:"base_sha,omitempty"`
	// BaseLink is a link to the commit identified by BaseSHA.
	BaseLink string `json:"base_link,omitempty"`

	Pulls []Pull `json:"pulls,omitempty"`

	// PathAlias is the location under <root-dir>/src
	// where this repository is cloned. If this is not
	// set, <root-dir>/src/github.com/org/repo will be
	// used as the default.
	PathAlias string `json:"path_alias,omitempty"`
	// CloneURI is the URI that is used to clone the
	// repository. If unset, will default to
	// `https://github.com/org/repo.git`.
	CloneURI string `json:"clone_uri,omitempty"`
	// SkipSubmodules determines if submodules should be
	// cloned when the job is run. Defaults to true.
	SkipSubmodules bool `json:"skip_submodules,omitempty"`
	// CloneDepth is the depth of the clone that will be used.
	// A depth of zero will do a full clone.
	CloneDepth int `json:"clone_depth,omitempty"`
}

// Pull describes a pull request at a particular point in time.
// TODO(yt3liu): Use the ReportMessage structure in "k8s.io/test-infra/prow/pubsub/reporter"
// https://github.com/knative/test-infra/issues/912
type Pull struct {
	Number int    `json:"number"`
	Author string `json:"author"`
	SHA    string `json:"sha"`
	Title  string `json:"title,omitempty"`

	// Ref is git ref can be checked out for a change
	// for example,
	// github: pull/123/head
	// gerrit: refs/changes/00/123/1
	Ref string `json:"ref,omitempty"`
	// Link links to the pull request itself.
	Link string `json:"link,omitempty"`
	// CommitLink links to the commit identified by the SHA.
	CommitLink string `json:"commit_link,omitempty"`
	// AuthorLink links to the author of the pull request.
	AuthorLink string `json:"author_link,omitempty"`
}

// ReportMessage is a message structure used to pass a prowjob status to Pub/Sub topic.
// TODO(yt3liu): Use the ReportMessage structure in "k8s.io/test-infra/prow/pubsub/reporter"
// https://github.com/knative/test-infra/issues/912
type ReportMessage struct {
	Project string `json:"project"`
	Topic   string `json:"topic"`
	RunID   string `json:"runid"`
	Status  string `json:"status"`
	URL     string `json:"url"`
	GCSPath string `json:"gcs_path"`
	Refs    []Refs `json:"refs,omitempty"`
	JobType string `json:"job_type"`
	JobName string `json:"job_name"`
}

// Client is a wrapper on the subscriber Operation
type Client struct {
	Operation
}

// Operation defines a list of methods for subscribing messages
type Operation interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
	String() string
}

// NewSubscriberClient returns a new SubscriberClient used to read crier pubsub messages
func NewSubscriberClient(ctx context.Context, projectID string, subName string) (*Client, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Client{c.Subscription(subName)}, nil
}

// ReceiveMessageAckAll acknowledges all incoming pusub messages and convert the pubsub message to ReportMessage.
// It executes `f` only if the pubsub message can be converted to ReportMessage. Otherwise, ignore the message.
func (c *Client) ReceiveMessageAckAll(ctx context.Context, f func(*ReportMessage)) error {
	return c.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if rmsg, err := c.toReportMessage(msg); err != nil {
			log.Printf("Cannot convert pubsub message (%v) to Report message %v", msg, err)
		} else if rmsg != nil {
			f(rmsg)
		}
		msg.Ack()
	})
}

func (c *Client) toReportMessage(msg *pubsub.Message) (*ReportMessage, error) {
	rmsg := &ReportMessage{}
	if err := json.Unmarshal(msg.Data, rmsg); err != nil {
		return nil, err
	}
	return rmsg, nil
}
