module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go/pubsub v1.2.0
	cloud.google.com/go/storage v1.6.0
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200123184029-53ce695e4179
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20190926221455-842c0d70d702
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.20.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.3
	k8s.io/test-infra v0.0.0-20200405022946-02a6dffa9047 // indirect
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	knative.dev/pkg v0.0.0-20200404181734-92cdec5b3593
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
)
