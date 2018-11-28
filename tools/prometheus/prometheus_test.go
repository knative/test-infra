/*
Copyright 2018 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package prometheus_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/knative/pkg/test/logging"
	"github.com/knative/test-infra/tools/prometheus"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	expected = 1.0
	query    = "test"
)

type TestPromAPI struct {
}

// AlertManagers returns an overview of the current state of the Prometheus alert manager discovery.
func (tpa *TestPromAPI) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	return v1.AlertManagersResult{}, nil
}

// CleanTombstones removes the deleted data from disk and cleans up the existing tombstones.
func (tpa *TestPromAPI) CleanTombstones(ctx context.Context) error {
	return nil
}

// Config returns the current Prometheus configuration.
func (tpa *TestPromAPI) Config(ctx context.Context) (v1.ConfigResult, error) {
	return v1.ConfigResult{}, nil
}

// DeleteSeries deletes data for a selection of series in a time range.
func (tpa *TestPromAPI) DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	return nil
}

// Flags returns the flag values that Prometheus was launched with.
func (tpa *TestPromAPI) Flags(ctx context.Context) (v1.FlagsResult, error) {
	return v1.FlagsResult{}, nil
}

// LabelValues performs a query for the values of the given label.
func (tpa *TestPromAPI) LabelValues(ctx context.Context, label string) (model.LabelValues, error) {
	return nil, nil
}

// Query performas a query on the prom api
func (tpa *TestPromAPI) Query(c context.Context, query string, ts time.Time) (model.Value, error) {
	fmt.Println("yay1")
	s := model.Sample{Value: expected}
	var v []*model.Sample
	v = append(v, &s)

	fmt.Println("yay")
	return model.Vector(v), nil
}

// QueryRange performs a query for the given range.
func (tpa *TestPromAPI) QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, error) {
	return nil, nil
}

// Series finds series by label matchers.
func (tpa *TestPromAPI) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, error) {
	return nil, nil
}

// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
// under the TSDB's data directory and returns the directory as response.
func (tpa *TestPromAPI) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	return v1.SnapshotResult{}, nil
}

// Targets returns an overview of the current state of the Prometheus target discovery.
func (t *TestPromAPI) Targets(ctx context.Context) (v1.TargetsResult, error) {
	return v1.TargetsResult{}, nil
}

// getTestAPI gets the test api implementation for prometheus api
func getTestAPI() *TestPromAPI {
	return &TestPromAPI{}
}

func TestRunQuery(t *testing.T) {
	logging.InitializeLogger(true)
	logger := logging.GetContextLogger("TestRunQuery")

	r := prometheus.RunQuery(context.Background(), logger, getTestAPI(), query)

	if r != expected {
		t.Fatalf("Expected: %f Actual: %f", expected, r)
	}
}
