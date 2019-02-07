A pre-submit check, to see how the coverage of your local repository compare to the lastest postsubmit run, can be run locally in the following way

1. download https://github.com/kubernetes/test-infra/tree/master/robots/coverage

2. build the binary
  - `go build -o coverage`
  
3. run the binary to download coverage profile from the latest healthy build of the given job
  - `./coverage download knative-prow [prowjob-name] -o $BASE_PROFILE`
  - for more help, run `./coverage download`
  
4. run `go test -coverprofile` on your local repository to produce the coverage profile. Assuming you store it in location $NEW_PROFILE

5. compare the new profile with the base profile
  - `./coverage diff $BASE_PROFILE $NEW_PROFILE`
