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

package cloudmail

import (
	"context"
	"fmt"
	"strings"

	mail "cloud.google.com/go/mail/apiv1alpha3"
	mailpb "google.golang.org/genproto/googleapis/cloud/mail/v1alpha3"
	"google.golang.org/genproto/protobuf/field_mask"
)

const (
	projectID = "knative-tests"
	region    = "us-central1"

	addressSetID   = "monitoring-address-set"
	addressPattern = "monitoring-alert"
	senderID       = "monitoring-alert-sender"
	receiptRuleID  = "monitoring-receipt-drop-rule"
)

// CreateDomain creates a new project domain for cloud main
func CreateDomain(ctx context.Context, client *mail.CloudMailClient) error {
	domain := &mailpb.Domain{
		ProjectDomain: true,
		DomainName:    "",
	}
	req := &mailpb.CreateDomainRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Region: region,
		Domain: domain,
	}
	resp, err := client.CreateDomain(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("Domain created.\n")
	fmt.Printf("Domain name: %s", resp.GetDomainName())
	fmt.Printf("Domain region: %s", region)
	fmt.Printf("Domain ID: %s", strings.Replace(resp.GetName(), fmt.Sprintf("regions/%s/domains/", region), "", 1))
	return nil
}

// CreateAddressSet enables email addresses under the domain
func CreateAddressSet(ctx context.Context, client *mail.CloudMailClient, domainID string) error {
	addressSet := &mailpb.AddressSet{
		AddressPatterns: []string{addressPattern},
	}
	req := &mailpb.CreateAddressSetRequest{
		Parent:       fmt.Sprintf("regions/%s/domains/%s", region, domainID),
		AddressSetId: addressSetID,
		AddressSet:   addressSet,
	}

	if _, err := client.CreateAddressSet(ctx, req); err != nil {
		return err
	}
	fmt.Printf("Address set created.\nAddress set ID: %s\n", addressSetID)
	return nil
}

// CreateSenderDomain configures the sender email
func CreateSenderDomain(ctx context.Context, client *mail.CloudMailClient, domainID string) error {
	addressSetPath := fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", region, domainID, addressSetID)
	sender := &mailpb.Sender{
		DefaultEnvelopeFromAuthority: addressSetPath,
		DefaultHeaderFromAuthority:   addressSetPath,
	}
	req := &mailpb.CreateSenderRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		Region:   region,
		SenderId: senderID,
		Sender:   sender,
	}

	if _, err := client.CreateSender(ctx, req); err != nil {
		return err
	}
	fmt.Printf("Sender created.\nSender ID: %s\n", senderID)
	return nil
}

// CreateAndApplyReceiptRuleDrop configures the bounce message behaviour to do nothing when an email cannot be delivered
func CreateAndApplyReceiptRuleDrop(ctx context.Context, client *mail.CloudMailClient, domainID string) error {
	matchMode := "PREFIX"

	envelopeToPattern := &mailpb.ReceiptRule_Pattern{
		Pattern:   addressPattern,
		MatchMode: mailpb.ReceiptRule_Pattern_MatchMode(mailpb.ReceiptRule_Pattern_MatchMode_value[matchMode]),
	}
	dropAction := &mailpb.ReceiptAction_Drop{
		Drop: &mailpb.DropAction{},
	}
	action := &mailpb.ReceiptAction{
		Action: dropAction,
	}
	receiptRule := &mailpb.ReceiptRule{
		EnvelopeToPatterns: []*mailpb.ReceiptRule_Pattern{envelopeToPattern},
		Action:             action,
	}
	createReceiptRuleReq := &mailpb.CreateReceiptRuleRequest{
		Parent:      fmt.Sprintf("regions/%s/domains/%s", region, domainID),
		RuleId:      receiptRuleID,
		ReceiptRule: receiptRule,
	}

	if _, err := client.CreateReceiptRule(ctx, createReceiptRuleReq); err != nil {
		return err
	}
	fmt.Printf("Receipt rule %s created.\n", receiptRuleID)

	receiptRuleset := &mailpb.ReceiptRuleset{
		ReceiptRules: []string{fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", region, domainID, receiptRuleID)},
	}
	domain := &mailpb.Domain{
		Name:           fmt.Sprintf("regions/%s/domains/%s", region, domainID),
		ReceiptRuleset: receiptRuleset,
	}

	mask := &field_mask.FieldMask{
		Paths: []string{"receipt_ruleset"},
	}
	updateDomainReq := &mailpb.UpdateDomainRequest{
		Domain:     domain,
		UpdateMask: mask,
	}

	if _, updateDomainErr := client.UpdateDomain(ctx, updateDomainReq); updateDomainErr != nil {
		return updateDomainErr
	}
	fmt.Printf("New receipt rule %s applied to domain %s.\n", receiptRuleID, domainID)
	return nil
}

// SendTestMessage sends a test email
func SendTestMessage(ctx context.Context, client *mail.CloudMailClient, domainName string, toAddress string) error {
	return SendEmailMessage(ctx, client, domainName, toAddress,
		"Knative Monitoring Cloud Mail Test",
		"This is a test message.")
}

// SendEmailMessage sends an email
func SendEmailMessage(ctx context.Context, client *mail.CloudMailClient, domainName string, toAddress string, subject string, body string) error {
	fromAddress := fmt.Sprintf("%s@%s", addressPattern, domainName)

	from := &mailpb.Address{
		AddressSpec: fromAddress,
	}
	to := &mailpb.Address{
		AddressSpec: toAddress,
	}
	message := &mailpb.SendMessageRequest_SimpleMessage{
		SimpleMessage: &mailpb.SimpleMessage{
			From:     from,
			To:       []*mailpb.Address{to},
			Subject:  subject,
			TextBody: body,
		},
	}

	req := &mailpb.SendMessageRequest{
		Sender: fmt.Sprintf("projects/%s/regions/%s/senders/%s", projectID, region, senderID),
		// If EnvelopeFromAuthority or HeaderFromAuthority is empty,
		// Cloud Mail will fall back on the default value in the provided
		// sender. There must be an envelope from authority and a header from
		// authority that Cloud Mail can use; otherwise the system will
		// return an error.
		EnvelopeFromAuthority: "",
		HeaderFromAuthority:   "",
		EnvelopeFromAddress:   fromAddress,
		Message:               message,
	}

	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("Messsage sent.\nMessage ID: %s\n", resp.GetRfc822MessageId())
	return nil
}
