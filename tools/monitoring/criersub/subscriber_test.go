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
	"errors"
	"reflect"
	"testing"

	"cloud.google.com/go/pubsub"
)

type contextKey int

const (
	keyError contextKey = iota
	keyExist
)

type fakeSubscriber struct {
	SubscriberOperation
	name string
}

func getFakeSubscriber(n string) *SubscriberClient {
	return &SubscriberClient{&fakeSubscriber{name: n}}
}

func (fs *fakeSubscriber) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	if err := ctx.Value(keyError); err != nil {
		return err.(error)
	}
	return nil
}

func (fs *fakeSubscriber) Exists(ctx context.Context) (bool, error) {
	exist := false
	if e := ctx.Value(keyExist); e != nil {
		exist = e.(bool)
	}

	if err := ctx.Value(keyError); err != nil {
		return exist, err.(error)
	}

	return exist, nil
}

func (fs *fakeSubscriber) String() string {
	return fs.name
}

func TestSubscriberClient_ReceiveMessageAckAll(t *testing.T) {
	receivedMsgs := make([]*ReportMessage, 3)

	type arguments struct {
		ctx context.Context
		f   func(*ReportMessage)
	}
	tests := []struct {
		name string
		args arguments
		want error
	}{
		{
			name: "Message Received",
			args: arguments{
				ctx: context.Background(),
				f: func(message *ReportMessage) {
					receivedMsgs = append(receivedMsgs, message)
				},
			},
			want: nil,
		},
		{
			name: "ReceiveError",
			args: arguments{
				ctx: context.WithValue(context.Background(), keyError, errors.New("receive failed")),
				f: func(message *ReportMessage) {
					receivedMsgs = append(receivedMsgs, message)
				},
			},
			want: errors.New("receive failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := getFakeSubscriber("fake subscriber")
			if got := fs.ReceiveMessageAckAll(tt.args.ctx, tt.args.f); !isSameError(got, tt.want) {
				t.Errorf("ReceiveMessageAutoAck(%v), got: %v, want: %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestValidateSubscription(t *testing.T) {
	type arguments struct {
		ctx context.Context
		sub SubscriberOperation
	}
	tests := []struct {
		name string
		sub  SubscriberOperation
		args arguments
		want error
	}{
		{
			name: "Valid Subscription",
			args: arguments{
				ctx: context.WithValue(context.Background(), keyExist, true),
				sub: getFakeSubscriber("test-sub"),
			},
			want: nil,
		},
		{
			name: "Non-Existent Subscription",
			args: arguments{
				ctx: context.WithValue(context.Background(), keyExist, false),
				sub: getFakeSubscriber("non-exist"),
			},
			want: errors.New("the subscription (non-exist) does not exist"),
		},
		{
			name: "Exists Error",
			args: arguments{
				ctx: context.WithValue(context.Background(), keyError, errors.New("subscriber Exists failed")),
				sub: getFakeSubscriber("non-exist"),
			},
			want: errors.New("subscriber Exists failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateSubscription(tt.args.ctx, tt.args.sub); !isSameError(got, tt.want) {
				t.Errorf("validateSubscription(%v), got: %v, want: %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestToReportMessage(t *testing.T) {
	tests := []struct {
		name string
		arg  *pubsub.Message
		want *ReportMessage
	}{
		{
			name: "Valid report",
			arg: &pubsub.Message{
				Data: []byte(`{"project":"knative-tests","topic":"knative-monitoring","runid":"post-knative-serving-go-coverage-dev","status":"triggered","url":"","gcs_path":"gs://","refs":[{"org":"knative","repo":"serving","base_ref":"master","base_sha":"ce96dd74b1c85f024d63ce0991d4bf61aced582a","clone_uri":"https://github.com/knative/serving.git"}],"job_type":"postsubmit","job_name":"post-knative-serving-go-coverage-dev"}`)},
			want: &ReportMessage{
				Project: "knative-tests",
				Topic:   "knative-monitoring",
				RunID:   "post-knative-serving-go-coverage-dev",
				Status:  "triggered",
				URL:     "",
				GCSPath: "gs://",
			},
		},
		{
			name: "Invalid Report",
			arg: &pubsub.Message{
				Data: []byte(`Random Weird Format`)},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toReportMessage(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toReportMessage(%v), got: %v, want: %v", tt.arg, got, tt.want)
			}
		})
	}
}

func isSameError(err1 error, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	} else if err1 == nil || err2 == nil {
		return false
	}
	return err1.Error() == err2.Error()
}
