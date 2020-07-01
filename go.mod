module knative.dev/test-infra

go 1.14

require (
	cloud.google.com/go v0.60.0 // indirect
	cloud.google.com/go/pubsub v1.4.0
	cloud.google.com/go/storage v1.10.0
	github.com/go-logr/logr v0.2.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.0
	github.com/google/go-containerregistry v0.1.1
	github.com/google/go-github/v27 v27.0.6
	github.com/google/licenseclassifier v0.0.0-20200402202327-879cb1424de0
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/tools v0.0.0-20200701000337-a32c0cb1d5b2 // indirect
	google.golang.org/api v0.28.0
	google.golang.org/genproto v0.0.0-20200701001935-0939c5918c31 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.5 // indirect
	k8s.io/apimachinery v0.18.5
	k8s.io/klog/v2 v2.2.0 // indirect
	k8s.io/test-infra v0.0.0-20200630233406-1dca6122872e // indirect
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
	knative.dev/pkg v0.0.0-20200630170034-2c1a029eb97f
	sigs.k8s.io/boskos v0.0.0-20200617235605-f289ba6555ba // indirect
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.0.0-20190815212128-ab0dd09aa10e

	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.0.0
)
