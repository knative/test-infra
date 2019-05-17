## Cloud Mail Setup

### `setup.go`

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

### Cloud Mail Command Line Utility

For more complex usage of cloud mail, use the command line interface directly. It is required
to sign up for [cloud mail alpha](go/cloud-mail-alpha-form) before the command line utility
can be used.

##### Useful Commands

List resources (domains, sender, address-set, receipt rules)
```bash
gcloud alpha mail [domains|senders] list --location="<region>"
gcloud alpha mail domains [address-sets|receipt-rules] list  --location="<region>" --domain="<domainId>"
```

Delete resources (domains, sender, address-set, receipt rules)

```bash
gcloud alpha mail domains delete <domainid> --location="<region>"
gcloud alpha mail senders delete <senderid> --location="<region>"
gcloud alpha mail domains address-sets delete <address-set-id>  --location="<region>" --domain="<domainId>"
gcloud alpha mail domains receipt-rules delete <receipt-rule-id>  ---location="<region>" --domain="<domainId>"
```
