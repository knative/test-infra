module knative.dev/test-infra

go 1.18

require (
	cloud.google.com/go/storage v1.22.1
	github.com/blang/semver/v4 v4.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-git/go-git-fixtures/v4 v4.2.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.8
	github.com/google/go-containerregistry v0.8.1-0.20220216220642-00c59d91847c
	github.com/google/go-github/v32 v32.1.1-0.20201004213705-76c3c3d7c6e7 // HEAD as of Nov 6
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.4.0
	go.uber.org/atomic v1.9.0
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	golang.org/x/net v0.0.0-20220607020251-c690dde0001d
	golang.org/x/oauth2 v0.0.0-20220608161450-d0670ef3b1eb
	google.golang.org/api v0.83.0
	k8s.io/apimachinery v0.24.2
	knative.dev/hack v0.0.0-20220902140220-97fd663bf52b
	sigs.k8s.io/boskos v0.0.0-20210730172138-093b54882439
	sigs.k8s.io/yaml v1.3.0
)

// TODO: remove this once we have a fixed knative-hack (https://github.com/knative/hack/issues/216)
replace knative.dev/hack => github.com/cardil/knative-hack v0.0.0-20220902210110-5d40dd4a1f82

require (
	cloud.google.com/go v0.102.0 // indirect
	cloud.google.com/go/compute v1.6.1 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/docker/cli v20.10.12+incompatible // indirect
	github.com/docker/distribution v2.8.0+incompatible // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go4.org v0.0.0-20201209231011-d4a079459e60 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220608164250-635b8c9b7f68 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220608133413-ed9918b62aac // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.1.0 // indirect
	k8s.io/test-infra v0.0.0-20220729194712-07f065e00e9b // indirect
	k8s.io/utils v0.0.0-20220725171434-9bab9ef40391 // indirect
)
