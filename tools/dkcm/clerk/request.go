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

package clerk

import (
	"time"
)

// Request stores a request for the "Request" table
type Request struct {
	*ClusterParams
	accessToken string
	requestTime time.Time
	ProwJobID   string
	ClusterID   int64
}

// Function option that modify a field of Request
type RequestOption func(*Request)

func NewRequest(opts ...RequestOption) *Request {
	r := &Request{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func AddRequestTime(requestTime time.Time) RequestOption {
	return func(r *Request) {
		r.requestTime = requestTime
	}
}

func AddProwJobID(prowJobID string) RequestOption {
	return func(r *Request) {
		r.ProwJobID = prowJobID
	}
}
