module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.62.0 // indirect
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/storage v1.10.0
	github.com/blang/semver/v4 v4.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-git/go-git-fixtures/v4 v4.0.1
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.1
	github.com/google/go-containerregistry v0.1.4
	github.com/google/go-github/v27 v27.0.6
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200802091954-4b90ce9b60b3 // indirect
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.8
	k8s.io/test-infra v0.0.0-20200803112140-d8aa4e063646 // indirect
	knative.dev/hack v0.0.0-20201125230335-c46a6498e9ed
	sigs.k8s.io/boskos v0.0.0-20200729174948-794df80db9c9
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
)
