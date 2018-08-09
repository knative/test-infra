FROM golang:1.10.2
RUN apt-get update && apt-get install -y --no-install-recommends

COPY apicoverage /apicoverage
ENTRYPOINT ["/apicoverage"]