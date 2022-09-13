module knative.dev/test-infra

go 1.15

require (
	cloud.google.com/go/iam v0.1.0 // indirect
	cloud.google.com/go/storage v1.10.0
	github.com/blang/semver/v4 v4.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-git/go-git-fixtures/v4 v4.0.1
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.7
	github.com/google/go-containerregistry v0.8.0
	github.com/google/go-github/v32 v32.1.1-0.20201004213705-76c3c3d7c6e7 // HEAD as of Nov 6
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.7.0
	github.com/wavesoftware/go-commandline v1.0.0
	go.uber.org/atomic v1.7.0
	golang.org/x/mod v0.6.0-dev.0.20220818022119-ed83ed61efb9
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	google.golang.org/api v0.70.0
	google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf // indirect
	k8s.io/apimachinery v0.20.6
	knative.dev/hack v0.0.0-20220913095247-7556452c2b54
	sigs.k8s.io/boskos v0.0.0-20200729174948-794df80db9c9
	sigs.k8s.io/yaml v1.2.0
)
