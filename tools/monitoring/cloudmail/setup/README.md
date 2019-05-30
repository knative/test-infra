# Cloud Mail Setup

## `setup.go`

`setup.go` sets up the cloud mail domain with the knative monitoring configuration. The setup actions
only need to be done once.

There are three actions in `setup.go`:

1. `setup-domain`: Sets up the email domain. This prints out the domain ID and domain name.
The domain id and the region together uniquely identifies the cloud mail domain. The domain
name is used to send the emails.
    ```bash
    go run setup.go -setup-domain
    ```

1. `setup-sender`: Sets up a sender mail address. It requires the `-domainID` to be set and create
all the resources under the specified domain.
   ```bash
   go run setup.go -setup-sender -domain-id "<DomainID generated in (1)"
   ```

1. `send-test-mail`: Send a test mail from the resources created in the previous setup steps
   ```bash
   go run setup.go -send-test-mail -domain-name "<Domain Name generated in (1)>" -to-address "<recipient email>"
   ```

## Set up Go Client Library

Follow the instructions at [Quickstart go client library setup](https://cloud.google.com/mail/docs/quickstart-client-libraries#cloud-mail-client-libraries-go)

## Cloud Mail Command Line Utility

For more complex usage of cloud mail, use the [gcloud command line interface](https://cloud.google.com/mail/docs/quickstart-cli)
directly. It is required to [sign up](https://goo.gl/UC8Eb4) for cloud mail alpha before the command line utility can be used.
