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
	"log"
	"math"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"

	"knative.dev/test-infra/pkg/mysql"
)

type Operations interface {
	// check cluster available, if available, return cluster id and access token
	CheckAvail(cp ClusterParams) (bool, int64)
	// check number of clusters of a status specific configurations
	CheckNumStatus(cp *ClusterParams, status string) int64
	// get with cluster id stored in the Cluster database
	GetCluster(clusterID int64) (Response, error)
	// delete a cluster entry
	DeleteCluster(clusterID int64) error
	// Insert a cluster entry
	InsertCluster(c *Cluster) (int64, error)
	// Update a cluster entry
	UpdateCluster(clusterID int64, opts ...UpdateOption) error
	// List clutsers (use for checking after downtime to see stale clusters)
	ListClusters() ([]Cluster, error)
	// get with accessToken stored in the Request database
	GetRequest(accessToken string) (Request, error)
	// Generate Unique access token for Prow
	generateToken() string
	// Insert a request entry
	InsertRequest(r *Request) error
	// Update a request entry
	UpdateRequest(requestID int64, opts ...UpdateOption) error
	// List requests within a time interval (use for checking after downtime to see stale requests)
	ListRequests(window time.Duration) ([]Request, error)
	// Rank priority based on time
	PriorityRanking(r *Request) int64
	// Detect Timeout for requests and mark clusterID as -1 for timeout
	ClearTimeOut(timeOut time.Duration) error
}

// compatible with both Row and Rows and unit test friendly
type scannable interface {
	Scan(dest ...interface{}) error
}

// DBClient holds an active database connection. This implements all the fuunctions of Operations.
type DBClient struct {
	*sql.DB
}

func (r Request) String() string {
	return fmt.Sprintf("Request Info: (RequestTime: %v, NodesCount: %v, NodeType: %s, ProwJobID: %s, Zone: %s)",
		r.requestTime, r.Nodes, r.NodeType, r.ProwJobID, r.Zone)
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
	c := &Cluster{ClusterParams: &ClusterParams{}}
	err := sc.Scan(&c.ID, &c.ProjectID, &c.Status, &c.Zone, &c.Nodes, &c.NodeType)
	return c, err
}

// Populate fields of Request
func populateRequest(sc scannable) (*Request, error) {
	r := &Request{ClusterParams: &ClusterParams{}}
	err := sc.Scan(&r.ID, &r.accessToken, &r.requestTime, &r.Zone, &r.Nodes, &r.NodeType, &r.ProwJobID, &r.ClusterID)
	return r, err
}

// NewDB returns the DB object with an active database connection
func NewDB(c *mysql.DBConfig) (*DBClient, error) {
	db, err := c.Connect()
	return &DBClient{db}, err
}

// check available clusters in Cluster db and return zone and projectID that would be used by ProwJob
func (db *DBClient) CheckAvail(cp *ClusterParams) (bool, int64) {
	// getting query string ready
	conditions := generateAND(cp.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType()))
	queryString := fmt.Sprintf("SELECT * FROM Clusters WHERE Status = 'Ready' AND %s", conditions)
	// check whether available cluster exists
	row := db.QueryRow(queryString)
	c, err := populateCluster(row)
	// no available cluster is found
	if err != nil {
		return false, -1
	}
	return true, c.ID
}

func (db *DBClient) CheckNumStatus(cp *ClusterParams, status string) int64 {
	var count int64
	// getting query string ready
	conditions := generateAND(cp.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType()))
	queryString := fmt.Sprintf("SELECT COUNT(*) FROM Clusters WHERE Status = '%s' AND %s", status, conditions)
	log.Printf("query string is: %s", queryString)
	err := db.QueryRow(queryString).Scan(&count)
	if err != nil {
		log.Printf("error is: %v", err)
		return 0
	}
	return count
}

// insert a cluster entry into db
func (db *DBClient) InsertCluster(c *Cluster) (int64, error) {
	// query only rows that haven't been assigned a cluster and of the same config
	stmt, err := db.Prepare(`INSERT INTO Clusters(Nodes, NodeType, Zone, ProjectID, Status)
							VALUES (?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Nodes, c.NodeType, c.Zone, c.ProjectID, c.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateQueryString(dbName string, id int64, opts ...UpdateOption) string {
	var fieldStatements []string
	for _, opt := range opts {
		fieldStatements = append(fieldStatements, opt())
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE ID = %d", dbName, strings.Join(fieldStatements, ","), id)
}

// Update a cluster entry on different fields
func (db *DBClient) UpdateCluster(clusterID int64, opts ...UpdateOption) error {
	queryString := updateQueryString("Clusters", clusterID, opts...)
	_, err := db.Exec(queryString)
	return err
}

// List all clusters
func (db *DBClient) ListClusters() ([]Cluster, error) {
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
func (db *DBClient) DeleteCluster(clusterID int64) error {
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

func (db *DBClient) generateToken() string {
	return string(uuid.NewUUID())
}

// assumption: get cluster only when it is ready, the clusterid that shows up on request db
func (db *DBClient) GetCluster(clusterID int64) (*Response, error) {
	// check in with token
	queryString := "SELECT * FROM Clusters WHERE ID = ?"
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
		return fmt.Errorf("statement not executable: %v ", err)
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("fail to get rows affected: %v ", err)
	} else if rowsAffected != numRows {
		return fmt.Errorf("expected %d row affected, got %d", numRows, rowsAffected)
	}
	return nil
}

// name cluster in the following format: e2e-cluster{id}-{prowJobID}
func nameCluster(clusterID int64) string {
	return fmt.Sprintf("e2e-cluster%v", clusterID)
}

// insert a request entry
func (db *DBClient) InsertRequest(r *Request) (string, error) {
	stmt, err := db.Prepare(`INSERT INTO Requests(AccessToken, RequestTime, ProwJobID, Nodes, NodeType, Zone)
							VALUES (?,?,?,?,?,?)`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	accessToken := db.generateToken()
	_, err = stmt.Exec(accessToken, r.requestTime, r.ProwJobID, r.Nodes, r.NodeType, r.Zone)
	return accessToken, err
}

// Update a request entry
func (db *DBClient) UpdateRequest(requestID int64, opts ...UpdateOption) error {
	queryString := updateQueryString("Requests", requestID, opts...)
	_, err := db.Exec(queryString)
	return err
}

// Get a request row in Request db by accessToken
func (db *DBClient) GetRequest(accessToken string) (*Request, error) {
	// check in with token
	queryString := "SELECT * FROM Requests WHERE AccessToken = ?"
	row := db.QueryRow(queryString, accessToken)
	r, err := populateRequest(row)
	if err != nil {
		return &Request{}, err
	}
	return r, nil
}

// List requests within a time interval (use for checking after downtime to see stale requests)
func (db *DBClient) ListRequests(window time.Duration) ([]Request, error) {
	var result []Request
	startTime := time.Now().Add(-1 * window)
	rows, err := db.Query(`
	SELECT *
	FROM Requests
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

// rank request priority so that available clusters are always assigned to the requests that come first
func (db *DBClient) PriorityRanking(r *Request) int64 {
	var rank int64
	// getting query string ready
	conditions := generateAND(r.ClusterParams.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType()))
	// query only rows that haven't been assigned a cluster and of the same config
	queryString := fmt.Sprintf("Select Rnk From (SELECT ID, RANK() OVER (ORDER BY RequestTime) Rnk FROM Requests WHERE ClusterID = 0 AND %s) Ranking WHERE Ranking.ID = %d", conditions, r.ID)
	log.Printf("query string is: %s", queryString)
	err := db.QueryRow(queryString).Scan(&rank)
	if err != nil {
		log.Printf("Got an error in the query: %v", err)
		return math.MaxInt64
	}
	return rank
}

// disable request due to timeout
func (db *DBClient) ClearTimeOut(timeOut time.Duration) error {
	startTime := time.Now().Add(-1 * timeOut)
	queryString := "UPDATE Requests SET ClusterID = -1 WHERE RequestTime <= ? AND ClusterID = 0"
	_, err := db.Exec(queryString, startTime)
	return err
}
