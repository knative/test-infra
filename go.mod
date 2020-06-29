module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.55.0 // indirect
	cloud.google.com/go/pubsub v1.2.0
	cloud.google.com/go/storage v1.6.0
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200123184029-53ce695e4179
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/cobra v0.0.6
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200327173247-9dae0f8f5775 // indirect
	golang.org/x/tools v0.0.0-20200329025819-fd4102a86c65 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200326112834-f447254575fd // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.6
	sigs.k8s.io/boskos v0.0.0-20200530174753-71e795271860
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
