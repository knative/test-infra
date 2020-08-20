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
	"log"
	"net/http"
	"time"
)

type Options struct {
	Region        string
	BackupRegions []string
	XXTimeout     time.Time
}

type requestInfo struct {
	minNodes int
	maxNodes int
	nodeType string
	zone     string
}

func releaseProject() bool {
	return true
}

func pollProject() string {
	return "Boskos Project"
}

func updateClerk(config requestInfo, prowid string) string {
	return "response token"
}

func handleProw(w http.ResponseWriter, req *http.Request) {
	// handle prow requests

}

func Start(o *Options) {
	log.Printf("Start running the main service with options: %v", o)

	// clustercli, _ := cluster.NewClient()
	// clustercli.Boskos.Release("", "")
	// clerkcli := clerk.NewClient()
}

// func main() {
//
// 	http.HandleFunc("/createcluster", handleProw)
// 	http.HandleFunc("/getcluster", handleProw)
// 	http.ListenAndServe(":8090", nil)
// }
