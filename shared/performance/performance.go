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

package performance

import (
	"fmt"
	"io/ioutil"

	"github.com/knative/test-infra/shared/junit"
	"github.com/knative/test-infra/shared/mysql"
)

const (
	dbName     = "knative_performance"
	dbInstance = "knative-tests:us-central1:knative-monitoring"

	// Path to secrets for username and password
	userSecret = "/secrets/cloudsql/monitoringdb/username"
	passSecret = "/secrets/cloudsql/monitoringdb/password"

	// Property name used by testgrid
	perfLatency = "perf_latency"
)

// CreatePerfTestCase creates a perf test case with the provided name and value
func CreatePerfTestCase(metricValue float32, metricName, testName string) junit.TestCase {
	tp := []junit.TestProperty{{Name: perfLatency, Value: fmt.Sprintf("%f", metricValue)}}
	tc := junit.TestCase{
		ClassName:  testName,
		Name:       fmt.Sprintf("%s/%s", testName, metricName),
		Properties: junit.TestProperties{Properties: tp}}
	return tc
}

type DBConfig struct {
	mysql.DBConfig
}

// Configure the db instance to store metrics information.
// This will be later used to show the trending metrics on our grafana dashboard.
func ConfigureDB() (*DBConfig, error) {
	user, err := ioutil.ReadFile(userSecret)
	if err != nil {
		return nil, err
	}

	pass, err := ioutil.ReadFile(passSecret)
	if err != nil {
		return nil, err
	}

	config := mysql.DBConfig{
		Username:     string(user),
		Password:     string(pass),
		DatabaseName: dbName,
		Instance:     dbInstance,
	}

	return &DBConfig{config}, nil
}

func (c *DBConfig) StoreMetrics(name string, value float64) error {
	conn, err := c.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
