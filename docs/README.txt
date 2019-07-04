The docs dir is the source of the GitHub Pages for knative/test-infra. See https://pages.github.com/ for details.

Main contents are rendered from README.md, which is a markdown+HTML file.

In order to allow README.md to read data from the oncall GCS bucket, proper permissions must be granted:

$ gsutil cors set cors-json-file.json gs://knative-infra-oncall

For more details, see https://cloud.google.com/storage/docs/configuring-cors

