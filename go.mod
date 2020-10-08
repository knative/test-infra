module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.68.0 // indirect
	cloud.google.com/go/pubsub v1.8.1
	cloud.google.com/go/storage v1.12.0
	github.com/blang/semver/v4 v4.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-git/go-git-fixtures/v4 v4.0.1
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.2
	github.com/google/go-containerregistry v0.1.3
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.0.0
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/net v0.0.0-20201006153459-a7d1128ccaa0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sync v0.0.0-20200930132711-30421366ff76 // indirect
	golang.org/x/sys v0.0.0-20201008064518-c1f3e3309c71 // indirect
	golang.org/x/tools v0.0.0-20201008025239-9df69603baec // indirect
	google.golang.org/api v0.32.0
	google.golang.org/genproto v0.0.0-20201007142714-5c0e72c5e71e // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.19.2
	k8s.io/test-infra v0.0.0-20201008101616-21f5a5568c79 // indirect
	sigs.k8s.io/boskos v0.0.0-20201002225104-ae3497d24cd7
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
