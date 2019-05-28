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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/knative/test-infra/tools/monitoring/mysql"
)

var dbConfig mysql.DBConfig

const (
	yamlURL              = "https://raw.githubusercontent.com/knative/test-infra/master/tools/monitoring/sample.yaml"
	dbUserSecretFile     = "/secrets/cloudsql/monitoringdb/username"
	dbPasswordSecretFile = "/secrets/cloudsql/monitoringdb/password"
)

func main() {
	var err error

	dbName := flag.String("database-name", "", "The monitoring database name")
	dbInst := flag.String("database-instance", "", "The monitoring CloudSQL instance connection name")
	flag.Parse()

	dbConfig, err = configureMonitoringDatabase(*dbName, *dbInst)
	if err != nil {
		log.Fatal(err)
	}

	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/hello", hello)
	server.HandleFunc("/test-conn", testCloudSQLConn)

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)

	yamlFile, err := ParseYaml(yamlURL)
	if err != nil {
		log.Fatalf("Cannot parse yaml: %v", err)
	}

	errorPatterns := yamlFile.CollectErrorPatterns()
	fmt.Fprintf(w, "error patterns collected from yaml:%s", errorPatterns)
}

func testCloudSQLConn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	log.Println("Testing mysql database connection.")

	err := dbConfig.TestConn()
	if err != nil {
		fmt.Fprintf(w, "Failed to ping the database %v", err)
		return
	}
	fmt.Fprintf(w, "Success\n")
}

func configureMonitoringDatabase(dbName string, dbInst string) (mysql.DBConfig, error) {
	var config mysql.DBConfig

	user, err := ioutil.ReadFile(dbUserSecretFile)
	if err != nil {
		return config, err
	}

	pass, err := ioutil.ReadFile(dbPasswordSecretFile)
	if err != nil {
		return config, err
	}

	config = mysql.DBConfig{
		Username:     string(user),
		Password:     string(pass),
		DatabaseName: dbName,
		Instance:     dbInst,
	}

	return config, nil
}
