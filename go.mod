module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.60.0 // indirect
	cloud.google.com/go/pubsub v1.4.0
	cloud.google.com/go/storage v1.10.0
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.0
	github.com/google/go-containerregistry v0.1.1
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/cobra v1.0.0
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200701000337-a32c0cb1d5b2 // indirect
	google.golang.org/api v0.28.0
	google.golang.org/genproto v0.0.0-20200701001935-0939c5918c31 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.5
	k8s.io/test-infra v0.0.0-20200630233406-1dca6122872e // indirect
	sigs.k8s.io/boskos v0.0.0-20200617235605-f289ba6555ba
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
