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

	mail "cloud.google.com/go/mail/apiv1alpha3"
	"github.com/knative/test-infra/tools/monitoring/cloudmail"
)

func main() {
	actionCreateDomain := flag.Bool("setup-domain", false, "Create the cloud mail domain")
	actionSetupSender := flag.Bool("setup-sender", false, "Setup the address set, sender, receipt rule with default settings")
	actionSendTestMail := flag.Bool("send-test-mail", false, "Send a test message")

	domainID := flag.String("domain-id", "", "Cloud Mail domain ID")
	domainName := flag.String("domain-name", "", "The domain name used to send the test email")
	toAddr := flag.String("to-address", "", "The test email recipient address")

	flag.Parse()

	ctx := context.Background()

	client, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		handleError("Failed to create Cloud Mail client: %s", err)
	}

	if *actionCreateDomain {
		fmt.Println("Creating the email domain")
		handleError("Failed to create the domain", cloudmail.CreateDomain(ctx, client))
	}

	if *actionSetupSender {
		fmt.Println("Setting up the sender")
		handleError("Failed to create address set", cloudmail.CreateAddressSet(ctx, client, *domainID))
		handleError("Failed to create sender domain", cloudmail.CreateSenderDomain(ctx, client, *domainID))
		handleError("Failed to setup receipt rule", cloudmail.CreateAndApplyReceiptRuleDrop(ctx, client, *domainID))
	}

	if *actionSendTestMail {
		fmt.Println("Sending a Test Email")
		handleError("Failed to send test message", cloudmail.SendTestMessage(ctx, client, *domainName, *toAddr))
	}
}

func handleError(errFmtMsg string, err error) {
	if err != nil {
		log.Fatalf(errFmtMsg, err)
	}
}
