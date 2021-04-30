module knative.dev/test-infra

go 1.15

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
	github.com/google/go-github/v32 v32.1.1-0.20201004213705-76c3c3d7c6e7 // HEAD as of Nov 6
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	go.uber.org/atomic v1.6.0
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.19.7
	knative.dev/hack v0.0.0-20210428122153-93ad9129c268
	sigs.k8s.io/boskos v0.0.0-20200729174948-794df80db9c9
)
