module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go/pubsub v1.8.2
	cloud.google.com/go/storage v1.12.0
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v20.10.0-beta1+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.2
	github.com/google/go-containerregistry v0.1.4
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20201027133719-8eef5233e2a1
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20201028094953-708e7fb298ac // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20201028111035-eafbe7b904eb // indirect
	google.golang.org/api v0.34.0
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201028140639-c77dae4b0522 // indirect
	google.golang.org/grpc v1.33.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools/v3 v3.0.3 // indirect
	k8s.io/apimachinery v0.19.3
	k8s.io/test-infra v0.0.0-20201028132156-1e7ec95bcb40 // indirect
	sigs.k8s.io/boskos v0.0.0-20201002225104-ae3497d24cd7
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
