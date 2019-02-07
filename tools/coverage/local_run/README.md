# Local Pre-submit Diff

A pre-submit check, to see how the coverage of your local repository compare to the latest successful postsubmit run, can be run locally in the following way

## Steps

1. `go get github.com/kubernetes/test-infra/robots/coverage`

2. build the binary
  - `go build -o coverage`
  
3. go to the target directory (where you want to run code coverage on), call the local_presubmit.sh with the following two arguments: 
  - name of the prow job. E.g. post-knative-serving-go-coverage
  - path the the binary built in step 2
