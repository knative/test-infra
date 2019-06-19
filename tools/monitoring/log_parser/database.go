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

// DB holds an active database connection created in `config`
type DB struct {
	*sql.DB
	Config *mysql.DBConfig
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

// NewDB returns the DB object with an active database connection
func NewDB(c *mysql.DBConfig) (*DB, error) {
	db, err := c.Connect()
	return &DB{db, c}, err
}

// InsertErrorLog insert a new error to the ErrorLogs table
func (db *DB) InsertErrorLog(errPat string, errMsg string, jobName string, prNum int, blogURL string) error {
	stmt, err := db.Prepare(`INSERT INTO ErrorLogs(ErrorPattern, ErrorMsg, JobName, PRNumber, BuildLogURL, TimeStamp)
				VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(errPat, errMsg, jobName, prNum, blogURL, time.Now())
	return err
}
