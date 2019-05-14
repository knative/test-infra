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

// MySQLConfig is the configuration used to connection to a MySQL database
type MySQLConfig struct {
	DatabaseName string
	Username     string
	Password     string
	Host         string
	Port         int
}

func testConn(config MySQLConfig) error {
	conn, err := getConn(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func getConn(config MySQLConfig) (*sql.DB, error) {
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

func (c MySQLConfig) dataStoreName(databaseName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, databaseName)
}
