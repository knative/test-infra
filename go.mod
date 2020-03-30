module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go/pubsub v1.2.0
	cloud.google.com/go/storage v1.6.0
	github.com/docker/cli v0.0.0-20191105005515-99c5edceb48d // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200123184029-53ce695e4179
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/licenseclassifier v0.0.0-20181010185715-e979a0b10eeb
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.20.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.0
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	knative.dev/pkg v0.0.0-20200323231609-0840da9555a3
)

replace knative.dev/pkg => github.com/chizhg/pkg v0.0.0-20200330020211-4643096970a8
