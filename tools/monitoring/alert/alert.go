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
	"log"

	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/mail"
	"github.com/knative/test-infra/tools/monitoring/mysql"
)

type MailConfig struct {
	*mail.Config
	Recipients []string
}

func (m *MailConfig) sendAlert(c *mailContent) error {
	log.Printf("sending alert...")
	return m.Send(m.Recipients, c.subject(), c.body())
}

// Alert checks alert condition and alerts table and send alert mail conditionally
func (m *MailConfig) Alert(errorPattern string, s *config.SelectedConfig, db *mysql.DB) (bool, error) {
	log.Println("Fetcing error logs")
	errorLogs, err := db.GetErrorLogs(s, errorPattern)
	if err != nil {
		return false, err
	}

	log.Println("Building Report and checking alert conditions")
	report := newReport(errorLogs)
	if !report.CheckAlertCondition(s) {
		return false, nil
	}

	log.Println("checking if the alert is a fresh alert pattern")
	ok, err := db.IsFreshAlertPattern(errorPattern, s.Duration())
	if err != nil || !ok {
		return false, err
	}

	log.Println("Adding the new alert pattern to the database")
	if err := db.AddAlert(errorPattern); err != nil {
		return false, err
	}

	log.Println("Generating and sending the alert email")
	content := mailContent{*report, errorPattern, s.Hint, s.Duration()}
	err = m.sendAlert(&content)
	return err == nil, err
}
