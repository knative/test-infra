# knative.dev redirect app

This is a simple redirection app for redirecting subdomains in `knative.dev`. It
expects the destination URL to be set in the `REDIR_TO` environment variable.
Each subdomain has its own `*subdomain*.knative.dev.yaml` deployment file.

The `Makefile` contains ready-to-use rules to update the projects in the
`knative-dns` GCP project.

The `dispatch.yaml` file defines dispatching rules for all apps hosted in the
`knative.dev` subdomain.
