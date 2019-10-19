# perf-ops

perf-ops is a tool that maintains clusters running performance benchmarks in
each repo. 

> This tool has to be added in test-infra since we use `go run` to run it in
> `scripts/performance-tests.sh`. After the shell script is converted to pure Go
> code and moved out of `test-infra`, this tool can be removed.

## Basic Usage

Flags for this tool are:

flag.StringVar(&gcpProjectName, "gcp-project", "", "name of the GCP project for cluster operations")
	flag.StringVar(&repoName, "repository", "", "name of the repository")
	flag.StringVar(&benchmarkRootFolder, "benchmark-root", "", "root folder of the benchmarks")
	flag.BoolVar(&isRecreate, "recreate", false, "is recreate operation or not")
	flag.BoolVar(&isReconcile, "reconcile", false, "is reconcile operation or not")

- `--gcp-project` [Required] specifies the GCP project that holds the clusters.
- `--repository` [Required] specifies the Github repository of the performance benchmarks.
- `--benchmark-root` [Required] specifies root folder path of the benchmarks in
  the repo. It should contains and only contains all the benchmark folders.
- `--recreate` or `--reconcile` [Required] specifies operation for the
  benchmarks, only one of them can be specified.
