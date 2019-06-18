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

package workflow

import (
	"context"
	"log"
	"strings"

	"github.com/knative/test-infra/shared/gcs"
	"github.com/knative/test-infra/tools/monitoring/config"
	"github.com/knative/test-infra/tools/monitoring/log_parser"
	"github.com/knative/test-infra/tools/monitoring/prowapi"
	"github.com/knative/test-infra/tools/monitoring/subscriber"
)

const yamlURL = "https://raw.githubusercontent.com/knative/test-infra/master/tools/monitoring/config/sample.yaml"

// RunAlerting start the alerting workflow
func RunAlerting(client *subscriber.Client) {
	log.Println("Starting alerting workflow")
	go func() {
		err := client.ReceiveMessageAckAll(context.Background(), handleReportMessage)
		if err != nil {
			log.Printf("Failed to retrieve messages due to %v", err)
		}
	}()
}

func handleReportMessage(rmsg *prowapi.ReportMessage) {
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
