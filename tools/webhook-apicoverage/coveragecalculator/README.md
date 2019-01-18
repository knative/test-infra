# Coverage Calculator

`coveragecalculator` package contains types and helper methods pertaining to
coverage calculation. Package includes type [TypeCoverage](coveragedata.go)
to represent coverage data for a particular API object type. This is the
wire contract between the webhook server running inside the K8 cluster and any
client using the API-Coverage tool. All API calls into the webhook-server would
return response containing this object to represent coverage data.