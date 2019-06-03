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
		setupSender(client, ctx)
	}

	if *actionSetupAll {
		createDomain(client, ctx)
		setupSender(client, ctx)
	}

	if *actionSendTestMail {
		fmt.Println("Sending a Test Email")
		failIfError("Failed to send test message", client.SendTestMessage(ctx, *toAddr))
	}
}

func createDomain(client *cloudmail.MailClient, ctx context.Context) {
	fmt.Println("Creating the email domain")
	if err := client.CreateDomain(ctx); err != nil {
		log.Fatalf("Failed to create domain %v", err)
	}
}

func setupSender(client *cloudmail.MailClient, ctx context.Context) {
	fmt.Println("Setting up the sender")
	failIfError("Failed to create address set %v", client.CreateAddressSet(ctx))
	failIfError("Failed to create sender domain %v", client.CreateSenderDomain(ctx))
	failIfError("Failed to setup receipt rule %v", client.CreateAndApplyReceiptRuleDrop(ctx))
}

func failIfError(errFmtMsg string, err error) {
	if err != nil {
		log.Fatalf(errFmtMsg, err)
	}
}
