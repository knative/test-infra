# Observed concurrency
This tool is designed to show the latency numbers calculated by the perf tests.

## Read from GCS
This tool reads the logs from the latest performance test run of knative/serving. 
It uses the service account passed in or by default will use the GOOGLE_APPLICATION_CREDENTIALS variable to get the logs. 

## Creating Output
This tool creates an output xml in the prow artifacts directory. The prow artifacts directory is passed in or by default will use `./artifacts` directory.

This output xml will be read by testgrid and displayed on the [dashboard](https://testgrid.knative.dev/knative-serving#perf-latency).

## Prow Job
There is a daily prow job that triggers this tool that is run at 02:05 AM PST. This tool will then generate the output xml which is then displayed in the testgrid dashboard.
