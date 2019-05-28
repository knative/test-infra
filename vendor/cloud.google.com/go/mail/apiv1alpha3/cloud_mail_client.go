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
	"time"

	"cloud.google.com/go/internal/version"
	"github.com/golang/protobuf/proto"
	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	mailpb "google.golang.org/genproto/googleapis/cloud/mail/v1alpha3"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// CloudMailCallOptions contains the retry settings for each method of CloudMailClient.
type CloudMailCallOptions struct {
	ListDomains                  []gax.CallOption
	GetDomain                    []gax.CallOption
	CreateDomain                 []gax.CallOption
	UpdateDomain                 []gax.CallOption
	DeleteDomain                 []gax.CallOption
	UndeleteDomain               []gax.CallOption
	ExpungeDomain                []gax.CallOption
	TestReceiptRules             []gax.CallOption
	VerifyDomain                 []gax.CallOption
	ListEmailVerifiedAddresses   []gax.CallOption
	GetEmailVerifiedAddress      []gax.CallOption
	CreateEmailVerifiedAddress   []gax.CallOption
	UpdateEmailVerifiedAddress   []gax.CallOption
	DeleteEmailVerifiedAddress   []gax.CallOption
	UndeleteEmailVerifiedAddress []gax.CallOption
	ExpungeEmailVerifiedAddress  []gax.CallOption
	RequestEmailVerification     []gax.CallOption
	VerifyEmail                  []gax.CallOption
	ListSenders                  []gax.CallOption
	GetSender                    []gax.CallOption
	CreateSender                 []gax.CallOption
	UpdateSender                 []gax.CallOption
	DeleteSender                 []gax.CallOption
	UndeleteSender               []gax.CallOption
	ExpungeSender                []gax.CallOption
	SendMessage                  []gax.CallOption
	ListSmtpCredentials          []gax.CallOption
	GetSmtpCredential            []gax.CallOption
	CreateSmtpCredential         []gax.CallOption
	UpdateSmtpCredential         []gax.CallOption
	DeleteSmtpCredential         []gax.CallOption
	ListReceiptRules             []gax.CallOption
	GetReceiptRule               []gax.CallOption
	CreateReceiptRule            []gax.CallOption
	UpdateReceiptRule            []gax.CallOption
	DeleteReceiptRule            []gax.CallOption
	ListAddressSets              []gax.CallOption
	GetAddressSet                []gax.CallOption
	CreateAddressSet             []gax.CallOption
	UpdateAddressSet             []gax.CallOption
	DeleteAddressSet             []gax.CallOption
	UndeleteAddressSet           []gax.CallOption
	ExpungeAddressSet            []gax.CallOption
	GetIamPolicy                 []gax.CallOption
	SetIamPolicy                 []gax.CallOption
	TestIamPermissions           []gax.CallOption
}

func defaultCloudMailClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("cloudmail.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultCloudMailCallOptions() *CloudMailCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &CloudMailCallOptions{
		ListDomains:                  retry[[2]string{"default", "idempotent"}],
		GetDomain:                    retry[[2]string{"default", "idempotent"}],
		CreateDomain:                 retry[[2]string{"default", "non_idempotent"}],
		UpdateDomain:                 retry[[2]string{"default", "non_idempotent"}],
		DeleteDomain:                 retry[[2]string{"default", "idempotent"}],
		UndeleteDomain:               retry[[2]string{"default", "non_idempotent"}],
		ExpungeDomain:                retry[[2]string{"default", "non_idempotent"}],
		TestReceiptRules:             retry[[2]string{"default", "non_idempotent"}],
		VerifyDomain:                 retry[[2]string{"default", "non_idempotent"}],
		ListEmailVerifiedAddresses:   retry[[2]string{"default", "idempotent"}],
		GetEmailVerifiedAddress:      retry[[2]string{"default", "idempotent"}],
		CreateEmailVerifiedAddress:   retry[[2]string{"default", "non_idempotent"}],
		UpdateEmailVerifiedAddress:   retry[[2]string{"default", "non_idempotent"}],
		DeleteEmailVerifiedAddress:   retry[[2]string{"default", "idempotent"}],
		UndeleteEmailVerifiedAddress: retry[[2]string{"default", "non_idempotent"}],
		ExpungeEmailVerifiedAddress:  retry[[2]string{"default", "non_idempotent"}],
		RequestEmailVerification:     retry[[2]string{"default", "non_idempotent"}],
		VerifyEmail:                  retry[[2]string{"default", "non_idempotent"}],
		ListSenders:                  retry[[2]string{"default", "idempotent"}],
		GetSender:                    retry[[2]string{"default", "idempotent"}],
		CreateSender:                 retry[[2]string{"default", "non_idempotent"}],
		UpdateSender:                 retry[[2]string{"default", "non_idempotent"}],
		DeleteSender:                 retry[[2]string{"default", "idempotent"}],
		UndeleteSender:               retry[[2]string{"default", "non_idempotent"}],
		ExpungeSender:                retry[[2]string{"default", "non_idempotent"}],
		SendMessage:                  retry[[2]string{"default", "non_idempotent"}],
		ListSmtpCredentials:          retry[[2]string{"default", "idempotent"}],
		GetSmtpCredential:            retry[[2]string{"default", "idempotent"}],
		CreateSmtpCredential:         retry[[2]string{"default", "non_idempotent"}],
		UpdateSmtpCredential:         retry[[2]string{"default", "non_idempotent"}],
		DeleteSmtpCredential:         retry[[2]string{"default", "idempotent"}],
		ListReceiptRules:             retry[[2]string{"default", "idempotent"}],
		GetReceiptRule:               retry[[2]string{"default", "idempotent"}],
		CreateReceiptRule:            retry[[2]string{"default", "non_idempotent"}],
		UpdateReceiptRule:            retry[[2]string{"default", "non_idempotent"}],
		DeleteReceiptRule:            retry[[2]string{"default", "idempotent"}],
		ListAddressSets:              retry[[2]string{"default", "idempotent"}],
		GetAddressSet:                retry[[2]string{"default", "idempotent"}],
		CreateAddressSet:             retry[[2]string{"default", "non_idempotent"}],
		UpdateAddressSet:             retry[[2]string{"default", "non_idempotent"}],
		DeleteAddressSet:             retry[[2]string{"default", "idempotent"}],
		UndeleteAddressSet:           retry[[2]string{"default", "non_idempotent"}],
		ExpungeAddressSet:            retry[[2]string{"default", "non_idempotent"}],
		GetIamPolicy:                 retry[[2]string{"default", "idempotent"}],
		SetIamPolicy:                 retry[[2]string{"default", "non_idempotent"}],
		TestIamPermissions:           retry[[2]string{"default", "non_idempotent"}],
	}
}

// CloudMailClient is a client for interacting with Google Cloud Mail API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type CloudMailClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	cloudMailClient mailpb.CloudMailClient

	// The call options for this service.
	CallOptions *CloudMailCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewCloudMailClient creates a new cloud mail client.
//
// Provides Google Cloud Mail customers a way to manage the service and to send
// mail.  Use this to add or remove domains, define handling for incoming
// messages, and register senders for outbound messages.
func NewCloudMailClient(ctx context.Context, opts ...option.ClientOption) (*CloudMailClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultCloudMailClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &CloudMailClient{
		conn:        conn,
		CallOptions: defaultCloudMailCallOptions(),

		cloudMailClient: mailpb.NewCloudMailClient(conn),
	}
	c.setGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *CloudMailClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *CloudMailClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *CloudMailClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", version.Go()}, keyval...)
	kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// ListDomains lists domains with the given parent.
func (c *CloudMailClient) ListDomains(ctx context.Context, req *mailpb.ListDomainsRequest, opts ...gax.CallOption) (*mailpb.ListDomainsResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListDomains[0:len(c.CallOptions.ListDomains):len(c.CallOptions.ListDomains)], opts...)
	var resp *mailpb.ListDomainsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.ListDomains(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetDomain gets the specified domain.
func (c *CloudMailClient) GetDomain(ctx context.Context, req *mailpb.GetDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetDomain[0:len(c.CallOptions.GetDomain):len(c.CallOptions.GetDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateDomain registers the specified domain for Cloud Mail.
// Cloudmail can provide a regional, per-project domain name e.g.
// my-project.us-east1.cloudsmtp.net
// where cloudmail manages all of the dns. These can be
// created by calling CreateDomain with Domain.project_domain == true
// and Domain.domain_name == "".
func (c *CloudMailClient) CreateDomain(ctx context.Context, req *mailpb.CreateDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateDomain[0:len(c.CallOptions.CreateDomain):len(c.CallOptions.CreateDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateDomain updates the given domain.
func (c *CloudMailClient) UpdateDomain(ctx context.Context, req *mailpb.UpdateDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateDomain[0:len(c.CallOptions.UpdateDomain):len(c.CallOptions.UpdateDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteDomain marks a domain as deleted. It will be automatically expunged after 30 days
// unless it is undeleted with UndeleteDomain.
func (c *CloudMailClient) DeleteDomain(ctx context.Context, req *mailpb.DeleteDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteDomain[0:len(c.CallOptions.DeleteDomain):len(c.CallOptions.DeleteDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.DeleteDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UndeleteDomain removes the deleted status for a domain that was previously deleted with
// DeleteDomain.
func (c *CloudMailClient) UndeleteDomain(ctx context.Context, req *mailpb.UndeleteDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UndeleteDomain[0:len(c.CallOptions.UndeleteDomain):len(c.CallOptions.UndeleteDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UndeleteDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExpungeDomain permanently expunges a domain. Will only succeed on domains already marked
// deleted using the DeleteDomain call.
func (c *CloudMailClient) ExpungeDomain(ctx context.Context, req *mailpb.ExpungeDomainRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ExpungeDomain[0:len(c.CallOptions.ExpungeDomain):len(c.CallOptions.ExpungeDomain)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.ExpungeDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// TestReceiptRules evaluates a recipient address against the domain's receipt ruleset and
// returns the list of rules that would fire.  Clients may provide an optional
// alternative candidate ruleset to be evaluated instead of the service's
// active ruleset.  This method can be used to verify Cloud Mail behavior for
// incoming messages.
func (c *CloudMailClient) TestReceiptRules(ctx context.Context, req *mailpb.TestReceiptRulesRequest, opts ...gax.CallOption) (*mailpb.TestReceiptRulesResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.TestReceiptRules[0:len(c.CallOptions.TestReceiptRules):len(c.CallOptions.TestReceiptRules)], opts...)
	var resp *mailpb.TestReceiptRulesResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.TestReceiptRules(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyDomain checks the domain's DNS TXT record for the verification token, and updates
// the status to ACTIVE if valid.
func (c *CloudMailClient) VerifyDomain(ctx context.Context, req *mailpb.VerifyDomainRequest, opts ...gax.CallOption) (*mailpb.Domain, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.VerifyDomain[0:len(c.CallOptions.VerifyDomain):len(c.CallOptions.VerifyDomain)], opts...)
	var resp *mailpb.Domain
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.VerifyDomain(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListEmailVerifiedAddresses lists EmailVerifiedAddresses with the given parent.
func (c *CloudMailClient) ListEmailVerifiedAddresses(ctx context.Context, req *mailpb.ListEmailVerifiedAddressesRequest, opts ...gax.CallOption) (*mailpb.ListEmailVerifiedAddressesResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListEmailVerifiedAddresses[0:len(c.CallOptions.ListEmailVerifiedAddresses):len(c.CallOptions.ListEmailVerifiedAddresses)], opts...)
	var resp *mailpb.ListEmailVerifiedAddressesResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.ListEmailVerifiedAddresses(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetEmailVerifiedAddress gets the specified EmailVerifiedAddress.
func (c *CloudMailClient) GetEmailVerifiedAddress(ctx context.Context, req *mailpb.GetEmailVerifiedAddressRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetEmailVerifiedAddress[0:len(c.CallOptions.GetEmailVerifiedAddress):len(c.CallOptions.GetEmailVerifiedAddress)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateEmailVerifiedAddress creates the given EmailVerifiedAddress.
func (c *CloudMailClient) CreateEmailVerifiedAddress(ctx context.Context, req *mailpb.CreateEmailVerifiedAddressRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateEmailVerifiedAddress[0:len(c.CallOptions.CreateEmailVerifiedAddress):len(c.CallOptions.CreateEmailVerifiedAddress)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateEmailVerifiedAddress updates the given EmailVerifiedAddress.
func (c *CloudMailClient) UpdateEmailVerifiedAddress(ctx context.Context, req *mailpb.UpdateEmailVerifiedAddressRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateEmailVerifiedAddress[0:len(c.CallOptions.UpdateEmailVerifiedAddress):len(c.CallOptions.UpdateEmailVerifiedAddress)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteEmailVerifiedAddress marks the specified EmailVerifiedAddress as deleted. It will be
// automatically expunged after 30 days unless it is undeleted with
// UndeleteEmailVerifiedAddress.
func (c *CloudMailClient) DeleteEmailVerifiedAddress(ctx context.Context, req *mailpb.DeleteEmailVerifiedAddressRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteEmailVerifiedAddress[0:len(c.CallOptions.DeleteEmailVerifiedAddress):len(c.CallOptions.DeleteEmailVerifiedAddress)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.DeleteEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UndeleteEmailVerifiedAddress undeletes the specified EmailVerifiedAddress.
func (c *CloudMailClient) UndeleteEmailVerifiedAddress(ctx context.Context, req *mailpb.UndeleteEmailVerifiedAddressRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UndeleteEmailVerifiedAddress[0:len(c.CallOptions.UndeleteEmailVerifiedAddress):len(c.CallOptions.UndeleteEmailVerifiedAddress)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UndeleteEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExpungeEmailVerifiedAddress permanently expunges an EmailVerifiedAddress. Will only succeed on
// resources already marked deleted using the DeleteEmailVerifiedAddress call.
func (c *CloudMailClient) ExpungeEmailVerifiedAddress(ctx context.Context, req *mailpb.ExpungeEmailVerifiedAddressRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ExpungeEmailVerifiedAddress[0:len(c.CallOptions.ExpungeEmailVerifiedAddress):len(c.CallOptions.ExpungeEmailVerifiedAddress)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.ExpungeEmailVerifiedAddress(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// RequestEmailVerification emails a verification token to an unverified EmailVerifiedAddress.
func (c *CloudMailClient) RequestEmailVerification(ctx context.Context, req *mailpb.RequestEmailVerificationRequest, opts ...gax.CallOption) (*mailpb.RequestEmailVerificationResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.RequestEmailVerification[0:len(c.CallOptions.RequestEmailVerification):len(c.CallOptions.RequestEmailVerification)], opts...)
	var resp *mailpb.RequestEmailVerificationResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.RequestEmailVerification(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyEmail checks token and verifies EmailVerifiedAddress
func (c *CloudMailClient) VerifyEmail(ctx context.Context, req *mailpb.VerifyEmailRequest, opts ...gax.CallOption) (*mailpb.EmailVerifiedAddress, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.VerifyEmail[0:len(c.CallOptions.VerifyEmail):len(c.CallOptions.VerifyEmail)], opts...)
	var resp *mailpb.EmailVerifiedAddress
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.VerifyEmail(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListSenders lists senders for the given parent.
func (c *CloudMailClient) ListSenders(ctx context.Context, req *mailpb.ListSendersRequest, opts ...gax.CallOption) *SenderIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListSenders[0:len(c.CallOptions.ListSenders):len(c.CallOptions.ListSenders)], opts...)
	it := &SenderIterator{}
	req = proto.Clone(req).(*mailpb.ListSendersRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*mailpb.Sender, string, error) {
		var resp *mailpb.ListSendersResponse
		req.PageToken = pageToken
		//if pageSize > math.MaxInt32 {
		//	req.PageSize = math.MaxInt32
		//} else {
		//	req.PageSize = int32(pageSize)
		//}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.cloudMailClient.ListSenders(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Senders, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	//it.pageInfo.MaxSize = int(req.PageSize)
	return it
}

// GetSender gets the specified sender.
func (c *CloudMailClient) GetSender(ctx context.Context, req *mailpb.GetSenderRequest, opts ...gax.CallOption) (*mailpb.Sender, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetSender[0:len(c.CallOptions.GetSender):len(c.CallOptions.GetSender)], opts...)
	var resp *mailpb.Sender
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateSender creates the specified sender.
func (c *CloudMailClient) CreateSender(ctx context.Context, req *mailpb.CreateSenderRequest, opts ...gax.CallOption) (*mailpb.Sender, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateSender[0:len(c.CallOptions.CreateSender):len(c.CallOptions.CreateSender)], opts...)
	var resp *mailpb.Sender
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateSender updates the specified sender.
func (c *CloudMailClient) UpdateSender(ctx context.Context, req *mailpb.UpdateSenderRequest, opts ...gax.CallOption) (*mailpb.Sender, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateSender[0:len(c.CallOptions.UpdateSender):len(c.CallOptions.UpdateSender)], opts...)
	var resp *mailpb.Sender
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteSender marks the specified sender as deleted. It will be automatically expunged
// after 30 days unless it is undeleted with UndeleteSender.
func (c *CloudMailClient) DeleteSender(ctx context.Context, req *mailpb.DeleteSenderRequest, opts ...gax.CallOption) (*mailpb.Sender, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteSender[0:len(c.CallOptions.DeleteSender):len(c.CallOptions.DeleteSender)], opts...)
	var resp *mailpb.Sender
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.DeleteSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UndeleteSender undeletes the specified sender.
func (c *CloudMailClient) UndeleteSender(ctx context.Context, req *mailpb.UndeleteSenderRequest, opts ...gax.CallOption) (*mailpb.Sender, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UndeleteSender[0:len(c.CallOptions.UndeleteSender):len(c.CallOptions.UndeleteSender)], opts...)
	var resp *mailpb.Sender
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UndeleteSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExpungeSender permanently expunges a Sender. Will only succeed on resources already
// marked deleted using the DeleteSender call.
func (c *CloudMailClient) ExpungeSender(ctx context.Context, req *mailpb.ExpungeSenderRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ExpungeSender[0:len(c.CallOptions.ExpungeSender):len(c.CallOptions.ExpungeSender)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.ExpungeSender(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// SendMessage sends a message using the specified sender.  The "From" address in the
// message headers must be a registered and verified domain with the service,
// and it must also match the sender's list of allowed "From" patterns;
// otherwise, the request will fail with a FAILED_PRECONDITION error.
func (c *CloudMailClient) SendMessage(ctx context.Context, req *mailpb.SendMessageRequest, opts ...gax.CallOption) (*mailpb.SendMessageResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.SendMessage[0:len(c.CallOptions.SendMessage):len(c.CallOptions.SendMessage)], opts...)
	var resp *mailpb.SendMessageResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.SendMessage(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListSmtpCredentials lists SMTP credentials for the specified sender.
func (c *CloudMailClient) ListSmtpCredentials(ctx context.Context, req *mailpb.ListSmtpCredentialsRequest, opts ...gax.CallOption) (*mailpb.ListSmtpCredentialsResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListSmtpCredentials[0:len(c.CallOptions.ListSmtpCredentials):len(c.CallOptions.ListSmtpCredentials)], opts...)
	var resp *mailpb.ListSmtpCredentialsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.ListSmtpCredentials(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSmtpCredential gets the specified SMTP credential.
func (c *CloudMailClient) GetSmtpCredential(ctx context.Context, req *mailpb.GetSmtpCredentialRequest, opts ...gax.CallOption) (*mailpb.SmtpCredential, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetSmtpCredential[0:len(c.CallOptions.GetSmtpCredential):len(c.CallOptions.GetSmtpCredential)], opts...)
	var resp *mailpb.SmtpCredential
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetSmtpCredential(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateSmtpCredential creates the specified SMTP credential.
func (c *CloudMailClient) CreateSmtpCredential(ctx context.Context, req *mailpb.CreateSmtpCredentialRequest, opts ...gax.CallOption) (*mailpb.SmtpCredential, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateSmtpCredential[0:len(c.CallOptions.CreateSmtpCredential):len(c.CallOptions.CreateSmtpCredential)], opts...)
	var resp *mailpb.SmtpCredential
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateSmtpCredential(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateSmtpCredential updates the specified SMTP credential.
func (c *CloudMailClient) UpdateSmtpCredential(ctx context.Context, req *mailpb.UpdateSmtpCredentialRequest, opts ...gax.CallOption) (*mailpb.SmtpCredential, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateSmtpCredential[0:len(c.CallOptions.UpdateSmtpCredential):len(c.CallOptions.UpdateSmtpCredential)], opts...)
	var resp *mailpb.SmtpCredential
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateSmtpCredential(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteSmtpCredential deletes the specified SMTP credential.
func (c *CloudMailClient) DeleteSmtpCredential(ctx context.Context, req *mailpb.DeleteSmtpCredentialRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteSmtpCredential[0:len(c.CallOptions.DeleteSmtpCredential):len(c.CallOptions.DeleteSmtpCredential)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.DeleteSmtpCredential(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListReceiptRules lists receipt rules for the specified Cloud Mail domain.
func (c *CloudMailClient) ListReceiptRules(ctx context.Context, req *mailpb.ListReceiptRulesRequest, opts ...gax.CallOption) *ReceiptRuleIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListReceiptRules[0:len(c.CallOptions.ListReceiptRules):len(c.CallOptions.ListReceiptRules)], opts...)
	it := &ReceiptRuleIterator{}
	req = proto.Clone(req).(*mailpb.ListReceiptRulesRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*mailpb.ReceiptRule, string, error) {
		var resp *mailpb.ListReceiptRulesResponse
		req.PageToken = pageToken
		//if pageSize > math.MaxInt32 {
		//	req.PageSize = math.MaxInt32
		//} else {
		//	req.PageSize = int32(pageSize)
		//}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.cloudMailClient.ListReceiptRules(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.ReceiptRules, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	//it.pageInfo.MaxSize = int(req.PageSize)
	return it
}

// GetReceiptRule gets the specified receipt rule.
func (c *CloudMailClient) GetReceiptRule(ctx context.Context, req *mailpb.GetReceiptRuleRequest, opts ...gax.CallOption) (*mailpb.ReceiptRule, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetReceiptRule[0:len(c.CallOptions.GetReceiptRule):len(c.CallOptions.GetReceiptRule)], opts...)
	var resp *mailpb.ReceiptRule
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetReceiptRule(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateReceiptRule creates the specified receipt rule.
func (c *CloudMailClient) CreateReceiptRule(ctx context.Context, req *mailpb.CreateReceiptRuleRequest, opts ...gax.CallOption) (*mailpb.ReceiptRule, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateReceiptRule[0:len(c.CallOptions.CreateReceiptRule):len(c.CallOptions.CreateReceiptRule)], opts...)
	var resp *mailpb.ReceiptRule
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateReceiptRule(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateReceiptRule updates the specified receipt rule.
func (c *CloudMailClient) UpdateReceiptRule(ctx context.Context, req *mailpb.UpdateReceiptRuleRequest, opts ...gax.CallOption) (*mailpb.ReceiptRule, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateReceiptRule[0:len(c.CallOptions.UpdateReceiptRule):len(c.CallOptions.UpdateReceiptRule)], opts...)
	var resp *mailpb.ReceiptRule
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateReceiptRule(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteReceiptRule deletes the specified receipt rule.  If the rule is part of the domain's
// active ruleset, the rule reference is also removed from the ruleset.
func (c *CloudMailClient) DeleteReceiptRule(ctx context.Context, req *mailpb.DeleteReceiptRuleRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteReceiptRule[0:len(c.CallOptions.DeleteReceiptRule):len(c.CallOptions.DeleteReceiptRule)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.DeleteReceiptRule(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListAddressSets lists AddressSets for the specified Cloud Mail domain.
func (c *CloudMailClient) ListAddressSets(ctx context.Context, req *mailpb.ListAddressSetsRequest, opts ...gax.CallOption) (*mailpb.ListAddressSetsResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListAddressSets[0:len(c.CallOptions.ListAddressSets):len(c.CallOptions.ListAddressSets)], opts...)
	var resp *mailpb.ListAddressSetsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.ListAddressSets(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAddressSet gets the specified AddressSet.
func (c *CloudMailClient) GetAddressSet(ctx context.Context, req *mailpb.GetAddressSetRequest, opts ...gax.CallOption) (*mailpb.AddressSet, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetAddressSet[0:len(c.CallOptions.GetAddressSet):len(c.CallOptions.GetAddressSet)], opts...)
	var resp *mailpb.AddressSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateAddressSet creates the specified AddressSet.
func (c *CloudMailClient) CreateAddressSet(ctx context.Context, req *mailpb.CreateAddressSetRequest, opts ...gax.CallOption) (*mailpb.AddressSet, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateAddressSet[0:len(c.CallOptions.CreateAddressSet):len(c.CallOptions.CreateAddressSet)], opts...)
	var resp *mailpb.AddressSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.CreateAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateAddressSet updates the specified AddressSet.
func (c *CloudMailClient) UpdateAddressSet(ctx context.Context, req *mailpb.UpdateAddressSetRequest, opts ...gax.CallOption) (*mailpb.AddressSet, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateAddressSet[0:len(c.CallOptions.UpdateAddressSet):len(c.CallOptions.UpdateAddressSet)], opts...)
	var resp *mailpb.AddressSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UpdateAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteAddressSet marks the specified AddressSet as deleted. It will be automatically
// expunged after 30 days unless it is undeleted with UndeleteAddressSet.
func (c *CloudMailClient) DeleteAddressSet(ctx context.Context, req *mailpb.DeleteAddressSetRequest, opts ...gax.CallOption) (*mailpb.AddressSet, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteAddressSet[0:len(c.CallOptions.DeleteAddressSet):len(c.CallOptions.DeleteAddressSet)], opts...)
	var resp *mailpb.AddressSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.DeleteAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UndeleteAddressSet undeletes the specified AddressSet.
func (c *CloudMailClient) UndeleteAddressSet(ctx context.Context, req *mailpb.UndeleteAddressSetRequest, opts ...gax.CallOption) (*mailpb.AddressSet, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UndeleteAddressSet[0:len(c.CallOptions.UndeleteAddressSet):len(c.CallOptions.UndeleteAddressSet)], opts...)
	var resp *mailpb.AddressSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.UndeleteAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExpungeAddressSet permanently expunges an AddressSet. Will only succeed on resources
// already marked deleted using the DeleteAddressSet call.
func (c *CloudMailClient) ExpungeAddressSet(ctx context.Context, req *mailpb.ExpungeAddressSetRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ExpungeAddressSet[0:len(c.CallOptions.ExpungeAddressSet):len(c.CallOptions.ExpungeAddressSet)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.cloudMailClient.ExpungeAddressSet(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// GetIamPolicy gets the access control policy for Cloud Mail resources.
// Returns an empty policy if the resource exists and does not have a policy
// set.
func (c *CloudMailClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetIamPolicy[0:len(c.CallOptions.GetIamPolicy):len(c.CallOptions.GetIamPolicy)], opts...)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.GetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetIamPolicy sets the access control policy for a Cloud Mail Resources. Replaces
// any existing policy.
func (c *CloudMailClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.SetIamPolicy[0:len(c.CallOptions.SetIamPolicy):len(c.CallOptions.SetIamPolicy)], opts...)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.SetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TestIamPermissions returns permissions that a caller has on a Cloud Mail Resource.
// If the resource does not exist, this will return an empty set of
// permissions, not a [NOT_FOUND][google.rpc.Code.NOT_FOUND] error.
//
// Note: This operation is designed to be used for building permission-aware
// UIs and command-line tools, not for authorization checking. This operation
// may "fail open" without warning.
func (c *CloudMailClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.TestIamPermissions[0:len(c.CallOptions.TestIamPermissions):len(c.CallOptions.TestIamPermissions)], opts...)
	var resp *iampb.TestIamPermissionsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.cloudMailClient.TestIamPermissions(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ReceiptRuleIterator manages a stream of *mailpb.ReceiptRule.
type ReceiptRuleIterator struct {
	items    []*mailpb.ReceiptRule
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*mailpb.ReceiptRule, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ReceiptRuleIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ReceiptRuleIterator) Next() (*mailpb.ReceiptRule, error) {
	var item *mailpb.ReceiptRule
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ReceiptRuleIterator) bufLen() int {
	return len(it.items)
}

func (it *ReceiptRuleIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// SenderIterator manages a stream of *mailpb.Sender.
type SenderIterator struct {
	items    []*mailpb.Sender
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*mailpb.Sender, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *SenderIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *SenderIterator) Next() (*mailpb.Sender, error) {
	var item *mailpb.Sender
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *SenderIterator) bufLen() int {
	return len(it.items)
}

func (it *SenderIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
