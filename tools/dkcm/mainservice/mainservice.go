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
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"knative.dev/test-infra/pkg/clustermanager/e2e-tests/boskos"
	"knative.dev/test-infra/pkg/mysql"
	"knative.dev/test-infra/tools/dkcm/clerk"
)

var (
	boskosClient         *boskos.Client
	dbConfig             *mysql.DBConfig
	dbClient             *clerk.DBClient
	DefaultClusterParams = clerk.ClusterParams{Zone: DefaultZone, Nodes: DefaultNodesCount, NodeType: DefaultNodeType}
)

// Reponse to Prow
type ServiceResponse struct {
	isReady  bool
	message  string
	response *clerk.Response
}

func Start() {
	var err error
	boskosClient, err = boskos.NewClient("", "", "")
	if err != nil {
		log.Fatalf("Failed to create Boskos Client %v", err)
	}
	server := http.NewServeMux()
	server.HandleFunc("/request-cluster", handleNewClusterRequest)
	server.HandleFunc("/get-cluster", handleGetCluster)
	server.HandleFunc("/clean-cluster", handleCleanCluster)
	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

// handle cleaning cluster request after usage
func handleCleanCluster(w http.ResponseWriter, r *http.Request) {
	err := boskosClient.ReleaseGKEProject("")
	if err != nil {
		http.Error(w, "there is an error releasing Boskos's project. Please try again.", http.StatusInternalServerError)
		return
	}
}

// check the pool capacity
func checkPoolCap(cp *clerk.ClusterParams) {
	numAvail := dbClient.CheckNumStatus(cp, "Ready")
	numWIP := dbClient.CheckNumStatus(cp, "WIP")
	if numAvail+numWIP < DefaultOverProvision {
		// create cluster if not meeting overprovisioning criteria
		CreateCluster(cp)
	}
}

// assign clusters if available upon request
func CreateCluster(cp *clerk.ClusterParams) error {
	// TODO: kubetest2 create cluster
	name, err := GetProject()
	if err != nil {
		return err
	}
	c := clerk.NewCluster(clerk.AddProjectID(name))
	c.ClusterParams = cp
	err = dbClient.InsertCluster(c)
	return err
}

// assign clusters if available upon request
func AssignCluster(token string, w http.ResponseWriter) {
	r, err := dbClient.GetRequest(token)
	if err != nil {
		http.Error(w, "there is an error getting the request with the token. Please try again.", http.StatusForbidden)
		return
	}
	// check if the Prow job has enough priority to get an existing cluster
	ranking := dbClient.PriorityRanking(r)
	numAvail := dbClient.CheckNumStatus(r.ClusterParams, "Ready")
	available, clusterID := dbClient.CheckAvail(r.ClusterParams)
	var serviceResponse *ServiceResponse
	if available && ranking <= numAvail {
		response, err := dbClient.GetCluster(clusterID)
		if err != nil {
			http.Error(w, "there is an error getting available clusters. Please try again.", http.StatusInternalServerError)
			return
		}
		dbClient.UpdateRequest(r.RequestID, clerk.UpdateNumField("ClusterID", clusterID))
		checkPoolCap(&DefaultClusterParams)
		serviceResponse = &ServiceResponse{isReady: true, message: "Your cluster is ready!", response: response}
	} else {
		serviceResponse = &ServiceResponse{isReady: false, message: "Your cluster isn't ready yet! Please check back later."}
	}
	responseJson, err := json.Marshal(serviceResponse)
	if err != nil {
		http.Error(w, "there is an error getting parsing response. Please try again.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJson))
}

// main go
func LoadDB() {

}

func GetProject() (string, error) {
	resource, err := boskosClient.AcquireGKEProject(boskos.GKEProjectResource)
	if err != nil {
		return "", err
	}
	return resource.Name, err
}

// handle new cluster request
func handleNewClusterRequest(w http.ResponseWriter, req *http.Request) {
	prowJobID := req.PostFormValue("prowjobid")
	nodesCount, err := strconv.Atoi(req.PostFormValue("nodes"))
	if err != nil {
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
	if err != nil || nodesCount <= 0 {
		http.Error(w, "there is an error creating new request. Please try again.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(accessToken))
	checkPoolCap(cp)
}

// handle get cluster request
func handleGetCluster(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	AssignCluster(token, w)
}

// run timeout check
func runTimeOut() {
	for {
		dbClient.ClearTimeOut(DefaultTimeOut)
		time.Sleep(2 * time.Second)
	}
}
