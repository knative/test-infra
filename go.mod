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
	github.com/google/go-containerregistry v0.10.0
	github.com/google/go-github/v32 v32.1.1-0.20201004213705-76c3c3d7c6e7 // HEAD as of Nov 6
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.5.0
	go.uber.org/atomic v1.9.0
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2
	google.golang.org/api v0.86.0
	k8s.io/apimachinery v0.24.2
	k8s.io/test-infra v0.0.0-20220801075428-527a7b720677 // indirect
	knative.dev/hack v0.0.0-20220725145124-782bbaabb8a1
	sigs.k8s.io/boskos v0.0.0-20220516053617-d007e4465eb6
	sigs.k8s.io/yaml v1.3.0
)

// COPIED FROM kubernetes/test-infra

// Pin all k8s.io staging repositories to kubernetes v0.18.6
// When bumping Kubernetes dependencies, you should update each of these lines
// to point to the same kubernetes v0.KubernetesMinor.KubernetesPatch version
// before running update-deps.sh.
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190301231843-5614ed5bae6f

	// Upstream is unmaintained. This fork introduces two important changes:
	// * We log an error if writing a cache key fails (e.G. because disk is full)
	// * We inject a header that allows ghproxy to detect if the response was revalidated or a cache miss
	github.com/gregjones/httpcache => github.com/alvaroaleman/httpcache v0.0.0-20210618195546-ab9a1a3f8a38

	golang.org/x/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20190709130402-674ba3eaed22
	k8s.io/client-go => k8s.io/client-go v0.23.9 // consider the version used by knative/pkg repo
)

require (
	cloud.google.com/go v0.102.1 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/docker/cli v20.10.16+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.16+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go4.org v0.0.0-20201209231011-d4a079459e60 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220706132729-d86698d07c53 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/utils v0.0.0-20220725171434-9bab9ef40391 // indirect
)
