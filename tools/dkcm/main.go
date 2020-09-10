package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/spf13/pflag"

	"knative.dev/test-infra/pkg/mysql"
	"knative.dev/test-infra/tools/dkcm/mainservice"
)

func initFlags() *mainservice.Options {
	var o mainservice.Options
	var regions []string
	pflag.StringSliceVar(&regions, "region", []string{}, "")
	pflag.Parse()
	if len(regions) > 0 {
		o.Region = regions[0]
		if len(regions) > 1 {
			o.BackupRegions = regions[1:]
		}
	} else {
		o.Region = mainservice.DefaultZone
	}
	timeOut := flag.String("timeout", "", "The mainservice database name")
	dbName := flag.String("database-name", "mainservice", "The mainservice database name")
	dbPort := flag.String("database-port", "3306", "The mainservice database port")
	dbUserSF := flag.String("database-user", "/secrets/cloudsql/clerkdb/username", "Database user secret file")
	dbPassSF := flag.String("database-password", "/secrets/cloudsql/clerkdb/password", "Database password secret file")
	dbHost := flag.String("database-host", "/secrets/cloudsql/clerkdb/host", "Database host secret file")
	flag.Parse()
	dbConfig, err := mysql.ConfigureDB(*dbUserSF, *dbPassSF, *dbHost, *dbPort, *dbName)
	if err != nil {
		log.Fatal(err)
	}
	if *timeOut == "" {
		o.Timeout = mainservice.DefaultTimeOut
	} else {
		userTimeOut, err := strconv.Atoi(*timeOut)
		if err != nil {
			log.Fatal(err)
		}
		o.Timeout = time.Duration(userTimeOut)
	}
	o.DBConfig = dbConfig
	return &o
}

func main() {
	o := initFlags()
	mainservice.Start(o)
}
