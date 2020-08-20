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
	"math/rand"
	"time"
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	db        *sql.DB
	randomGen *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

/*
func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=newPassword")
	if err != nil {
		log.Fatal(err)
	}
	// check connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

} */

func generateID(idSize int) string {
	bytes := make([]byte, idSize)
	for i := range bytes {
		bytes[i] = charset[randomGen.Intn(len(charset))]
	}
	return string(bytes)
}

func generateUnique(idSize int, key string) string {
	var randomid string
	var count int
	for {
		randomid = generateID(idSize)
		db.QueryRow("SELECT count(*) FROM table Where $1 = $2", key, randomid).Scan(&count)
		if count == 0 {
			return randomid
		}
	}
}

func clerkInfo() {
	// print the database
}

func clerkQuery() (bool, string, string) {
	// check whether available cluster exists
	return true, "", ""
}

func getWithToken(token string, status chan string, errc chan string) {
	// check in with token

}

func updateCluster(zone string, prowid string, boskosid string, infoChan chan string) {
	// assign cluster if available

}
