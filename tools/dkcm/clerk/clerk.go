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
	"database/sql"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"
	"knative.dev/test-infra/pkg/mysql"
)

// Cluster stores a row in the "Cluster" db table
type Cluster struct {
	ClusterID string
	ProjectID string
	Status    string
	Zone      string
	Nodes     int
	NodeType  string
}

// Request stores a request for the "Request" table
type Request struct {
	accessToken string
	requestTime time.Time
	Zone        string
	Nodes       int
	NodeType    string
	ProwJobID   string
	ClusterID   string
}

// ClusterConfig is the struct to return to Prow once the cluster is available
type ClusterConfig struct {
	ClusterName string
	ProjectID   string
	Zone        string
}

type ClerkOperations interface {
	//check cluster available, if available, return cluster id and access token
	CheckAvail(nodes int, nodeType string, zone string) (bool, string, string)
	// get with token
	GetCluster(accessToken string) (ClusterConfig, error)
	// delete a cluster entry
	DeleteCluster(accessToken string) error
	// Insert a cluster entry
	InsertCluster(nodes int, nodeType string, projectID string, status string, zone string) error
	// Update a cluster entry
	UpdateCluster(clusterID string, fields []string, values []string) error
	// List clutsers within a time interval (use for checking after downtime to see stale clusters)
	ListClusters(window time.Duration) ([]Cluster, error)
	// Generate Unique access token for Prow
	generateToken() string
	// Insert a request entry
	InsertRequest(Zone string, Nodes int, NodeType string, ProwJobID string) error
	// Update a request entry
	UpdateRequest(requestID string, fields []string, values []string) error
	// List requests within a time interval (use for checking after downtime to see stale requests)
	ListRequests(window time.Duration) ([]Request, error)
}

// DB holds an active database connection. This implements all the fuunctions of ClerkOperations.
type DB struct {
	*sql.DB
}

func (c Cluster) String() string {
	return ""
}

// NewDB returns the DB object with an active database connection
func NewDB(c *mysql.DBConfig) (*DB, error) {
	return nil, nil
}

func (db *DB) CheckAvail() (bool, string, string) {
	return true, "", ""
}

// insert a cluster entry into db
func (db *DB) InsertCluster(nodes int, nodeType string, accessToken string, projectID string, prowJobID string, status string, zone string) error {
	return nil
}

// Update a cluster entry on different fields len(fields) == len(values)
func (db *DB) UpdateCluster(clusterID string, fields []string, values []string) error {
	return nil
}

// List all clusters within a certain time window
func (db *DB) ListClusters(window time.Duration) ([]Cluster, error) {
	var result []Cluster
	return result, nil

}

// DeleteCluster deletes a row from Cluster db
func (db *DB) DeleteCluster(accessToken string) error {
	return nil
}

func (db *DB) generateToken() string {
	return string(uuid.NewUUID())
}

func (db *DB) getCluster(token string) (*ClusterConfig, error) {
	// check in with token
	return &ClusterConfig{}, nil
}

// checkRowAffected number of checks row affected by a given query.
func checkRowAffected(stmt *sql.Stmt, args ...interface{}) error {
	return nil
}

// name cluster in the following format: e2e-cluster{id}-{prowJobID}
func nameCluster(clusterID string, prowJobID string) string {
	return fmt.Sprintf("e2e-cluster%s-%s", clusterID, prowJobID)
}

// insert a request entry
func (db *DB) InsertRequest(Zone string, Nodes int, NodeType string, ProwJobID string) error {
	return nil
}

// Update a request entry  len(fields) == len(values)
func (db *DB) UpdateRequest(requestID string, fields []string, values []string) error {
	return nil
}

// List requests within a time interval (use for checking after downtime to see stale requests)
func (db *DB) ListRequests(window time.Duration) ([]Request, error) {
	var result []Request
	return result, nil
}
