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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/knative/test-infra/tools/monitoring/cloudmail"
)

func main() {
	actionCreateDomain := flag.Bool("setup-domain", false, "Create the cloud mail domain")
	actionSetupSender := flag.Bool("setup-sender", false, "Setup the address set, sender, receipt rule with default settings")
	actionSetupAll := flag.Bool("setup-all", false, "Setup up the domain, address set, sender, receipt rules with default settings")
	actionSendTestMail := flag.Bool("send-test-mail", false, "Send a test message")

	domainID := flag.String("domain-id", "", "Cloud Mail domain ID")
	domainName := flag.String("domain-name", "", "The domain name used to send the test email")
	toAddr := flag.String("to-address", "", "The test email recipient address")

	flag.Parse()

	ctx := context.Background()
	client, err := cloudmail.NewMailClient(ctx)
	if err != nil {
		failIfError("Failed to create Cloud Mail client: %s", err)
	}

	if *actionCreateDomain {
		createDomain(client, ctx)
	}

	if *actionSetupSender {
		setupSender(client, ctx, *domainID)
	}

	if *actionSetupAll {
		domainID := createDomain(client, ctx)
		setupSender(client, ctx, domainID)
	}

	if *actionSendTestMail {
		fmt.Println("Sending a Test Email")
		failIfError("Failed to send test message", client.SendTestMessage(ctx, *domainName, *toAddr))
	}
}

func createDomain(client *cloudmail.MailClient, ctx context.Context) string {
	fmt.Println("Creating the email domain")
	domainID, err := client.CreateDomain(ctx)
	if err != nil {
		log.Fatalf("Failed to create domain %v", err)
	}
	return domainID
}

func setupSender(client *cloudmail.MailClient, ctx context.Context, domainID string) {
	fmt.Println("Setting up the sender")
	failIfError("Failed to create address set %v", client.CreateAddressSet(ctx, domainID))
	failIfError("Failed to create sender domain %v", client.CreateSenderDomain(ctx, domainID))
	failIfError("Failed to setup receipt rule %v", client.CreateAndApplyReceiptRuleDrop(ctx, domainID))
}

func failIfError(errFmtMsg string, err error) {
	if err != nil {
		log.Fatalf(errFmtMsg, err)
	}
}
