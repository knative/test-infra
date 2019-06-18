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
	"fmt"
	"time"

	"github.com/knative/test-infra/shared/mysql"
)

// MonitoringDatabase holds an active database connection created in `config`
type MonitoringDatabase struct {
	*sql.DB
	config mysql.DBConfig
}

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

// String returns the string representation of the struct used in alert message
func (e ErrorLog) String() string {
	return fmt.Sprintf("[%v] %s (Job: %s, PR: %v, BuildLog: %s)",
		e.TimeStamp, e.Msg, e.JobName, e.PRNumber, e.BuildLogURL)
}

func (db *MonitoringDatabase) insertErrorLog(errPat string, errMsg string, jobName string, prNum int, blogURL string) error {
	stmt, err := db.Prepare(`INSERT INTO ErrorLogs(ErrorPattern, ErrorMsg, JobName, PRNumber, BuildLogURL, TimeStamp)
				VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err = execAffectingOneRow(stmt, errPat, errMsg, jobName, prNum, blogURL, time.Now())
	return err
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}
	return r, nil
}
