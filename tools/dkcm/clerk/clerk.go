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
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"
	"knative.dev/test-infra/pkg/mysql"
)

// ClusterParams is a struct for common cluster parameters shared by both Cluster and Request
type ClusterParams struct {
	Zone     string
	Nodes    int64
	NodeType string
}

// Cluster stores a row in the "Cluster" db table
type Cluster struct {
	*ClusterParams
	ClusterID int64
	ProjectID string
	Status    string
}

// Request stores a request for the "Request" table
type Request struct {
	*ClusterParams
	accessToken string
	requestTime time.Time
	ProwJobID   string
	ClusterID   int64
}

// Response is the struct to return(end product) to Prow once the cluster is available
type Response struct {
	ClusterName string
	ProjectID   string
	Zone        string
}

type ClerkOperations interface {
	//check cluster available, if available, return cluster id and access token
	CheckAvail(cp ClusterParams) (bool, int64)
	// get with cluster id stored in the Cluster database
	GetCluster(clusterID int64) (Response, error)
	// delete a cluster entry
	DeleteCluster(clusterID int64) error
	// Insert a cluster entry
	InsertCluster(c *Cluster) error
	// Update a cluster entry
	UpdateCluster(clusterID int64, opts ...UpdateOption) error
	// List clutsers (use for checking after downtime to see stale clusters)
	ListClusters() ([]Cluster, error)
	// Generate Unique access token for Prow
	generateToken() string
	// Insert a request entry
	InsertRequest(r *Request) error
	// Update a request entry
	UpdateRequest(requestID int64, opts ...UpdateOption) error
	// List requests within a time interval (use for checking after downtime to see stale requests)
	ListRequests(window time.Duration) ([]Request, error)
}

// compatible with both Row and Rows and unit test friendly
type scannable interface {
	Scan(dest ...interface{}) error
}

// DB holds an active database connection. This implements all the fuunctions of ClerkOperations.
type DB struct {
	*sql.DB
}

func (c Cluster) String() string {
	return fmt.Sprintf("Cluster Info: (ProjectID: %s, NodesCount: %d, NodeType: %s, Status: %s, Zone: %s)",
		c.ProjectID, c.Nodes, c.NodeType, c.Status, c.Zone)
}

func (r Request) String() string {
	return fmt.Sprintf("Request Info: (RequestTime: %v, NodesCount: %v, NodeType: %s, ProwJobID: %s, ClusterID: %d,Zone: %s)",
		r.requestTime, r.Nodes, r.NodeType, r.ProwJobID, r.ClusterID, r.Zone)
}

// generate a list of strings that could be used in a query for ClusterParams
func (cp *ClusterParams) generateParamsConditions(opts ...QueryClusterParamsOption) []string {
	var fieldStatements []string
	for _, opt := range opts {
		fieldStatements = append(fieldStatements, opt(cp))
	}
	return fieldStatements
}

func generateAND(fieldStatements []string) string {
	return strings.Join(fieldStatements, " AND ")
}

// Populate fields of Cluster
func populateCluster(sc scannable) (*Cluster, error) {
	c := &Cluster{}
	err := sc.Scan(&c.ClusterID, &c.ProjectID, &c.Nodes, &c.NodeType, &c.Status, &c.Zone)
	return c, err
}

// Populate fields of Request
func populateRequest(sc scannable) (*Request, error) {
	r := &Request{}
	err := sc.Scan(&r.accessToken, &r.requestTime, &r.Nodes, &r.NodeType, &r.ProwJobID, &r.ClusterID, &r.Zone)
	return r, err
}

// NewDB returns the DB object with an active database connection
func NewDB(c *mysql.DBConfig) (*DB, error) {
	db, err := c.Connect()
	return &DB{db}, err
}

// check available clusters in Cluster db and return zone and projectID that would be used by ProwJob
func (db *DB) CheckAvail(cp *ClusterParams) (bool, int64) {
	// getting query string ready
	conditions := generateAND(cp.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType()))
	queryString := fmt.Sprintf("SELECT * FROM Cluster WHERE Status = 'Ready' AND %v", conditions)
	// check whether available cluster exists
	row := db.QueryRow(queryString)
	c, err := populateCluster(row)
	// no available cluster is found
	if err != nil {
		return false, -1
	}
	return true, c.ClusterID

}

// insert a cluster entry into db
func (db *DB) InsertCluster(c *Cluster) error {
	stmt, err := db.Prepare(`INSERT INTO Cluster(Nodes, NodeType, Zone, ProjectID, Status)
							VALUES (?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Nodes, c.NodeType, c.Zone, c.ProjectID, c.Status)
	return err
}

func updateQueryString(dbName string, id int64, opts ...UpdateOption) string {
	var fieldStatements []string
	for _, opt := range opts {
		fieldStatements = append(fieldStatements, opt())
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE ID = %d", dbName, strings.Join(fieldStatements, ","), id)
}

// Update a cluster entry on different fields len(fields) == len(values)
func (db *DB) UpdateCluster(clusterID int64, opts ...UpdateOption) error {
	queryString := updateQueryString("Cluster", clusterID, opts...)
	_, err := db.Exec(queryString)
	return err
}

// List all clusters
func (db *DB) ListClusters() ([]Cluster, error) {
	var result []Cluster
	rows, err := db.Query(`
	SELECT ID, Nodes, NodeType, ProjectID, Status, Zone
	FROM Cluster`)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		c, err := populateCluster(rows)
		if err != nil {
			return result, err
		}
		result = append(result, *c)
	}
	return result, nil

}

// DeleteCluster deletes a row from Cluster db
func (db *DB) DeleteCluster(clusterID int64) error {
	// make sure the cluster to deleted exist
	_, err := db.GetCluster(clusterID)
	if err != nil {
		return err
	}
	delString := "DELETE FROM Cluster WHERE ID = ?"
	stmt, err := db.Prepare(delString)
	if err == nil {
		// check exactly one row is deleted
		err = checkAffected(stmt, 1, clusterID)
		defer stmt.Close()
	}
	return err
}

func (db *DB) generateToken() string {
	return string(uuid.NewUUID())
}

// assumption: get cluster only when it is ready, the clusterid that shows up on request db
func (db *DB) GetCluster(clusterID int64) (*Response, error) {
	// check in with token
	queryString := "SELECT * FROM Cluster WHERE ID = ?"
	row := db.QueryRow(queryString, clusterID)
	c, err := populateCluster(row)
	if err != nil {
		return &Response{}, err
	}
	clusterName := nameCluster(clusterID)
	cc := &Response{ClusterName: clusterName, ProjectID: c.ProjectID, Zone: c.Zone}
	return cc, err

}

// checkRowAffected expects a certain number of rows in the db to be affected.
func checkAffected(stmt *sql.Stmt, numRows int64, args ...interface{}) error {
	res, err := stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("Statement not executable: %v ", err)
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Fail to get rows affected: %v ", err)
	} else if rowsAffected != numRows {
		return fmt.Errorf("Expected %d row affected, got %d", numRows, rowsAffected)
	}
	return nil
}

// name cluster in the following format: e2e-cluster{id}-{prowJobID}
func nameCluster(clusterID int64) string {
	return fmt.Sprintf("e2e-cluster%v", clusterID)
}

// insert a request entry
func (db *DB) InsertRequest(r *Request) error {
	stmt, err := db.Prepare(`INSERT INTO Requests(AccessToken, RequestTime, ProwJobID, Nodes, NodeType, Zone)
							VALUES (?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	accessToken := db.generateToken()
	_, err = stmt.Exec(accessToken, r.requestTime, r.ProwJobID, r.Nodes, r.NodeType, r.Zone)
	return err
}

// Update a request entry  len(fields) == len(values)
func (db *DB) UpdateRequest(requestID int64, opts ...UpdateOption) error {
	queryString := updateQueryString("Request", requestID, opts...)
	_, err := db.Exec(queryString)
	return err
}

// List requests within a time interval (use for checking after downtime to see stale requests)
func (db *DB) ListRequests(window time.Duration) ([]Request, error) {
	var result []Request
	startTime := time.Now().Add(-1 * window)
	rows, err := db.Query(`
	SELECT ID, RequestTime, Nodes, NodeType, ProwJobID, Zone, ClusterID
	FROM Request
	WHERE RequestTime > ?`, startTime)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		r, err := populateRequest(rows)
		if err != nil {
			return result, err
		}
		result = append(result, *r)
	}
	return result, nil
}
