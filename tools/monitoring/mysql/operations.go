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
	"time"

	"github.com/knative/test-infra/shared/mysql"
	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"
)

const (
	insertStmt = `
	INSERT INTO ErrorLogs (
		ErrorPattern, ErrorMsg, JobName, PRNumber, BuildLogURL, TimeStamp
		) VALUES (?,?,?,?,?,?)`
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

	stmt, err := tx.Prepare(insertStmt)
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
