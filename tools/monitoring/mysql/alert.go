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
			Sent, ErrorPattern
		) VALUES (?,?)`

	alertUpdateStmt = `
	UPDATE Alerts
	SET Sent = (?)
	WHERE ErrorPattern = (?)`

	emailTemplate = `In the past %v, 
The number of occurrences of the following error pattern reached threshold:
%s

Hint for diagnose & recovery: %s
`
)

type MailConfig struct {
	*mail.Config
	recipients []string
}

func (m *MailConfig) sendAlert(errorPattern string, config *config.SelectedConfig) error {
	log.Printf("sending alert...")
	subject := fmt.Sprintf("Error pattern reached alerting threshold: %s", errorPattern)
	body := fmt.Sprintf(emailTemplate, config.Duration(), errorPattern, config.Hint)

	return m.mailConfig.Send(m.recipients, subject, body)
}

// Alert checks alert condition and alerts table and send alert mail conditionally
func (m *MailConfig) Alert(errorPattern string, config *config.SelectedConfig, db *sql.DB) (bool, error) {
	if ok, err := config.CheckAlertCondition(errorPattern, db); err != nil || !ok {
		return false, err
	}

	queryTemplate, err := checkAlertsTable(errorPattern, config.Duration(), db)
	if err != nil || queryTemplate == "" {
		return false, err
	}

	if err := updateAlertsTable(queryTemplate, errorPattern, db); err != nil {
		return false, err
	}

	err = m.sendAlert(errorPattern, config)
	return err == nil, err
}

// checkAlertsTable checks alert table and see if it is necessary to send alert email
// returns a sql statement that contains db update action to perform. An empty string indicates
// no db operations to perform and no alert email needs to be sent
func checkAlertsTable(errorPattern string, window time.Duration, db *sql.DB) (string, error) {
	var id int
	var sent time.Time

	row := db.QueryRow(`
		SELECT ID, Sent
		FROM Alerts
		WHERE ErrorPattern = ?`,
		errorPattern)

	if err := row.Scan(&id, &sent); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		// if no record found, instruct to add a record
		return alertInsertStmt, nil
	}

	if sent.Add(window).Before(time.Now()) {
		// if previous alert expires. Instruct to update the timestamp
		log.Printf("previous alert timestamp=%v expired, alert window size=%v", sent, window)
		return alertUpdateStmt, nil
	}

	log.Printf("previous alert not expired (timestamp=%v), "+
		"alert window size=%v, no alert will be sent", sent, window)
	return "", nil
}

func updateAlertsTable(queryTemplate, errorPattern string, db *sql.DB) error {
	if queryTemplate == "" {
		return nil
	}
	_, err := db.Query(queryTemplate, time.Now(), errorPattern)
	return err
}
