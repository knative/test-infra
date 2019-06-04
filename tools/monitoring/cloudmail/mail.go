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
	matchMode      = "PREFIX"
)

// MailClient holds instance of CloudMailClient for making mail-related requests
type MailClient struct {
	*mail.CloudMailClient
	DomainName string
	DomainID   string
}

// NewMailClient creates a new MailClient to handle all the mail interaction
func NewMailClient(ctx context.Context) (*MailClient, error) {
	c, err := mail.NewCloudMailClient(ctx)
	return &MailClient{
		CloudMailClient: c,
		// TODO(yt3liu): Update it with the domain name created in knative-tests
		DomainName: "",
		DomainID:   "",
	}, err
}

// CreateDomain creates a new project domain for cloud mail
func (c *MailClient) CreateMailDomain(ctx context.Context) error {
	domain := &mailpb.Domain{
		ProjectDomain: true,
		DomainName:    "",
	}
	req := &mailpb.CreateDomainRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Region: region,
		Domain: domain,
	}
	resp, err := c.CreateDomain(ctx, req)
	if err != nil {
		return err
	}

	c.DomainID = strings.Replace(resp.GetName(), fmt.Sprintf("regions/%s/domains/", region), "", 1)
	c.DomainName = resp.GetDomainName()

	fmt.Printf("Domain created.\n")
	fmt.Printf("Domain name: %s\n", c.DomainName)
	fmt.Printf("Domain region: %s\n", region)
	fmt.Printf("Domain ID: %s\n", c.DomainID)

	return nil
}

// CreateAddressSet enables email addresses under the domain
func (c *MailClient) CreateMailAddressSet(ctx context.Context) error {
	addressSet := &mailpb.AddressSet{
		AddressPatterns: []string{addressPattern},
	}
	req := &mailpb.CreateAddressSetRequest{
		Parent:       fmt.Sprintf("regions/%s/domains/%s", region, c.DomainID),
		AddressSetId: addressSetID,
		AddressSet:   addressSet,
	}

	if _, err := c.CreateAddressSet(ctx, req); err != nil {
		return err
	}
	fmt.Printf("Address set created.\nAddress set ID: %s\n", addressSetID)
	return nil
}

// CreateSenderDomain configures the sender email
func (c *MailClient) CreateSenderDomain(ctx context.Context) error {
	addressSetPath := fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", region, c.DomainID, addressSetID)
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

	if _, err := c.CreateSender(ctx, req); err != nil {
		return err
	}
	fmt.Printf("Sender created.\nSender ID: %s\n", senderID)
	return nil
}

// CreateAndApplyReceiptRuleDrop configures the bounce message behaviour to do nothing when an email cannot be delivered
func (c *MailClient) CreateAndApplyReceiptRuleDrop(ctx context.Context) error {
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
		Parent:      fmt.Sprintf("regions/%s/domains/%s", region, c.DomainID),
		RuleId:      receiptRuleID,
		ReceiptRule: receiptRule,
	}

	if _, err := c.CreateReceiptRule(ctx, createReceiptRuleReq); err != nil {
		return err
	}
	fmt.Printf("Receipt rule %s created.\n", receiptRuleID)

	receiptRuleset := &mailpb.ReceiptRuleset{
		ReceiptRules: []string{fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", region, c.DomainID, receiptRuleID)},
	}
	domain := &mailpb.Domain{
		Name:           fmt.Sprintf("regions/%s/domains/%s", region, c.DomainID),
		ReceiptRuleset: receiptRuleset,
	}

	mask := &field_mask.FieldMask{
		Paths: []string{"receipt_ruleset"},
	}
	updateDomainReq := &mailpb.UpdateDomainRequest{
		Domain:     domain,
		UpdateMask: mask,
	}

	if _, err := c.UpdateDomain(ctx, updateDomainReq); err != nil {
		return err
	}
	fmt.Printf("New receipt rule %s applied to domain %s.\n", receiptRuleID, c.DomainID)
	return nil
}

// SendTestMessage sends a test email
func (c *MailClient) SendTestMessage(ctx context.Context, toAddress string) error {
	return c.SendEmailMessage(ctx, toAddress,
		"Knative Monitoring Cloud Mail Test",
		"This is a test message.")
}

// SendEmailMessage sends an email
func (c *MailClient) SendEmailMessage(ctx context.Context, toAddress string, subject string, body string) error {
	fromAddress := fmt.Sprintf("%s@%s", addressPattern, c.DomainName)

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

	resp, err := c.SendMessage(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("Messsage sent.\nMessage ID: %s\n", resp.GetRfc822MessageId())
	return nil
}
