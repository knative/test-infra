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

package mail

import (
	emptypb "github.com/golang/protobuf/ptypes/empty"
	mailpb "google.golang.org/genproto/googleapis/cloud/mail/v1alpha3"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	field_maskpb "google.golang.org/genproto/protobuf/field_mask"
)

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
)

var _ = io.EOF
var _ = ptypes.MarshalAny
var _ status.Status

type mockCloudMailServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	mailpb.CloudMailServer

	reqs []proto.Message

	// If set, all calls return this error.
	err error

	// responses to return if err == nil
	resps []proto.Message
}

func (s *mockCloudMailServer) ListDomains(ctx context.Context, req *mailpb.ListDomainsRequest) (*mailpb.ListDomainsResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListDomainsResponse), nil
}

func (s *mockCloudMailServer) GetDomain(ctx context.Context, req *mailpb.GetDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) CreateDomain(ctx context.Context, req *mailpb.CreateDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) UpdateDomain(ctx context.Context, req *mailpb.UpdateDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) DeleteDomain(ctx context.Context, req *mailpb.DeleteDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) UndeleteDomain(ctx context.Context, req *mailpb.UndeleteDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) ExpungeDomain(ctx context.Context, req *mailpb.ExpungeDomainRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) TestReceiptRules(ctx context.Context, req *mailpb.TestReceiptRulesRequest) (*mailpb.TestReceiptRulesResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.TestReceiptRulesResponse), nil
}

func (s *mockCloudMailServer) VerifyDomain(ctx context.Context, req *mailpb.VerifyDomainRequest) (*mailpb.Domain, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Domain), nil
}

func (s *mockCloudMailServer) ListEmailVerifiedAddresses(ctx context.Context, req *mailpb.ListEmailVerifiedAddressesRequest) (*mailpb.ListEmailVerifiedAddressesResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListEmailVerifiedAddressesResponse), nil
}

func (s *mockCloudMailServer) GetEmailVerifiedAddress(ctx context.Context, req *mailpb.GetEmailVerifiedAddressRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) CreateEmailVerifiedAddress(ctx context.Context, req *mailpb.CreateEmailVerifiedAddressRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) UpdateEmailVerifiedAddress(ctx context.Context, req *mailpb.UpdateEmailVerifiedAddressRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) DeleteEmailVerifiedAddress(ctx context.Context, req *mailpb.DeleteEmailVerifiedAddressRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) UndeleteEmailVerifiedAddress(ctx context.Context, req *mailpb.UndeleteEmailVerifiedAddressRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) ExpungeEmailVerifiedAddress(ctx context.Context, req *mailpb.ExpungeEmailVerifiedAddressRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) RequestEmailVerification(ctx context.Context, req *mailpb.RequestEmailVerificationRequest) (*mailpb.RequestEmailVerificationResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.RequestEmailVerificationResponse), nil
}

func (s *mockCloudMailServer) VerifyEmail(ctx context.Context, req *mailpb.VerifyEmailRequest) (*mailpb.EmailVerifiedAddress, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.EmailVerifiedAddress), nil
}

func (s *mockCloudMailServer) ListSenders(ctx context.Context, req *mailpb.ListSendersRequest) (*mailpb.ListSendersResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListSendersResponse), nil
}

func (s *mockCloudMailServer) GetSender(ctx context.Context, req *mailpb.GetSenderRequest) (*mailpb.Sender, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Sender), nil
}

func (s *mockCloudMailServer) CreateSender(ctx context.Context, req *mailpb.CreateSenderRequest) (*mailpb.Sender, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Sender), nil
}

func (s *mockCloudMailServer) UpdateSender(ctx context.Context, req *mailpb.UpdateSenderRequest) (*mailpb.Sender, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Sender), nil
}

func (s *mockCloudMailServer) DeleteSender(ctx context.Context, req *mailpb.DeleteSenderRequest) (*mailpb.Sender, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Sender), nil
}

func (s *mockCloudMailServer) UndeleteSender(ctx context.Context, req *mailpb.UndeleteSenderRequest) (*mailpb.Sender, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.Sender), nil
}

func (s *mockCloudMailServer) ExpungeSender(ctx context.Context, req *mailpb.ExpungeSenderRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) SendMessage(ctx context.Context, req *mailpb.SendMessageRequest) (*mailpb.SendMessageResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.SendMessageResponse), nil
}

func (s *mockCloudMailServer) ListSmtpCredentials(ctx context.Context, req *mailpb.ListSmtpCredentialsRequest) (*mailpb.ListSmtpCredentialsResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListSmtpCredentialsResponse), nil
}

func (s *mockCloudMailServer) GetSmtpCredential(ctx context.Context, req *mailpb.GetSmtpCredentialRequest) (*mailpb.SmtpCredential, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.SmtpCredential), nil
}

func (s *mockCloudMailServer) CreateSmtpCredential(ctx context.Context, req *mailpb.CreateSmtpCredentialRequest) (*mailpb.SmtpCredential, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.SmtpCredential), nil
}

func (s *mockCloudMailServer) UpdateSmtpCredential(ctx context.Context, req *mailpb.UpdateSmtpCredentialRequest) (*mailpb.SmtpCredential, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.SmtpCredential), nil
}

func (s *mockCloudMailServer) DeleteSmtpCredential(ctx context.Context, req *mailpb.DeleteSmtpCredentialRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) ListReceiptRules(ctx context.Context, req *mailpb.ListReceiptRulesRequest) (*mailpb.ListReceiptRulesResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListReceiptRulesResponse), nil
}

func (s *mockCloudMailServer) GetReceiptRule(ctx context.Context, req *mailpb.GetReceiptRuleRequest) (*mailpb.ReceiptRule, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ReceiptRule), nil
}

func (s *mockCloudMailServer) CreateReceiptRule(ctx context.Context, req *mailpb.CreateReceiptRuleRequest) (*mailpb.ReceiptRule, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ReceiptRule), nil
}

func (s *mockCloudMailServer) UpdateReceiptRule(ctx context.Context, req *mailpb.UpdateReceiptRuleRequest) (*mailpb.ReceiptRule, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ReceiptRule), nil
}

func (s *mockCloudMailServer) DeleteReceiptRule(ctx context.Context, req *mailpb.DeleteReceiptRuleRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) ListAddressSets(ctx context.Context, req *mailpb.ListAddressSetsRequest) (*mailpb.ListAddressSetsResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.ListAddressSetsResponse), nil
}

func (s *mockCloudMailServer) GetAddressSet(ctx context.Context, req *mailpb.GetAddressSetRequest) (*mailpb.AddressSet, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.AddressSet), nil
}

func (s *mockCloudMailServer) CreateAddressSet(ctx context.Context, req *mailpb.CreateAddressSetRequest) (*mailpb.AddressSet, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.AddressSet), nil
}

func (s *mockCloudMailServer) UpdateAddressSet(ctx context.Context, req *mailpb.UpdateAddressSetRequest) (*mailpb.AddressSet, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.AddressSet), nil
}

func (s *mockCloudMailServer) DeleteAddressSet(ctx context.Context, req *mailpb.DeleteAddressSetRequest) (*mailpb.AddressSet, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.AddressSet), nil
}

func (s *mockCloudMailServer) UndeleteAddressSet(ctx context.Context, req *mailpb.UndeleteAddressSetRequest) (*mailpb.AddressSet, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*mailpb.AddressSet), nil
}

func (s *mockCloudMailServer) ExpungeAddressSet(ctx context.Context, req *mailpb.ExpungeAddressSetRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

func (s *mockCloudMailServer) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest) (*iampb.Policy, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*iampb.Policy), nil
}

func (s *mockCloudMailServer) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest) (*iampb.Policy, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*iampb.Policy), nil
}

func (s *mockCloudMailServer) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest) (*iampb.TestIamPermissionsResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*iampb.TestIamPermissionsResponse), nil
}

// clientOpt is the option tests should use to connect to the test server.
// It is initialized by TestMain.
var clientOpt option.ClientOption

var (
	mockCloudMail mockCloudMailServer
)

func TestMain(m *testing.M) {
	flag.Parse()

	serv := grpc.NewServer()
	mailpb.RegisterCloudMailServer(serv, &mockCloudMail)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go serv.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	clientOpt = option.WithGRPCConn(conn)

	os.Exit(m.Run())
}

func TestCloudMailListDomains(t *testing.T) {
	var expectedResponse *mailpb.ListDomainsResponse = &mailpb.ListDomainsResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var parent string = "parent-995424086"
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListDomainsRequest{
		Parent:      parent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListDomains(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListDomainsError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var parent string = "parent-995424086"
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListDomainsRequest{
		Parent:      parent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListDomains(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetDomain(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name2,
		Parent:            parent,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.GetDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.GetDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateDomain(t *testing.T) {
	var name string = "name3373707"
	var parent2 string = "parent21175163357"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name,
		Parent:            parent2,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var parent string = "parent-995424086"
	var region string = "region-934795532"
	var domain *mailpb.Domain = &mailpb.Domain{}
	var request = &mailpb.CreateDomainRequest{
		Parent: parent,
		Region: region,
		Domain: domain,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var parent string = "parent-995424086"
	var region string = "region-934795532"
	var domain *mailpb.Domain = &mailpb.Domain{}
	var request = &mailpb.CreateDomainRequest{
		Parent: parent,
		Region: region,
		Domain: domain,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateDomain(t *testing.T) {
	var name string = "name3373707"
	var parent string = "parent-995424086"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name,
		Parent:            parent,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var domain *mailpb.Domain = &mailpb.Domain{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateDomainRequest{
		Domain:     domain,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var domain *mailpb.Domain = &mailpb.Domain{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateDomainRequest{
		Domain:     domain,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteDomain(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name2,
		Parent:            parent,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.DeleteDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailDeleteDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.DeleteDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUndeleteDomain(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name2,
		Parent:            parent,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.UndeleteDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUndeleteDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.UndeleteDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailExpungeDomain(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.ExpungeDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailExpungeDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.ExpungeDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailTestReceiptRules(t *testing.T) {
	var expectedResponse *mailpb.TestReceiptRulesResponse = &mailpb.TestReceiptRulesResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedDomain string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var recipient string = "recipient820081177"
	var receiptRuleset *mailpb.ReceiptRuleset = &mailpb.ReceiptRuleset{}
	var request = &mailpb.TestReceiptRulesRequest{
		Domain:         formattedDomain,
		Recipient:      recipient,
		ReceiptRuleset: receiptRuleset,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.TestReceiptRules(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailTestReceiptRulesError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedDomain string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var recipient string = "recipient820081177"
	var receiptRuleset *mailpb.ReceiptRuleset = &mailpb.ReceiptRuleset{}
	var request = &mailpb.TestReceiptRulesRequest{
		Domain:         formattedDomain,
		Recipient:      recipient,
		ReceiptRuleset: receiptRuleset,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.TestReceiptRules(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailVerifyDomain(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var deleted bool = false
	var domainName string = "domainName104118566"
	var projectDomain bool = true
	var verificationToken string = "verificationToken-498552107"
	var expectedResponse = &mailpb.Domain{
		Name:              name2,
		Parent:            parent,
		Deleted:           deleted,
		DomainName:        domainName,
		ProjectDomain:     projectDomain,
		VerificationToken: verificationToken,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.VerifyDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.VerifyDomain(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailVerifyDomainError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &mailpb.VerifyDomainRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.VerifyDomain(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailListEmailVerifiedAddresses(t *testing.T) {
	var expectedResponse *mailpb.ListEmailVerifiedAddressesResponse = &mailpb.ListEmailVerifiedAddressesResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListEmailVerifiedAddressesRequest{
		Parent:      formattedParent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListEmailVerifiedAddresses(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListEmailVerifiedAddressesError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListEmailVerifiedAddressesRequest{
		Parent:      formattedParent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListEmailVerifiedAddresses(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetEmailVerifiedAddress(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name2,
		Parent:  parent,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.GetEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.GetEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateEmailVerifiedAddress(t *testing.T) {
	var name string = "name3373707"
	var parent2 string = "parent21175163357"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name,
		Parent:  parent2,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var emailVerifiedAddress *mailpb.EmailVerifiedAddress = &mailpb.EmailVerifiedAddress{}
	var request = &mailpb.CreateEmailVerifiedAddressRequest{
		Parent:               formattedParent,
		Region:               region,
		EmailVerifiedAddress: emailVerifiedAddress,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var emailVerifiedAddress *mailpb.EmailVerifiedAddress = &mailpb.EmailVerifiedAddress{}
	var request = &mailpb.CreateEmailVerifiedAddressRequest{
		Parent:               formattedParent,
		Region:               region,
		EmailVerifiedAddress: emailVerifiedAddress,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateEmailVerifiedAddress(t *testing.T) {
	var name string = "name3373707"
	var parent string = "parent-995424086"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name,
		Parent:  parent,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var emailVerifiedAddress *mailpb.EmailVerifiedAddress = &mailpb.EmailVerifiedAddress{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateEmailVerifiedAddressRequest{
		EmailVerifiedAddress: emailVerifiedAddress,
		UpdateMask:           updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var emailVerifiedAddress *mailpb.EmailVerifiedAddress = &mailpb.EmailVerifiedAddress{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateEmailVerifiedAddressRequest{
		EmailVerifiedAddress: emailVerifiedAddress,
		UpdateMask:           updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteEmailVerifiedAddress(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name2,
		Parent:  parent,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.DeleteEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailDeleteEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.DeleteEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUndeleteEmailVerifiedAddress(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name2,
		Parent:  parent,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.UndeleteEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUndeleteEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.UndeleteEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailExpungeEmailVerifiedAddress(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.ExpungeEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeEmailVerifiedAddress(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailExpungeEmailVerifiedAddressError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.ExpungeEmailVerifiedAddressRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeEmailVerifiedAddress(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailRequestEmailVerification(t *testing.T) {
	var rfc822MessageId string = "rfc822MessageId-427623191"
	var expectedResponse = &mailpb.RequestEmailVerificationResponse{
		Rfc822MessageId: rfc822MessageId,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.RequestEmailVerificationRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.RequestEmailVerification(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailRequestEmailVerificationError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var request = &mailpb.RequestEmailVerificationRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.RequestEmailVerification(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailVerifyEmail(t *testing.T) {
	var name2 string = "name2-1052831874"
	var parent string = "parent-995424086"
	var address string = "address-1147692044"
	var deleted bool = false
	var expectedResponse = &mailpb.EmailVerifiedAddress{
		Name:    name2,
		Parent:  parent,
		Address: address,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var token string = "token110541305"
	var request = &mailpb.VerifyEmailRequest{
		Name:  formattedName,
		Token: token,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.VerifyEmail(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailVerifyEmailError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/emailVerifiedAddresses/%s", "[PROJECT]", "[REGION]", "[EMAIL_VERIFIED_ADDRESS]")
	var token string = "token110541305"
	var request = &mailpb.VerifyEmailRequest{
		Name:  formattedName,
		Token: token,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.VerifyEmail(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailListSenders(t *testing.T) {
	var nextPageToken string = ""
	var sendersElement *mailpb.Sender = &mailpb.Sender{}
	var senders = []*mailpb.Sender{sendersElement}
	var expectedResponse = &mailpb.ListSendersResponse{
		NextPageToken: nextPageToken,
		Senders:       senders,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListSendersRequest{
		Parent:      formattedParent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListSenders(context.Background(), request).Next()

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	want := (interface{})(expectedResponse.Senders[0])
	got := (interface{})(resp)
	var ok bool

	switch want := (want).(type) {
	case proto.Message:
		ok = proto.Equal(want, got.(proto.Message))
	default:
		ok = want == got
	}
	if !ok {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListSendersError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListSendersRequest{
		Parent:      formattedParent,
		Region:      region,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListSenders(context.Background(), request).Next()

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetSender(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var defaultEnvelopeFromAuthority string = "defaultEnvelopeFromAuthority1550530879"
	var defaultHeaderFromAuthority string = "defaultHeaderFromAuthority-1184297630"
	var expectedResponse = &mailpb.Sender{
		Name:                         name2,
		Deleted:                      deleted,
		DefaultEnvelopeFromAuthority: defaultEnvelopeFromAuthority,
		DefaultHeaderFromAuthority:   defaultHeaderFromAuthority,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.GetSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.GetSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateSender(t *testing.T) {
	var name string = "name3373707"
	var deleted bool = false
	var defaultEnvelopeFromAuthority string = "defaultEnvelopeFromAuthority1550530879"
	var defaultHeaderFromAuthority string = "defaultHeaderFromAuthority-1184297630"
	var expectedResponse = &mailpb.Sender{
		Name:                         name,
		Deleted:                      deleted,
		DefaultEnvelopeFromAuthority: defaultEnvelopeFromAuthority,
		DefaultHeaderFromAuthority:   defaultHeaderFromAuthority,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var senderId string = "senderId32190309"
	var sender *mailpb.Sender = &mailpb.Sender{}
	var request = &mailpb.CreateSenderRequest{
		Parent:   formattedParent,
		Region:   region,
		SenderId: senderId,
		Sender:   sender,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s", "[PROJECT]")
	var region string = "region-934795532"
	var senderId string = "senderId32190309"
	var sender *mailpb.Sender = &mailpb.Sender{}
	var request = &mailpb.CreateSenderRequest{
		Parent:   formattedParent,
		Region:   region,
		SenderId: senderId,
		Sender:   sender,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateSender(t *testing.T) {
	var name string = "name3373707"
	var deleted bool = false
	var defaultEnvelopeFromAuthority string = "defaultEnvelopeFromAuthority1550530879"
	var defaultHeaderFromAuthority string = "defaultHeaderFromAuthority-1184297630"
	var expectedResponse = &mailpb.Sender{
		Name:                         name,
		Deleted:                      deleted,
		DefaultEnvelopeFromAuthority: defaultEnvelopeFromAuthority,
		DefaultHeaderFromAuthority:   defaultHeaderFromAuthority,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var sender *mailpb.Sender = &mailpb.Sender{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateSenderRequest{
		Sender:     sender,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var sender *mailpb.Sender = &mailpb.Sender{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateSenderRequest{
		Sender:     sender,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteSender(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var defaultEnvelopeFromAuthority string = "defaultEnvelopeFromAuthority1550530879"
	var defaultHeaderFromAuthority string = "defaultHeaderFromAuthority-1184297630"
	var expectedResponse = &mailpb.Sender{
		Name:                         name2,
		Deleted:                      deleted,
		DefaultEnvelopeFromAuthority: defaultEnvelopeFromAuthority,
		DefaultHeaderFromAuthority:   defaultHeaderFromAuthority,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.DeleteSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailDeleteSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.DeleteSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUndeleteSender(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var defaultEnvelopeFromAuthority string = "defaultEnvelopeFromAuthority1550530879"
	var defaultHeaderFromAuthority string = "defaultHeaderFromAuthority-1184297630"
	var expectedResponse = &mailpb.Sender{
		Name:                         name2,
		Deleted:                      deleted,
		DefaultEnvelopeFromAuthority: defaultEnvelopeFromAuthority,
		DefaultHeaderFromAuthority:   defaultHeaderFromAuthority,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.UndeleteSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUndeleteSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.UndeleteSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailExpungeSender(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.ExpungeSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeSender(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailExpungeSenderError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var request = &mailpb.ExpungeSenderRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeSender(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailSendMessage(t *testing.T) {
	var rfc822MessageId string = "rfc822MessageId-427623191"
	var expectedResponse = &mailpb.SendMessageResponse{
		Rfc822MessageId: rfc822MessageId,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedSender string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var envelopeFromAuthority string = "envelopeFromAuthority-735981251"
	var headerFromAuthority string = "headerFromAuthority-985559840"
	var envelopeFromAddress string = "envelopeFromAddress1388551278"
	var request = &mailpb.SendMessageRequest{
		Sender:                formattedSender,
		EnvelopeFromAuthority: envelopeFromAuthority,
		HeaderFromAuthority:   headerFromAuthority,
		EnvelopeFromAddress:   envelopeFromAddress,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.SendMessage(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailSendMessageError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedSender string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var envelopeFromAuthority string = "envelopeFromAuthority-735981251"
	var headerFromAuthority string = "headerFromAuthority-985559840"
	var envelopeFromAddress string = "envelopeFromAddress1388551278"
	var request = &mailpb.SendMessageRequest{
		Sender:                formattedSender,
		EnvelopeFromAuthority: envelopeFromAuthority,
		HeaderFromAuthority:   headerFromAuthority,
		EnvelopeFromAddress:   envelopeFromAddress,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.SendMessage(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailListSmtpCredentials(t *testing.T) {
	var expectedResponse *mailpb.ListSmtpCredentialsResponse = &mailpb.ListSmtpCredentialsResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var filter string = "filter-1274492040"
	var request = &mailpb.ListSmtpCredentialsRequest{
		Parent: formattedParent,
		Filter: filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListSmtpCredentials(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListSmtpCredentialsError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var filter string = "filter-1274492040"
	var request = &mailpb.ListSmtpCredentialsRequest{
		Parent: formattedParent,
		Filter: filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListSmtpCredentials(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetSmtpCredential(t *testing.T) {
	var name2 string = "name2-1052831874"
	var userId string = "userId-147132913"
	var password string = "password1216985755"
	var serviceAccountName string = "serviceAccountName235400871"
	var serviceAccountEmail string = "serviceAccountEmail-1300473088"
	var expectedResponse = &mailpb.SmtpCredential{
		Name:                name2,
		UserId:              userId,
		Password:            password,
		ServiceAccountName:  serviceAccountName,
		ServiceAccountEmail: serviceAccountEmail,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s/smtpCredentials/%s", "[PROJECT]", "[REGION]", "[SENDER]", "[SMTP_CREDENTIAL]")
	var request = &mailpb.GetSmtpCredentialRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetSmtpCredential(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetSmtpCredentialError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s/smtpCredentials/%s", "[PROJECT]", "[REGION]", "[SENDER]", "[SMTP_CREDENTIAL]")
	var request = &mailpb.GetSmtpCredentialRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetSmtpCredential(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateSmtpCredential(t *testing.T) {
	var name string = "name3373707"
	var userId string = "userId-147132913"
	var password string = "password1216985755"
	var serviceAccountName string = "serviceAccountName235400871"
	var serviceAccountEmail string = "serviceAccountEmail-1300473088"
	var expectedResponse = &mailpb.SmtpCredential{
		Name:                name,
		UserId:              userId,
		Password:            password,
		ServiceAccountName:  serviceAccountName,
		ServiceAccountEmail: serviceAccountEmail,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var smtpCredentialId string = "smtpCredentialId-1531115558"
	var smtpCredential *mailpb.SmtpCredential = &mailpb.SmtpCredential{}
	var request = &mailpb.CreateSmtpCredentialRequest{
		Parent:           formattedParent,
		SmtpCredentialId: smtpCredentialId,
		SmtpCredential:   smtpCredential,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateSmtpCredential(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateSmtpCredentialError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("projects/%s/regions/%s/senders/%s", "[PROJECT]", "[REGION]", "[SENDER]")
	var smtpCredentialId string = "smtpCredentialId-1531115558"
	var smtpCredential *mailpb.SmtpCredential = &mailpb.SmtpCredential{}
	var request = &mailpb.CreateSmtpCredentialRequest{
		Parent:           formattedParent,
		SmtpCredentialId: smtpCredentialId,
		SmtpCredential:   smtpCredential,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateSmtpCredential(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateSmtpCredential(t *testing.T) {
	var name string = "name3373707"
	var userId string = "userId-147132913"
	var password string = "password1216985755"
	var serviceAccountName string = "serviceAccountName235400871"
	var serviceAccountEmail string = "serviceAccountEmail-1300473088"
	var expectedResponse = &mailpb.SmtpCredential{
		Name:                name,
		UserId:              userId,
		Password:            password,
		ServiceAccountName:  serviceAccountName,
		ServiceAccountEmail: serviceAccountEmail,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var smtpCredential *mailpb.SmtpCredential = &mailpb.SmtpCredential{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateSmtpCredentialRequest{
		SmtpCredential: smtpCredential,
		UpdateMask:     updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateSmtpCredential(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateSmtpCredentialError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var smtpCredential *mailpb.SmtpCredential = &mailpb.SmtpCredential{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateSmtpCredentialRequest{
		SmtpCredential: smtpCredential,
		UpdateMask:     updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateSmtpCredential(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteSmtpCredential(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s/smtpCredentials/%s", "[PROJECT]", "[REGION]", "[SENDER]", "[SMTP_CREDENTIAL]")
	var request = &mailpb.DeleteSmtpCredentialRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.DeleteSmtpCredential(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailDeleteSmtpCredentialError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("projects/%s/regions/%s/senders/%s/smtpCredentials/%s", "[PROJECT]", "[REGION]", "[SENDER]", "[SMTP_CREDENTIAL]")
	var request = &mailpb.DeleteSmtpCredentialRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.DeleteSmtpCredential(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailListReceiptRules(t *testing.T) {
	var nextPageToken string = ""
	var totalSize int32 = 705419236
	var receiptRulesElement *mailpb.ReceiptRule = &mailpb.ReceiptRule{}
	var receiptRules = []*mailpb.ReceiptRule{receiptRulesElement}
	var expectedResponse = &mailpb.ListReceiptRulesResponse{
		NextPageToken: nextPageToken,
		TotalSize:     totalSize,
		ReceiptRules:  receiptRules,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var filter string = "filter-1274492040"
	var request = &mailpb.ListReceiptRulesRequest{
		Parent: formattedParent,
		Filter: filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListReceiptRules(context.Background(), request).Next()

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	want := (interface{})(expectedResponse.ReceiptRules[0])
	got := (interface{})(resp)
	var ok bool

	switch want := (want).(type) {
	case proto.Message:
		ok = proto.Equal(want, got.(proto.Message))
	default:
		ok = want == got
	}
	if !ok {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListReceiptRulesError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var filter string = "filter-1274492040"
	var request = &mailpb.ListReceiptRulesRequest{
		Parent: formattedParent,
		Filter: filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListReceiptRules(context.Background(), request).Next()

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetReceiptRule(t *testing.T) {
	var name2 string = "name2-1052831874"
	var expectedResponse = &mailpb.ReceiptRule{
		Name: name2,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", "[REGION]", "[DOMAIN]", "[RECEIPT_RULE]")
	var request = &mailpb.GetReceiptRuleRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetReceiptRule(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetReceiptRuleError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", "[REGION]", "[DOMAIN]", "[RECEIPT_RULE]")
	var request = &mailpb.GetReceiptRuleRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetReceiptRule(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateReceiptRule(t *testing.T) {
	var name string = "name3373707"
	var expectedResponse = &mailpb.ReceiptRule{
		Name: name,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var ruleId string = "ruleId1548659006"
	var receiptRule *mailpb.ReceiptRule = &mailpb.ReceiptRule{}
	var request = &mailpb.CreateReceiptRuleRequest{
		Parent:      formattedParent,
		RuleId:      ruleId,
		ReceiptRule: receiptRule,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateReceiptRule(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateReceiptRuleError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var ruleId string = "ruleId1548659006"
	var receiptRule *mailpb.ReceiptRule = &mailpb.ReceiptRule{}
	var request = &mailpb.CreateReceiptRuleRequest{
		Parent:      formattedParent,
		RuleId:      ruleId,
		ReceiptRule: receiptRule,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateReceiptRule(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateReceiptRule(t *testing.T) {
	var name string = "name3373707"
	var expectedResponse = &mailpb.ReceiptRule{
		Name: name,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var receiptRule *mailpb.ReceiptRule = &mailpb.ReceiptRule{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateReceiptRuleRequest{
		ReceiptRule: receiptRule,
		UpdateMask:  updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateReceiptRule(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateReceiptRuleError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var receiptRule *mailpb.ReceiptRule = &mailpb.ReceiptRule{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateReceiptRuleRequest{
		ReceiptRule: receiptRule,
		UpdateMask:  updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateReceiptRule(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteReceiptRule(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", "[REGION]", "[DOMAIN]", "[RECEIPT_RULE]")
	var request = &mailpb.DeleteReceiptRuleRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.DeleteReceiptRule(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailDeleteReceiptRuleError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/receiptRules/%s", "[REGION]", "[DOMAIN]", "[RECEIPT_RULE]")
	var request = &mailpb.DeleteReceiptRuleRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.DeleteReceiptRule(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailListAddressSets(t *testing.T) {
	var expectedResponse *mailpb.ListAddressSetsResponse = &mailpb.ListAddressSetsResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListAddressSetsRequest{
		Parent:      formattedParent,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListAddressSets(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailListAddressSetsError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var showDeleted bool = false
	var filter string = "filter-1274492040"
	var request = &mailpb.ListAddressSetsRequest{
		Parent:      formattedParent,
		ShowDeleted: showDeleted,
		Filter:      filter,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.ListAddressSets(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailGetAddressSet(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var expectedResponse = &mailpb.AddressSet{
		Name:    name2,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.GetAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.GetAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailCreateAddressSet(t *testing.T) {
	var name string = "name3373707"
	var deleted bool = false
	var expectedResponse = &mailpb.AddressSet{
		Name:    name,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var addressSetId string = "addressSetId549816515"
	var addressSet *mailpb.AddressSet = &mailpb.AddressSet{}
	var request = &mailpb.CreateAddressSetRequest{
		Parent:       formattedParent,
		AddressSetId: addressSetId,
		AddressSet:   addressSet,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailCreateAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedParent string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var addressSetId string = "addressSetId549816515"
	var addressSet *mailpb.AddressSet = &mailpb.AddressSet{}
	var request = &mailpb.CreateAddressSetRequest{
		Parent:       formattedParent,
		AddressSetId: addressSetId,
		AddressSet:   addressSet,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.CreateAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUpdateAddressSet(t *testing.T) {
	var name string = "name3373707"
	var deleted bool = false
	var expectedResponse = &mailpb.AddressSet{
		Name:    name,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var addressSet *mailpb.AddressSet = &mailpb.AddressSet{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateAddressSetRequest{
		AddressSet: addressSet,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUpdateAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var addressSet *mailpb.AddressSet = &mailpb.AddressSet{}
	var updateMask *field_maskpb.FieldMask = &field_maskpb.FieldMask{}
	var request = &mailpb.UpdateAddressSetRequest{
		AddressSet: addressSet,
		UpdateMask: updateMask,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailDeleteAddressSet(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var expectedResponse = &mailpb.AddressSet{
		Name:    name2,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.DeleteAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailDeleteAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.DeleteAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DeleteAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailUndeleteAddressSet(t *testing.T) {
	var name2 string = "name2-1052831874"
	var deleted bool = false
	var expectedResponse = &mailpb.AddressSet{
		Name:    name2,
		Deleted: deleted,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.UndeleteAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailUndeleteAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.UndeleteAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UndeleteAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailExpungeAddressSet(t *testing.T) {
	var expectedResponse *emptypb.Empty = &emptypb.Empty{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.ExpungeAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeAddressSet(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

}

func TestCloudMailExpungeAddressSetError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedName string = fmt.Sprintf("regions/%s/domains/%s/addressSets/%s", "[REGION]", "[DOMAIN]", "[ADDRESS_SET]")
	var request = &mailpb.ExpungeAddressSetRequest{
		Name: formattedName,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExpungeAddressSet(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
}
func TestCloudMailGetIamPolicy(t *testing.T) {
	var version int32 = 351608024
	var etag []byte = []byte("21")
	var expectedResponse = &iampb.Policy{
		Version: version,
		Etag:    etag,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &iampb.GetIamPolicyRequest{
		Resource: formattedResource,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetIamPolicy(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailGetIamPolicyError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var request = &iampb.GetIamPolicyRequest{
		Resource: formattedResource,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.GetIamPolicy(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailSetIamPolicy(t *testing.T) {
	var version int32 = 351608024
	var etag []byte = []byte("21")
	var expectedResponse = &iampb.Policy{
		Version: version,
		Etag:    etag,
	}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var policy *iampb.Policy = &iampb.Policy{}
	var request = &iampb.SetIamPolicyRequest{
		Resource: formattedResource,
		Policy:   policy,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.SetIamPolicy(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailSetIamPolicyError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var policy *iampb.Policy = &iampb.Policy{}
	var request = &iampb.SetIamPolicyRequest{
		Resource: formattedResource,
		Policy:   policy,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.SetIamPolicy(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
func TestCloudMailTestIamPermissions(t *testing.T) {
	var expectedResponse *iampb.TestIamPermissionsResponse = &iampb.TestIamPermissionsResponse{}

	mockCloudMail.err = nil
	mockCloudMail.reqs = nil

	mockCloudMail.resps = append(mockCloudMail.resps[:0], expectedResponse)

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var permissions []string = nil
	var request = &iampb.TestIamPermissionsRequest{
		Resource:    formattedResource,
		Permissions: permissions,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.TestIamPermissions(context.Background(), request)

	if err != nil {
		t.Fatal(err)
	}

	if want, got := request, mockCloudMail.reqs[0]; !proto.Equal(want, got) {
		t.Errorf("wrong request %q, want %q", got, want)
	}

	if want, got := expectedResponse, resp; !proto.Equal(want, got) {
		t.Errorf("wrong response %q, want %q)", got, want)
	}
}

func TestCloudMailTestIamPermissionsError(t *testing.T) {
	errCode := codes.PermissionDenied
	mockCloudMail.err = gstatus.Error(errCode, "test error")

	var formattedResource string = fmt.Sprintf("regions/%s/domains/%s", "[REGION]", "[DOMAIN]")
	var permissions []string = nil
	var request = &iampb.TestIamPermissionsRequest{
		Resource:    formattedResource,
		Permissions: permissions,
	}

	c, err := NewCloudMailClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.TestIamPermissions(context.Background(), request)

	if st, ok := gstatus.FromError(err); !ok {
		t.Errorf("got error %v, expected grpc error", err)
	} else if c := st.Code(); c != errCode {
		t.Errorf("got error code %q, want %q", c, errCode)
	}
	_ = resp
}
