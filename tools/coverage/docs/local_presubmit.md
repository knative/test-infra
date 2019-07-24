# Local Pre-submit Diff

A pre-submit check, to see how the coverage of your local repository compare to the latest successful postsubmit run, can be run locally in the following way

## Steps

1. `go get k8s.io/test-infra/robots/coverage`
  
2. Go to the target directory (where you want to run code coverage on), call the [local_presubmit.sh](../local_presubmit.sh) with the name of the Prow post-submit coverage job, e.g. `./local_presubmit.sh post-knative-serving-go-coverage`.
