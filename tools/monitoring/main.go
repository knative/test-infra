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
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/knative/test-infra/shared/mysql"
	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/mail"
	"github.com/knative/test-infra/tools/monitoring/prowapi"
	"github.com/knative/test-infra/tools/monitoring/subscriber"
)

var (
	dbConfig   *mysql.DBConfig
	mailConfig *mail.Config
	client     *subscriber.Client

	alertEmailRecipients = []string{"knative-productivity-oncall@googlegroups.com"}
)

const (
	projectID = "knative-tests"

	yamlURL = "https://raw.githubusercontent.com/knative/test-infra/master/tools/monitoring/sample.yaml"
	subName = "test-infra-monitoring-sub"
)

func main() {
	var err error

	dbName := flag.String("database-name", "monitoring", "The monitoring database name")
	dbInst := flag.String("database-instance", "knative-tests:us-central1:knative-monitoring", "The monitoring CloudSQL instance connection name")

	dbUserSF := flag.String("database-user", "/secrets/cloudsql/monitoringdb/username", "Database user secret file")
	dbPassSF := flag.String("database-password", "/secrets/cloudsql/monitoringdb/password", "Database password secret file")
	mailAddrSF := flag.String("sender-email", "/secrets/sender-email/mail", "Alert sender email address file")
	mailPassSF := flag.String("sender-password", "/secrets/sender-email/password", "Alert sender email password file")

	flag.Parse()

	dbConfig, err = mysql.ConfigureDB(*dbUserSF, *dbPassSF, *dbName, *dbInst)
	if err != nil {
		log.Fatal(err)
	}

	mailConfig, err = mail.NewMailConfig(*mailAddrSF, *mailPassSF)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err = subscriber.NewSubscriberClient(ctx, projectID, subName)
	if err != nil {
		log.Fatalf("Failed to initialize the subscriber %+v", err)
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
	server.HandleFunc("/send-mail", sendTestEmail)
	server.HandleFunc("/test-sub", testSubscriber)

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

// hello tests the as much completed steps in the entire monitoring workflow as possible
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)

	yamlFile, err := config.ParseYaml(yamlURL)
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

func sendTestEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	log.Println("Sending test email")

	err := mailConfig.Send(
		alertEmailRecipients,
		"Test Subject",
		"Test Content",
	)
	if err != nil {
		fmt.Fprintf(w, "Failed to send email %v", err)
		return
	}

	fmt.Fprintln(w, "Sent the Email")
}

func testSubscriber(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	log.Println("Start listening to messages")

	go func() {
		err := client.ReceiveMessageAckAll(context.Background(), func(rmsg *prowapi.ReportMessage) {
			log.Printf("Report Message: %+v\n", rmsg)
		})
		if err != nil {
			log.Printf("Failed to retrieve messages due to %v", err)
		}
	}()
}
