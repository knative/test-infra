package mysql

import (
	"database/sql"
	"github.com/knative/test-infra/shared/mysql"
	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"
	"time"
)

const (
	logInsertStmt = `
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
