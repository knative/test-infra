# Test

This directory contains tests and testing docs.

- [Test library](#test-library) contains code you can use in your `knative`
  tests
- [Flags](#flags) added by [the test library](#test-library)
- [Unit tests](#running-unit-tests) currently reside in the codebase alongside
  the code they test

## Running unit tests

To run all unit tests:

```bash
go test ./...
```

## Test library

You can use the test library in this dir to:

- [Use common test flags](#use-common-test-flags)
- [Output logs](#output-logs)
- [Emit metrics](#emit-metrics)
- [Ensure test cleanup](#ensure-test-cleanup)

### Use common test flags

These flags are useful for running against an existing cluster, making use of
your existing
[environment setup](https://github.com/knative/serving/blob/master/DEVELOPMENT.md#environment-setup).

By importing `github.com/knative/pkg/test` you get access to a global variable
called `test.Flags` which holds the values of
[the command line flags](/test/README.md#flags).

```go
logger.Infof("Using namespace %s", test.Flags.Namespace)
```

_See [e2e_flags.go](./e2e_flags.go)._

### Output logs

[When tests are run with `--logverbose` option](README.md#output-verbose-logs),
debug logs will be emitted to stdout.

We are using the common [e2e logging library](logging/logging.go) that uses the
[Knative logging library](../logging/) for structured logging. It is built on
top of [zap](https://github.com/uber-go/zap). Tests should initialize the global
logger to use a test specifc context with `logging.GetContextLogger`:

```go
// The convention is for the name of the logger to match the name of the test.
logging.GetContextLogger("TestHelloWorld")
```

Logs can then be emitted using the `logger` object which is required by many
functions in the test library. To emit logs:

```go
logger.Infof("Creating a new Route %s and Configuration %s", route, configuration)
logger.Debugf("The LogURL is %s, not yet verifying", logURL)
```

_See [logging.go](./logging/logging.go)._

### Emit metrics

You can emit metrics from your tests using
[the opencensus library](https://github.com/census-instrumentation/opencensus-go),
which
[is being used inside Knative as well](https://github.com/knative/serving/blob/master/docs/telemetry.md).
These metrics will be emitted by the test if the test is run with
[the `--emitmetrics` option](#metrics-flag).

You can record arbitrary metrics with
[`stats.Record`](https://github.com/census-instrumentation/opencensus-go#stats)
or measure latency with
[`trace.StartSpan`](https://github.com/census-instrumentation/opencensus-go#traces):

```go
ctx, span := trace.StartSpan(context.Background(), "MyMetric")
```

- These traces will be emitted automatically by
  [the generic crd polling functions](#check-knative-serving-resources).
- The traces are emitted by [a custom metric exporter](./logging/logging.go)
  that uses the global logger instance.

#### Metric format

When a `trace` metric is emitted, the format is
`metric <name> <startTime> <endTime> <duration>`. The name of the metric is
arbitrary and can be any string. The values are:

- `metric` - Indicates this log is a metric
- `<name>` - Arbitrary string indentifying the metric
- `<startTime>` - Unix time in nanoseconds when measurement started
- `<endTime>` - Unix time in nanoseconds when measurement ended
- `<duration>` - The difference in ms between the startTime and endTime

For example:

```bash
metric WaitForConfigurationState/prodxiparjxt/ConfigurationUpdatedWithRevision 1529980772357637397 1529980772431586609 73.949212ms
```

_The [`Wait` methods](#check-knative-serving-resources) (which poll resources)
will prefix the metric names with the name of the function, and if applicable,
the name of the resource, separated by `/`. In the example above,
`WaitForConfigurationState` is the name of the function, and `prodxiparjxt` is
the name of the configuration resource being polled.
`ConfigurationUpdatedWithRevision` is the string passed to
`WaitForConfigurationState` by the caller to identify what state is being polled
for._

### Check Knative Serving resources

_WARNING: this code also exists in
[`knative/serving`](https://github.com/knative/serving/blob/master/test/adding_tests.md#make-requests-against-deployed-services)._

After creating Knative Serving resources or making changes to them, you will
need to wait for the system to realize those changes. You can use the Knative
Serving CRD check and polling methods to check the resources are either in or
reach the desired state.

The `WaitFor*` functions use the kubernetes
[`wait` package](https://godoc.org/k8s.io/apimachinery/pkg/util/wait). To poll
they use
[`PollImmediate`](https://godoc.org/k8s.io/apimachinery/pkg/util/wait#PollImmediate)
and the return values of the function you provide behave the same as
[`ConditionFunc`](https://godoc.org/k8s.io/apimachinery/pkg/util/wait#ConditionFunc):
a `bool` to indicate if the function should stop or continue polling, and an
`error` to indicate if there has been an error.

For example, you can poll a `Configuration` object to find the name of the
`Revision` that was created for it:

```go
var revisionName string
err := test.WaitForConfigurationState(clients.ServingClient, configName, func(c *v1alpha1.Configuration) (bool, error) {
    if c.Status.LatestCreatedRevisionName != "" {
        revisionName = c.Status.LatestCreatedRevisionName
        return true, nil
    }
    return false, nil
}, "ConfigurationUpdatedWithRevision")
```

_[Metrics will be emitted](#emit-metrics) for these `Wait` method tracking how
long test poll for._

_See [kube_checks.go](./kube_checks.go)._

### Ensure test cleanup

To ensure your test is cleaned up, you should defer cleanup to execute after
your test completes and also ensure the cleanup occurs if the test is
interrupted:

```go
defer tearDown(clients)
test.CleanupOnInterrupt(func() { tearDown(clients) })
```

_See [cleanup.go](./simple_e2e_test.go)._

## Flags

Importing [the test library](#test-library) adds flags that are useful for end
to end tests that need to run against a cluster.

Tests importing [`github.com/knative/pkg/test`](#test-library) recognize these
flags:

- [`--kubeconfig`](#specifying-kubeconfig)
- [`--cluster`](#specifying-cluster)
- [`--namespace`](#specifying-namespace)
- [`--logverbose`](#output-verbose-logs)
- [`--emitmetrics`](#metrics-flag)

### Specifying kubeconfig

By default the tests will use the
[kubeconfig file](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)
at `~/.kube/config`. If there is an error getting the current user, it will use
`kubeconfig` instead as the default value. You can specify a different config
file with the argument `--kubeconfig`.

To run tests with a non-default kubeconfig file:

```bash
go test ./test --kubeconfig /my/path/kubeconfig
```

### Specifying cluster

The `--cluster` argument lets you use a different cluster than
[your specified kubeconfig's](#specifying-kubeconfig) active context. This will
default to the value of your
[`K8S_CLUSTER_OVERRIDE` environment variable](https://github.com/knative/serving/blob/master/DEVELOPMENT.md#environment-setup)
if not specified.

```bash
go test ./test --cluster your-cluster-name
```

The current cluster names can be obtained by running:

```bash
kubectl config get-clusters
```

### Specifying namespace

The `--namespace` argument lets you specify the namespace to use for the tests.
By default, tests will use `serving-tests`.

```bash
go test ./test --namespace your-namespace-name
```

### Output verbose logs

The `--logverbose` argument lets you see verbose test logs and k8s logs.

```bash
go test ./test --logverbose
```

### Metrics flag

Running tests with the `--emitmetrics` argument will cause latency metrics to be
emitted by the tests.

```bash
go test ./test --emitmetrics
```

- To add additional metrics to a test, see
  [emitting metrics](adding_tests.md#emit-metrics).
- For more info on the format of the metrics, see
  [metric format](adding_tests.md#metric-format).

[minikube]: https://kubernetes.io/docs/setup/minikube/

---

Except as otherwise noted, the content of this page is licensed under the
[Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/),
and code samples are licensed under the
[Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).
# Hello World - Elixir Sample

A simple web application written in [Elixir](https://elixir-lang.org/) using the
[Phoenix Framework](https://phoenixframework.org/). The application prints all
environment variables to the main page.

# Set up Elixir and Phoenix Locally

Following the
[Phoenix Installation Guide](https://hexdocs.pm/phoenix/installation.html) is
the best way to get your computer set up for developing, building, running, and
packaging Elixir Web applications.

# Running Locally

To start your Phoenix server:

- Install dependencies with `mix deps.get`
- Install Node.js dependencies with `cd assets && npm install`
- Start Phoenix endpoint with `mix phx.server`

Now you can visit [`localhost:4000`](http://localhost:4000) from your browser.

# Recreating the sample code

1. Generate a new project.

```shell
mix phoenix.new helloelixir
```

When asked, if you want to `Fetch and install dependencies? [Yn]` select `y`

1. Follow the direction in the output to change directories into start your
   local server with `mix phoenix.server`

1. In the new directory, create a new Dockerfile for packaging your application
   for deployment

   ```docker
   # Start from a base image for elixir
   # Phoenix works best on pre 1.7 at the moment.
   FROM elixir:1.6.6-alpine

   # Set up Elixir and Phoenix
   ARG APP_NAME=hello
   ARG PHOENIX_SUBDIR=.
   ENV MIX_ENV=prod REPLACE_OS_VARS=true TERM=xterm
   WORKDIR /opt/app

   # Update nodejs, rebar, and hex.
   RUN apk update \
       && apk --no-cache --update add nodejs nodejs-npm \
       && mix local.rebar --force \
       && mix local.hex --force
   COPY . .

   # Download and compile dependencies, then compile Web app.
   RUN mix do deps.get, deps.compile, compile
   RUN cd ${PHOENIX_SUBDIR}/assets \
       && npm install \
       && ./node_modules/brunch/bin/brunch build -p \
       && cd .. \
       && mix phx.digest

   # Create a release version of the application
   RUN mix release --env=prod --verbose \
       && mv _build/prod/rel/${APP_NAME} /opt/release \
       && mv /opt/release/bin/${APP_NAME} /opt/release/bin/start_server

   # Prepare final layer
   FROM alpine:latest
   RUN apk update && apk --no-cache --update add bash openssl-dev

   # Add a user so the server will run as a non-root user.
   RUN addgroup -g 1000 appuser && \
       adduser -S -u 1000 -G appuser appuser
   # Pre-create necessary temp directory for erlang and set permissions.
   RUN mkdir -p /opt/app/var
   RUN chown appuser /opt/app/var
   # Run everything else as 'appuser'
   USER appuser

   # Document that the service listens on port 8080.
   ENV PORT=8080 MIX_ENV=prod REPLACE_OS_VARS=true
   WORKDIR /opt/app
   EXPOSE 8080
   COPY --from=0 /opt/release .
   ENV RUNNER_LOG_DIR /var/log

   # Command to execute the application.
   CMD ["/opt/app/bin/start_server", "foreground", "boot_var=/tmp"]
   ```

1. Create a new file, `service.yaml` and copy the following Service definition
   into the file. Make sure to replace `{username}` with your Docker Hub
   username.

   ```yaml
   apiVersion: serving.knative.dev/v1alpha1
   kind: Service
   metadata:
     name: helloworld-elixir
     namespace: default
   spec:
     runLatest:
       configuration:
         revisionTemplate:
           spec:
             container:
               image: docker.io/{username}/helloworld-elixir
               env:
                 - name: TARGET
                   value: "elixir Sample v1"
   ```

# Building and deploying the sample

The sample in this directory is ready to build and deploy without changes. You
can deploy the sample as is, or use you created version following the directions
above.

1.  Generate a new `secret_key_base` in the `config/prod.secret.exs` file.
    Phoenix applications use a secrets file on production deployments and, by
    default, that file is not checked into source control. We have provides
    shell of an example on `config/prod.secret.exs.sample` and you can use the
    following command to generate a new prod secrets file.

    ```shell
    SECRET_KEY_BASE=$(elixir -e ":crypto.strong_rand_bytes(48) |> Base.encode64 |> IO.puts")
    sed "s|SECRET+KEY+BASE|$SECRET_KEY_BASE|" config/prod.secret.exs.sample >config/prod.secret.exs
    ```

1.  Use Docker to build the sample code into a container. To build and push with
    Docker Hub, run these commands replacing `{username}` with your Docker Hub
    username:

    ```shell
     # Build the container on your local machine
     docker build -t {username}/helloworld-elixir .

     # Push the container to docker registry
     docker push {username}/helloworld-elixir
    ```

1.  After the build has completed and the container is pushed to docker hub, you
    can deploy the app into your cluster. Ensure that the container image value
    in `service.yaml` matches the container you built in the previous step.
    Apply the configuration using `kubectl`:

    ```shell
    kubectl apply --filename service.yaml
    ```

1.  Now that your service is created, Knative will perform the following steps:

    - Create a new immutable revision for this version of the app.
    - Network programming to create a route, ingress, service, and load balance
      for your app.
    - Automatically scale your pods up and down (including to zero active pods).

1.  To find the IP address for your service, use
    `kubectl get svc knative-ingressgateway --namespace istio-system` to get the
    ingress IP for your cluster. If your cluster is new, it may take sometime
    for the service to get asssigned an external IP address.

        ```
        kubectl get svc knative-ingressgateway --namespace istio-system

        NAME                     TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                                      AGE

    knative-ingressgateway LoadBalancer 10.35.254.218 35.225.171.32
    80:32380/TCP,443:32390/TCP,32400:32400/TCP 1h

    ```

    ```

1.  To find the URL for your service, use

    ```
    kubectl get ksvc helloworld-elixir --output=custom-columns=NAME:.metadata.name,DOMAIN:.status.domain

    NAME                DOMAIN
    helloworld-elixir   helloworld-elixir.default.example.com
    ```

1.  Now you can make a request to your app to see the results. Replace
    `{IP_ADDRESS}` with the address you see returned in the previous step.

        ```shell
        curl -H "Host: helloworld-elixir.default.example.com" http://{IP_ADDRESS}

        ...
        # HTML from your application is returned.
        ```

    Here is the HTML returned from our deployed sample application:

    ```HTML
    <!DOCTYPE html>
    <html lang="en">
    <head>
     <meta charset="utf-8">
     <meta http-equiv="X-UA-Compatible" content="IE=edge">
     <meta name="viewport" content="width=device-width, initial-scale=1">
     <meta name="description" content="">
     <meta name="author" content="">

     <title>Hello Knative</title>
     <link rel="stylesheet" type="text/css" href="/css/app-833cc7e8eeed7a7953c5a02e28130dbd.css?vsn=d">
    </head>
    ```

  <body>
    <div class="container">
      <header class="header">
        <nav role="navigation">

        </nav>
      </header>

      <p class="alert alert-info" role="alert"></p>
      <p class="alert alert-danger" role="alert"></p>

      <main role="main">

<div class="jumbotron">
  <h2>Welcome to Knative and Elixir</h2>

  <p>$TARGET = elixir Sample v1</p>
</div>

  <h3>Environment</h3>
  <ul>
    <li>BINDIR = /opt/app/erts-9.3.2/bin</li>
    <li>DEST_SYS_CONFIG_PATH = /opt/app/var/sys.config</li>
    <li>DEST_VMARGS_PATH = /opt/app/var/vm.args</li>
    <li>DISTILLERY_TASK = foreground</li>
    <li>EMU = beam</li>
    <li>ERL_LIBS = /opt/app/lib</li>
    <li>ERL_OPTS = </li>
    <li>ERTS_DIR = /opt/app/erts-9.3.2</li>
    <li>ERTS_LIB_DIR = /opt/app/erts-9.3.2/../lib</li>
    <li>ERTS_VSN = 9.3.2</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_PORT = tcp://10.35.241.50:80</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_PORT_80_TCP = tcp://10.35.241.50:80</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_PORT_80_TCP_ADDR = 10.35.241.50</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_PORT_80_TCP_PORT = 80</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_PORT_80_TCP_PROTO = tcp</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_SERVICE_HOST = 10.35.241.50</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_SERVICE_PORT = 80</li>
    <li>HELLOWORLD_ELIXIR_00001_SERVICE_SERVICE_PORT_HTTP = 80</li>
    <li>HELLOWORLD_ELIXIR_PORT = tcp://10.35.253.90:80</li>
    <li>HELLOWORLD_ELIXIR_PORT_80_TCP = tcp://10.35.253.90:80</li>
    <li>HELLOWORLD_ELIXIR_PORT_80_TCP_ADDR = 10.35.253.90</li>
    <li>HELLOWORLD_ELIXIR_PORT_80_TCP_PORT = 80</li>
    <li>HELLOWORLD_ELIXIR_PORT_80_TCP_PROTO = tcp</li>
    <li>HELLOWORLD_ELIXIR_SERVICE_HOST = 10.35.253.90</li>
    <li>HELLOWORLD_ELIXIR_SERVICE_PORT = 80</li>
    <li>HELLOWORLD_ELIXIR_SERVICE_PORT_HTTP = 80</li>
    <li>HOME = /root</li>
    <li>HOSTNAME = helloworld-elixir-00001-deployment-84f68946b4-76hcv</li>
    <li>KUBERNETES_PORT = tcp://10.35.240.1:443</li>
    <li>KUBERNETES_PORT_443_TCP = tcp://10.35.240.1:443</li>
    <li>KUBERNETES_PORT_443_TCP_ADDR = 10.35.240.1</li>
    <li>KUBERNETES_PORT_443_TCP_PORT = 443</li>
    <li>KUBERNETES_PORT_443_TCP_PROTO = tcp</li>
    <li>KUBERNETES_SERVICE_HOST = 10.35.240.1</li>
    <li>KUBERNETES_SERVICE_PORT = 443</li>
    <li>KUBERNETES_SERVICE_PORT_HTTPS = 443</li>
    <li>LD_LIBRARY_PATH = /opt/app/erts-9.3.2/lib:</li>
    <li>MIX_ENV = prod</li>
    <li>NAME = hello@127.0.0.1</li>
    <li>NAME_ARG = -name hello@127.0.0.1</li>
    <li>NAME_TYPE = -name</li>
    <li>OLDPWD = /opt/app</li>
    <li>OTP_VER = 20</li>
    <li>PATH = /opt/app/erts-9.3.2/bin:/opt/app/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin</li>
    <li>PORT = 8080</li>
    <li>PROGNAME = opt/app/releases/0.0.1/hello.sh</li>
    <li>PWD = /opt/app</li>
    <li>RELEASES_DIR = /opt/app/releases</li>
    <li>RELEASE_CONFIG_DIR = /opt/app</li>
    <li>RELEASE_ROOT_DIR = /opt/app</li>
    <li>REL_NAME = hello</li>
    <li>REL_VSN = 0.0.1</li>
    <li>REPLACE_OS_VARS = true</li>
    <li>ROOTDIR = /opt/app</li>
    <li>RUNNER_LOG_DIR = /var/log</li>
    <li>RUN_ERL_ENV = </li>
    <li>SHLVL = 1</li>
    <li>SRC_SYS_CONFIG_PATH = /opt/app/releases/0.0.1/sys.config</li>
    <li>SRC_VMARGS_PATH = /opt/app/releases/0.0.1/vm.args</li>
    <li>SYS_CONFIG_PATH = /opt/app/var/sys.config</li>
    <li>TARGET = elixir Sample v1</li>
    <li>TERM = xterm</li>
    <li>VMARGS_PATH = /opt/app/var/vm.args</li>
  </ul>
      </main>

    </div> <!-- /container -->
    <script src="/js/app-930ab1950e10d7b5ab5083423c28f06e.js?vsn=d"></script>

  </body>
</html>
   ```

## Removing the sample app deployment

To remove the sample app from your cluster, delete the service record:

```shell
kubectl delete --filename service.yaml
```
