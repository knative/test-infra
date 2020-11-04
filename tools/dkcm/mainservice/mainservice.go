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

package mainservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"knative.dev/test-infra/pkg/clustermanager/e2e-tests/boskos"
	"knative.dev/test-infra/pkg/clustermanager/kubetest2"
	"knative.dev/test-infra/pkg/mysql"
	"knative.dev/test-infra/tools/dkcm/clerk"
)

var (
	// channel serves as a lock for go routine
	chanLock             = make(chan struct{})
	boskosClient         *boskos.Client
	dbClient             *clerk.DBClient
	serviceAccount       string
	DefaultClusterParams = clerk.ClusterParams{Zone: DefaultZone, Nodes: DefaultNodesCount, NodeType: DefaultNodeType}
)

// Response to Prow
type ServiceResponse struct {
	IsReady     bool            `json:"isReady"`
	Message     string          `json:"message"`
	ClusterInfo *clerk.Response `json:"clusterInfo"`
}

func Start(dbConfig *mysql.DBConfig, boskosClientHost, gcpServiceAccount string) error {
	var err error
	boskosClient, err = boskos.NewClient(boskosClientHost, "", "")
	if err != nil {
		return fmt.Errorf("failed to create Boskos client: %w", err)
	}
	dbClient, err = clerk.NewDB(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create Clerk client: %w", err)
	}
	serviceAccount = gcpServiceAccount
	server := http.NewServeMux()
	server.HandleFunc("/request-cluster", handleNewClusterRequest)
	server.HandleFunc("/get-cluster", handleGetCluster)
	server.HandleFunc("/clean-cluster", handleCleanCluster)
	// use PORT environment variable, or default to 8080
	port := DefaultPort
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	// start the web server on port and accept requests
	log.Printf("Server listening on port %q", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), server)
}

// handle cleaning cluster request after usage
func handleCleanCluster(w http.ResponseWriter, req *http.Request) {
	// add project name
	token := req.URL.Query().Get("token")
	r, err := dbClient.GetRequest(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error getting the request with the token: %v, please try again", err), http.StatusForbidden)
		return
	}
	c, err := dbClient.GetCluster(r.ClusterID)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error getting the cluster with the token: %v, please try again", err), http.StatusForbidden)
		return
	}
	err = dbClient.DeleteCluster(r.ClusterID)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error deleting the cluster with the token: %v, please try again", err), http.StatusForbidden)
		return
	}
	err = boskosClient.ReleaseGKEProject(c.ProjectID)
	if err != nil {
		http.Error(w, "there is an error releasing Boskos's project. Please try again.", http.StatusInternalServerError)
		return
	}
}

// check the pool capacity and create clusters if necessary
func checkPoolCap(cp *clerk.ClusterParams) {
	chanLock <- struct{}{}
	numAvail := dbClient.CheckNumStatus(cp, Ready)
	numWIP := dbClient.CheckNumStatus(cp, WIP)
	diff := DefaultOverProvision - numAvail - numWIP
	var wg sync.WaitGroup
	// create cluster if not meeting overprovisioning criteria
	for i := int64(0); i < diff; i++ {
		wg.Add(1)
		log.Printf("Creating a new cluster: %v", cp)
		go CreateCluster(cp, &wg)
	}
	wg.Wait()
	<-chanLock
}

// assign clusters if available upon request
func CreateCluster(cp *clerk.ClusterParams, wg *sync.WaitGroup) {
	project, err := boskosClient.AcquireGKEProject(boskos.GKEProjectResource)
	if err != nil {
		log.Printf("Failed to acquire a project from boskos: %v", err)
		return
	}
	projectName := project.Name
	c := clerk.NewCluster(clerk.AddProjectID(projectName))
	c.ClusterParams = cp
	clusterID, err := dbClient.InsertCluster(c)
	wg.Done()
	if err != nil {
		log.Printf("Failed to insert a new Cluster entry: %v", err)
		return
	}
	if err := kubetest2.Run(&kubetest2.Options{}, &kubetest2.GKEClusterConfig{
		GCPServiceAccount: serviceAccount,
		GCPProjectID:      projectName,
		Name:              DefaultClusterName,
		Region:            cp.Zone,
		Machine:           cp.NodeType,
		MinNodes:          int(cp.Nodes),
		MaxNodes:          int(cp.Nodes),
		Network:           DefaultNetworkName,
		Environment:       "prod",
		Version:           "latest",
		Scopes:            "cloud-platform",
	}); err != nil {
		if err := dbClient.UpdateCluster(clusterID, clerk.UpdateStringField(Status, Fail)); err != nil {
			log.Printf("Failed to insert delete new Cluster entry: %v", err)
			return
		}
		err = boskosClient.ReleaseGKEProject(c.ProjectID)
		if err != nil {
			log.Printf("Failed to release Boskos Project: %v", err)
			return
		}
		return
	}
	if err := dbClient.UpdateCluster(clusterID, clerk.UpdateStringField(Status, Ready)); err != nil {
		log.Printf("Failed to insert a new Cluster entry: %v", err)
		return
	}
}

// assign clusters if available upon request
func AssignCluster(token string, w http.ResponseWriter) {
	r, err := dbClient.GetRequest(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error getting the request with the token: %v, please try again", err), http.StatusForbidden)
		return
	}
	// check if the Prow job has enough priority to get an existing cluster
	ranking := dbClient.PriorityRanking(r)
	numAvail := dbClient.CheckNumStatus(r.ClusterParams, Ready)
	available, clusterID := dbClient.CheckAvail(r.ClusterParams)
	var serviceResponse *ServiceResponse
	if available && ranking <= numAvail {
		response, err := dbClient.GetCluster(clusterID)
		if err != nil {
			http.Error(w, fmt.Sprintf("there is an error getting available clusters: %v, please try again", err), http.StatusInternalServerError)
			return
		}
		dbClient.UpdateRequest(r.ID, clerk.UpdateNumField(ClusterID, clusterID))
		dbClient.UpdateCluster(clusterID, clerk.UpdateStringField(Status, InUse))
		serviceResponse = &ServiceResponse{IsReady: true, Message: "Your cluster is ready!", ClusterInfo: response}
	} else {
		serviceResponse = &ServiceResponse{IsReady: false, Message: "Your cluster isn't ready yet! Please check back later."}
	}
	responseJson, err := json.Marshal(serviceResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error getting parsing response: %v, please try again", err), http.StatusInternalServerError)
		return
	}
	w.Write(responseJson)
}

// handle new cluster request
func handleNewClusterRequest(w http.ResponseWriter, req *http.Request) {
	prowJobID := req.PostFormValue("prowjobid")
	nodesCount, err := strconv.Atoi(req.PostFormValue("nodes"))
	if err != nil || nodesCount <= 0 {
		nodesCount = DefaultNodesCount
	}
	nodesType := req.PostFormValue("nodeType")
	if nodesType == "" {
		nodesType = DefaultNodeType
	}
	zone := req.PostFormValue("zone")
	if zone == "" {
		zone = DefaultZone
	}
	cp := clerk.NewClusterParams(clerk.AddZone(zone), clerk.AddNodes(int64(nodesCount)), clerk.AddNodeType(nodesType))
	r := clerk.NewRequest(clerk.AddProwJobID(prowJobID), clerk.AddRequestTime(time.Now()))
	r.ClusterParams = cp
	accessToken, err := dbClient.InsertRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("there is an error creating new request: %v. Please try again.", err), http.StatusInternalServerError)
		return
	}
	go checkPoolCap(cp)
	w.Write([]byte(accessToken))
}

// handle get cluster request
func handleGetCluster(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	AssignCluster(token, w)
}
