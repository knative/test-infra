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
	"log"
	"time"

	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/mail"
)

const (
	alertInsertStmt = `
	INSERT INTO Alerts (
		ErrorPattern, Sent
	) VALUES (?,?)`

	emailTemplate = `In the past %v, 
The number of occurrences of the following error pattern reached threshold:
%s

Hint for diagnose & recovery: %s
`
)

func sendAlert(errorPattern string, config *config.SelectedConfig, mailConfig *mail.Config, recipients []string) error {
	log.Printf("sending alert...")
	subject := fmt.Sprintf("Error pattern reached alerting threshold: %s", errorPattern)
	body := fmt.Sprintf(emailTemplate, config.Duration(), errorPattern, config.Hint)

	return mailConfig.Send(recipients, subject, body)
}

// CheckAlertCondition checks whether the given error pattern meets
// the alert condition specified in config
func CheckAlertCondition(errorPattern string, config *config.SelectedConfig, db *sql.DB) (bool, error) {
	// the timestamp we want to start collecting logs
	startTime := time.Now().Add(config.Duration())

	var nMatches, nJobs, nPRs int

	row := db.QueryRow(`
		SELECT 
			COUNT(*),
			COUNT (DISTINCT Jobname),
			COUNT (DISTINCT PrNumber)
		FROM ErrorLogs
		WHERE ErrorPattern=? and TimeStamp > ?`,
		errorPattern, startTime)

	if err := row.Scan(&nMatches, &nJobs, &nPRs); err != nil {
		return false, err
	}

	return nMatches >= config.Occurrences && nJobs >= config.JobsAffected && nPRs >= config.PrsAffected, nil
}

// Alert checks alert condition and alerts table and send alert mail conditionally
func Alert(errorPattern string, config *config.SelectedConfig, db *sql.DB, mailConfig *mail.Config, recipients []string) (bool, error) {
	alertConditionMet, err := CheckAlertCondition(errorPattern, config, db)
	if err != nil {
		return false, err
	}
	if !alertConditionMet {
		return false, nil
	}

	shouldSendMail, err := checkAlertsTable(errorPattern, config.Duration(), db)
	if err != nil {
		return false, err
	}
	if !shouldSendMail {
		return false, nil
	}
	if err = sendAlert(errorPattern, config, mailConfig, recipients); err != nil {
		return false, err
	}
	return true, nil
}

// checkAlertsTable checks alert table and see if it is necessary to send alert email
// also updates the alerts table if email sent
func checkAlertsTable(errorPattern string, window time.Duration, db *sql.DB) (bool, error) {
	var id int
	var sent time.Time

	row := db.QueryRow(`
		SELECT ID, Sent
		FROM Alerts
		WHERE ErrorPattern = ?`,
		errorPattern)

	if err := row.Scan(&id, &sent); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		_, err := db.Query(alertInsertStmt, errorPattern, time.Now())
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if sent.Add(window).Before(time.Now()) {
		log.Printf("previous alert timestamp=%v expired, alert window size=%v", sent, window)
		return true, nil
	}

	log.Printf("previous alert not expired (timestamp=%v), "+
		"alert window size=%v, no alert will be sent", sent, window)
	return false, nil
}
