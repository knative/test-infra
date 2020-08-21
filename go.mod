module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.62.0 // indirect
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/storage v1.10.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.1
	github.com/google/go-containerregistry v0.1.1
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200802091954-4b90ce9b60b3 // indirect
	golang.org/x/tools v0.0.0-20200731060945-b5fad4ed8dd6 // indirect
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/test-infra v0.0.0-20200803112140-d8aa4e063646 // indirect
	sigs.k8s.io/boskos v0.0.0-20200729174948-794df80db9c9
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
