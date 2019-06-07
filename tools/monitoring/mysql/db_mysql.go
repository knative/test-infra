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

package mysql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"

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

func (c DBConfig) TestConn() error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func (c DBConfig) getConn() (*sql.DB, error) {
	conn, err := sql.Open(driverName, c.dataStoreName(c.DatabaseName))
	if err != nil {
		return nil, fmt.Errorf("could not get a connection: %v", err)
	}

	if conn.Ping() == driver.ErrBadConn {
		return nil, fmt.Errorf("could not connect to the datastore. " +
			"could be bad address, or this address is inaccessible from your host.\n")
	}

	return conn, nil
}

func (c DBConfig) dataStoreName(dbName string) string {
	var cred string
	// [username[:password]@]
	if len(c.Username) > 0 {
		cred = c.Username
		if len(c.Password) > 0 {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	return fmt.Sprintf("%sunix(%s)/%s", cred, "/cloudsql/"+c.Instance, dbName)
}

// PubsubMsgHandler adds record(s) to ErrorLogs table in database,
// after parsing build log and compares the result with config yaml
func PubsubMsgHandler(db *sql.DB, configURL, buildLogURL, jobname string, prNumber int) error {
	config, err := config.ParseYaml(configURL)
	if err != nil {
		return err
	}

	errorLogs, err := log_parser.ParseLog(buildLogURL, config.CollectErrorPatterns())
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO ErrorLogs('ErrorPattern', 'ErrorMsg', 'JobName', 'PRNumber', 'BuildLogURL', 'TimeStamp') VALUES(?,?,?,?,?,?)")
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr == nil {
			return fmt.Errorf("SQL Statement preparation failed: %v; rolled back", err)

		}
		return fmt.Errorf("SQL Statement preparation failed: %v; rollback failed: %v", err, rollbackErr)
	}

	defer stmt.Close()

	for _, errorLog := range errorLogs {
		_, err := stmt.Exec(errorLog.Pattern, errorLog.Msg, jobname, prNumber, buildLogURL, time.Now())
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr == nil {
				return fmt.Errorf("SQL Statement execution failed: %v; rolled back", err)
			}
			return fmt.Errorf("SQL Statement execution failed: %v; rollback failed: %v", err, rollbackErr)
		}
	}

	return tx.Commit()
}

// CheckAlertCondition checks whether the given error pattern meets
// the alert condition specified in config
func CheckAlertCondition(errorPattern string, config *config.SelectedConfig, db *sql.DB) (bool, error) {
	// the timestamp we want to start collecting logs
	startTime := time.Now().Add(time.Minute * time.Duration(config.Period))

	_, err := db.Query(`
		CREATE VIEW Matched AS
  	SELECT Jobname, PrNumber 
  	FROM ErrorLogs
  	WHERE ErrorPattern=? and TimeStamp > ?`,
		errorPattern, startTime)

	if err != nil {
		return false, err
	}

	var nMatches, nJobs, nPRs int

	row := db.QueryRow(`
		SELECT 
    	COUNT(*),
    	COUNT (DISTINCT Jobname),
    	COUNT (DISTINCT PrNumber)
		FROM Matched;`)

	if err = row.Scan(&nMatches, nJobs, nPRs); err != nil {
		return false, err
	}

	return nMatches >= config.Occurrences && nJobs >= config.JobsAffected && nPRs >= config.PrsAffected, nil
}
