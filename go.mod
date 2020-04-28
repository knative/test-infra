module knative.dev/test-infra

go 1.13

require (
	cloud.google.com/go/pubsub v1.2.0
	cloud.google.com/go/storage v1.6.0
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200123184029-53ce695e4179
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/cobra v0.0.6
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.20.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.3
	k8s.io/test-infra v0.0.0-20200407001919-bc7f71ef65b8 // indirect
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	knative.dev/pkg v0.0.0-20200406231459-466c676678dd
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4

	knative.dev/pkg => github.com/chizhg/pkg v0.0.0-20200427225849-8dd906c83e27
)
