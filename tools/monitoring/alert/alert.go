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

package alert

import (
	"database/sql"
	"log"
	"time"

	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/mail"
)

const (
	alertInsertStmt = `
		INSERT INTO Alerts (Sent, ErrorPattern) VALUES (?,?)
		ON DUPLICATE KEY UPDATE Sent = (?)`
)

type MailConfig struct {
	*mail.Config
	recipients []string
}

func (m *MailConfig) sendAlert(c *mailContent) error {
	log.Printf("sending alert...")
	return m.Send(m.recipients, c.subject(), c.body())
}

// Alert checks alert condition and alerts table and send alert mail conditionally
func (m *MailConfig) Alert(errorPattern string, s *config.SelectedConfig, db *sql.DB) (bool, error) {

	errorLogs, err := GetErrorLogs(s, errorPattern, db)
	if err != nil {
		return false, err
	}

	report := newReport(errorLogs)
	if !report.CheckAlertCondition(s) {
		return false, nil
	}

	ok, err := checkAlertsTable(errorPattern, s.Duration(), db)
	if err != nil || !ok {
		return false, err
	}

	if err := updateAlertsTable(errorPattern, db); err != nil {
		return false, err
	}

	content := mailContent{*report, errorPattern, s.Hint, s.Duration()}
	err = m.sendAlert(&content)
	return err == nil, err
}

// checkAlertsTable checks alert table and see if it is necessary to send alert email
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

		// if no record found
		return true, nil
	}

	if sent.Add(window).Before(time.Now()) {
		// if previous alert expires.
		log.Printf("previous alert timestamp=%v expired, alert window size=%v", sent, window)
		return true, nil
	}

	log.Printf("previous alert not expired (timestamp=%v), "+
		"alert window size=%v, no alert will be sent", sent, window)
	return false, nil
}

func updateAlertsTable(errorPattern string, db *sql.DB) error {
	now := time.Now()
	_, err := db.Query(alertInsertStmt, now, errorPattern, now)
	return err
}
