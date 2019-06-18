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

package log_parser

import (
	"database/sql"
	"time"

	"github.com/knative/test-infra/shared/mysql"
	"github.com/knative/test-infra/tools/monitoring/config"
)

const (
	logInsertStmt = `
	INSERT INTO ErrorLogs (
		ErrorPattern, ErrorMsg, JobName, PRNumber, BuildLogURL, TimeStamp
		) VALUES (?,?,?,?,?,?)`
)

// ErrorLog stores a row in the "ErrorLogs" db table
// Table schema: github.com/knative/test-infra/tools/monitoring/mysql/schema.sql
type ErrorLog struct {
	Pattern     string
	Msg         string
	JobName     string
	PRNumber    int
	BuildLogURL string
	TimeStamp   time.Time
}

// PubsubMsgHandler adds record(s) to ErrorLogs table in database,
// after parsing build log and compares the result with config yaml
func PubsubMsgHandler(db *sql.DB, configURL, buildLogURL, jobname string, prNumber int) error {
	config, err := config.ParseYaml(configURL)
	if err != nil {
		return err
	}

	errorLogs, err := ParseLog(buildLogURL, config.CollectErrorPatterns())
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(logInsertStmt)
	defer stmt.Close()

	if err != nil {
		return mysql.RollbackTx(tx, err)
	}

	for _, errorLog := range errorLogs {
		if _, err := stmt.Exec(errorLog.Pattern, errorLog.Msg, jobname, prNumber, buildLogURL, time.Now()); err != nil {
			return mysql.RollbackTx(tx, err)
		}
	}

	return tx.Commit()
}
