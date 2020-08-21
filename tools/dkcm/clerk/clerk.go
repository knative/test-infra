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

	"k8s.io/apimachinery/pkg/util/uuid"
	// _ "github.com/lib/pq"
)

type Cluster struct {
	clusterid   string
	accesstoken string
	boskosid    string
	prowid      string
	status      string
	zone        string
}

var (
	db *sql.DB
)

func generateUnique(idSize int, key string) string {
	var randomid string
	var count int
	for {
		randomid = uuid.NewUUID()
		db.QueryRow("SELECT count(*) FROM table Where $1 = $2", key, randomid).Scan(&count)
		if count == 0 {
			return randomid
		}
	}
}

func info() {
	// print the database
}

func query() (bool, string, string) {
	// check whether available cluster exists
	return true, "", ""
}

func getWithToken(token string, status chan string, errChan chan string) {
	// check in with token

}

func updateCluster(zone, prowid, boskosid string, infoChan chan string) {
	// assign cluster if available

}
