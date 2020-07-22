module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.61.0 // indirect
	cloud.google.com/go/pubsub v1.5.0
	cloud.google.com/go/storage v1.10.0
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.1
	github.com/google/go-containerregistry v0.1.1
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/cobra v1.0.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200720211630-cb9d2d5c5666 // indirect
	golang.org/x/tools v0.0.0-20200721223218-6123e77877b2 // indirect
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200722002428-88e341933a54 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.6
	k8s.io/test-infra v0.0.0-20200722010006-526277bee528 // indirect
	sigs.k8s.io/boskos v0.0.0-20200717180850-7299d535c033
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
