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
	"context"
	"log"
	"strings"

	"github.com/knative/test-infra/shared/gcs"
	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"
	"github.com/knative/test-infra/tools/monitoring/mysql"
	"github.com/knative/test-infra/tools/monitoring/prowapi"
	"github.com/knative/test-infra/tools/monitoring/subscriber"
)

const yamlURL = "https://raw.githubusercontent.com/knative/test-infra/master/tools/monitoring/config/sample.yaml"

// Client holds all the resources required to run alerting
type Client struct {
	pubsubClient *subscriber.Client
	db           *mysql.DB
}

// Setup sets up the client required to run alerting workflow
func Setup(psClient *subscriber.Client, db *mysql.DB) *Client {
	return &Client{
		pubsubClient: psClient,
		db:           db,
	}
}

// RunAlerting start the alerting workflow
func (c *Client) RunAlerting() {
	log.Println("Starting alerting workflow")
	go func() {
		err := c.pubsubClient.ReceiveMessageAckAll(context.Background(), c.handleReportMessage)
		if err != nil {
			log.Printf("Failed to retrieve messages due to %v", err)
		}
	}()
}

func (c *Client) handleReportMessage(rmsg *prowapi.ReportMessage) {
	if rmsg.Status == prowapi.SuccessState || rmsg.Status == prowapi.FailureState {
		config, err := config.ParseYaml(yamlURL)
		if err != nil {
			log.Printf("Failed to config yaml (%v): %v\n", config, err)
			return
		}

		content, err := gcs.ReadURL(context.Background(), buildLogPath(gubernatortoGcsLink(rmsg.GCSPath)))
		if err != nil {
			log.Printf("Failed to read from url %s. Error: %v\n", rmsg.GCSPath, err)
			return
		}
		log.Printf("Read content: %s\n", content)

		errorLogs, err := log_parser.ParseLog(content, config.CollectErrorPatterns())
		if err != nil {
			log.Printf("Failed to parse content %v. Error: %v\n", string(content), err)
			return
		}

		log.Printf("Parsed errorLogs: %v\n", errorLogs)

		for _, el := range errorLogs {
			if len(rmsg.Refs) <= 0 || len(rmsg.Refs[0].Pulls) <= 0 {
				err = c.db.InsertErrorLog(el.Pattern, el.Msg, rmsg.JobName, 0, rmsg.GCSPath)
			} else {
				err = c.db.InsertErrorLog(el.Pattern, el.Msg, rmsg.JobName, rmsg.Refs[0].Pulls[0].Number, rmsg.GCSPath)
			}
			if err != nil {
				log.Printf("Failed to insert error to db %+v\n", err)
			}
		}

		// TODO(yt3liu): check sending alert
	}
}

// TODO(yt3liu): Remove this hack after the gcs path does not contain the gubernator link
func gubernatortoGcsLink(link string) string {
	return strings.Replace(link, "https://gubernator.knative.dev/build/", "", 1)
}

func buildLogPath(gcsDir string) string {
	return gcsDir + "build-log.txt"
}
