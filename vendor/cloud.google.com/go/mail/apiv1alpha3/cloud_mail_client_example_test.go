// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AUTO-GENERATED CODE. DO NOT EDIT.

package mail_test

import (
	"cloud.google.com/go/mail/apiv1alpha3"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	mailpb "google.golang.org/genproto/googleapis/cloud/mail/v1alpha3"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

func ExampleNewCloudMailClient() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleCloudMailClient_ListDomains() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListDomainsRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ListDomains(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_GetDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.DeleteDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UndeleteDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UndeleteDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UndeleteDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ExpungeDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ExpungeDomainRequest{
		// TODO: Fill request struct fields.
	}
	err = c.ExpungeDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_TestReceiptRules() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.TestReceiptRulesRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.TestReceiptRules(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_VerifyDomain() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.VerifyDomainRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.VerifyDomain(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ListEmailVerifiedAddresses() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListEmailVerifiedAddressesRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ListEmailVerifiedAddresses(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_GetEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.DeleteEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UndeleteEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UndeleteEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UndeleteEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ExpungeEmailVerifiedAddress() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ExpungeEmailVerifiedAddressRequest{
		// TODO: Fill request struct fields.
	}
	err = c.ExpungeEmailVerifiedAddress(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_RequestEmailVerification() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.RequestEmailVerificationRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RequestEmailVerification(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_VerifyEmail() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.VerifyEmailRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.VerifyEmail(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ListSenders() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListSendersRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListSenders(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleCloudMailClient_GetSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetSenderRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateSenderRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateSenderRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteSenderRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.DeleteSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UndeleteSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UndeleteSenderRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UndeleteSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ExpungeSender() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ExpungeSenderRequest{
		// TODO: Fill request struct fields.
	}
	err = c.ExpungeSender(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_SendMessage() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.SendMessageRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.SendMessage(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ListSmtpCredentials() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListSmtpCredentialsRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ListSmtpCredentials(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_GetSmtpCredential() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetSmtpCredentialRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetSmtpCredential(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateSmtpCredential() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateSmtpCredentialRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateSmtpCredential(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateSmtpCredential() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateSmtpCredentialRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateSmtpCredential(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteSmtpCredential() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteSmtpCredentialRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteSmtpCredential(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_ListReceiptRules() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListReceiptRulesRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListReceiptRules(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleCloudMailClient_GetReceiptRule() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetReceiptRuleRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetReceiptRule(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateReceiptRule() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateReceiptRuleRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateReceiptRule(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateReceiptRule() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateReceiptRuleRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateReceiptRule(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteReceiptRule() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteReceiptRuleRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteReceiptRule(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_ListAddressSets() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ListAddressSetsRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ListAddressSets(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_GetAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.GetAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_CreateAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.CreateAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UpdateAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UpdateAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_DeleteAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.DeleteAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.DeleteAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_UndeleteAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.UndeleteAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UndeleteAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_ExpungeAddressSet() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &mailpb.ExpungeAddressSetRequest{
		// TODO: Fill request struct fields.
	}
	err = c.ExpungeAddressSet(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudMailClient_GetIamPolicy() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.GetIamPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetIamPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_SetIamPolicy() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.SetIamPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.SetIamPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudMailClient_TestIamPermissions() {
	ctx := context.Background()
	c, err := mail.NewCloudMailClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.TestIamPermissionsRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.TestIamPermissions(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
