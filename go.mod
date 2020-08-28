module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.65.0 // indirect
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/storage v1.11.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.2
	github.com/google/go-containerregistry v0.1.2
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200828081204-131dc92a58d5 // indirect
	golang.org/x/tools v0.0.0-20200828013309-97019fc2e64b // indirect
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200828030656-73b5761be4c5 // indirect
	google.golang.org/grpc v1.31.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.19.0
	k8s.io/test-infra v0.0.0-20200828131253-b23899a92dfa // indirect
	sigs.k8s.io/boskos v0.0.0-20200819010710-984516eae7e8
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
