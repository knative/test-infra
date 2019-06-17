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

package gcs

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
)

// Client is the client used to interact with google cloud storage
type Client struct {
	*storage.Client
}

// NewClient creates a new storage Client
func NewClient(ctx context.Context) (*Client, error) {
	c, err := storage.NewClient(ctx)
	return &Client{c}, err
}

// ReadFromLink reads from a gsUrl and return a log structure
func (c Client) ReadFromLink(ctx context.Context, gsURL string) ([]byte, error) {
	var data []byte

	bucket, obj, err := linkToBucketAndObject(gsURL)
	if err != nil {
		return data, err
	}

	return c.read(ctx, bucket, obj)
}

func (c Client) read(ctx context.Context, bucket string, object string) ([]byte, error) {
	rc, err := c.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func linkToBucketAndObject(gsURL string) (string, string, error) {
	var bucket, obj string
	gsURL = strings.Replace(gsURL, "gs://", "", 1)

	sIdx := strings.IndexByte(gsURL, '/')
	if sIdx == -1 || sIdx+1 >= len(gsURL) {
		return bucket, obj, fmt.Errorf("the gsUrl (%s) cannot be converted to bucket/object", gsURL)
	}

	return gsURL[:sIdx], gsURL[sIdx+1:], nil
}
