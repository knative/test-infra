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
	"database/sql"
	"database/sql/driver"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const driverName = "mysql"

// DBConfig is the configuration used to connection to database
type DBConfig struct {
	Username     string
	Password     string
	Instance     string
	DatabaseName string
}

func getConn(config DBConfig) (*sql.DB, error) {
	conn, err := sql.Open(driverName, config.dataStoreName(config.DatabaseName))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}

	if conn.Ping() == driver.ErrBadConn {
		return nil, fmt.Errorf("mysql: could not connect to the datastore. " +
			"could be bad address, or this address is not whitelisted for access.\n")
	}

	return conn, nil
}

func (c DBConfig) dataStoreName(dbName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	return fmt.Sprintf("%sunix(%s)/%s", cred, "/cloudsql/"+c.Instance, dbName)
}
