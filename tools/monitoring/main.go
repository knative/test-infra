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
	"log"
	"net/http"
	"os"

	"io/ioutil"
)

var dbConfig DBConfig

func main() {
	var err error

	fDbUser := flag.String("monitoring-database-user", "",
		"Text file containing the user name of to connect to the monitoring database")
	fDbPass := flag.String("monitoring-database-password", "",
		"Text file containing the password to connect to the monitoring database")
	flag.Parse()

	dbConfig, err = configureMonitoringDatabase(*fDbUser, *fDbPass)
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
}

func testCloudSQLConn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	log.Println("Testing mysql database connection.")

	conn, err := getConn(dbConfig)
	if err != nil {
		fmt.Fprintf(w, "Failed to ping database: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(w, "Success\n")
}

func configureMonitoringDatabase(fUser string, fPass string) (DBConfig, error) {
	var config DBConfig

	user, err := ioutil.ReadFile(fUser)
	if err != nil {
		return config, err
	}

	pass, err := ioutil.ReadFile(fPass)
	if err != nil {
		return config, err
	}

	config = DBConfig{
		Username:     string(user),
		Password:     string(pass),
		DatabaseName: os.Getenv("MONITORING_DB_NAME"),
		Instance:     os.Getenv("MONITORING_DB_INSTANCE"),
	}

	return config, nil
}
