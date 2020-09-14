package main

import (
	"flag"
	"log"
	"os"

	"knative.dev/test-infra/pkg/mysql"
	"knative.dev/test-infra/tools/dkcm/mainservice"
)

func main() {
	dbName := flag.String("database-name", "dkcm", "The dkcm database name")
	dbPort := flag.String("database-port", "3307", "The dkcm database port")

	dbUserSF := flag.String("database-user", "/secrets/cloudsql/dkcmdb/username", "Database user secret file")
	dbPassSF := flag.String("database-password", "/secrets/cloudsql/dkcmdb/password", "Database password secret file")
	dbHost := flag.String("database-host", "/secrets/cloudsql/dkcmdb/host", "Database host secret file")

	boskosClientHost := flag.String("boskos-client-host", "dkcm", "Boskos client host name")

	gcpServiceAccount := flag.String("gcp-service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCP service account")

	flag.Parse()

	dbConfig, err := mysql.ConfigureDB(*dbUserSF, *dbPassSF, *dbHost, *dbPort, *dbName)
	if err != nil {
		log.Fatal(err)
	}

	if err := mainservice.Start(dbConfig, *boskosClientHost, *gcpServiceAccount); err != nil {
		log.Fatalf("Failed to start main service: %v", err)
	}
}
