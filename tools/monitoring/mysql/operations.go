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
	"fmt"
	"time"

	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"
)

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
	defer stmt.Close()

	if err != nil {
		return rollbackTx(tx, err)
	}

	for _, errorLog := range errorLogs {
		if _, err := stmt.Exec(errorLog.Pattern, errorLog.Msg, jobname, prNumber, buildLogURL, time.Now()); err != nil {
			return rollbackTx(tx, err)
		}
	}

	return tx.Commit()
}

// rollbackTx will try to rollback the transaction and return an error message accordingly
func rollbackTx(tx *sql.Tx, err error) error {
	rollbackErr := tx.Rollback()
	if rollbackErr == nil {
		return fmt.Errorf("Statement execution failed: %v; rolled back", err)
	}
	return fmt.Errorf("Statement execution failed: %v; rollback failed: %v", err, rollbackErr)
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

	if err = row.Scan(&nMatches, &nJobs, &nPRs); err != nil {
		return false, err
	}

	return nMatches >= config.Occurrences && nJobs >= config.JobsAffected && nPRs >= config.PrsAffected, nil
}
