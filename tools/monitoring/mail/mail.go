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

package mail

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

// Config stores the sender information for mail
type Config struct {
	senderEmail    string
	senderPassword string
}

// NewMailConfig creates a config with a valid sender info
func NewMailConfig(mailAddrFile string, mailPassFile string) (*Config, error) {
	mail, err := ioutil.ReadFile(mailAddrFile)
	if err != nil {
		return nil, err
	}

	pass, err := ioutil.ReadFile(mailPassFile)
	if err != nil {
		return nil, err
	}

	return &Config{
		senderEmail:    string(mail),
		senderPassword: string(pass),
	}, nil
}

// Send sends an email
func (c *Config) Send(recipients []string, subject string, body string) error {
	msg := buildMessage(c.senderEmail, recipients, subject, body)

	err := smtp.SendMail(buildServerName(smtpHost, smtpPort),
		smtp.PlainAuth("", c.senderEmail, c.senderPassword, smtpHost),
		c.senderEmail, recipients, []byte(msg))
	if err != nil {
		return err
	}

	log.Print("Message sent.")
	return nil
}

func buildServerName(host string, port string) string {
	return host + ":" + port
}

func buildMessage(sender string, recipients []string, subject string, body string) string {
	message := ""
	message += fmt.Sprintf("From: %s\n", sender)
	if len(recipients) > 0 {
		message += fmt.Sprintf("To: %s\n", strings.Join(recipients, ";"))
	}

	message += fmt.Sprintf("Subject: %s\n", subject) + "\n"
	message += body

	return message
}
