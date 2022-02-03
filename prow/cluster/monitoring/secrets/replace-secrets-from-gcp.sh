#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

PROJECT="${1:-}"
if [[ -z "${PROJECT}" ]]; then
    echo "ERROR: GCP project name must be provided."
    exit 1
fi

echo "Updating secret with secret manager from project $PROJECT"
email="$(gcloud beta secrets versions access latest --secret prometheus-alerter-email-address --project ${PROJECT})"
token="$(gcloud beta secrets versions access latest --secret prometheus-alerter-app-password --project ${PROJECT})"

cat "${CUR_DIR}/alertmanager-prow_secret.yaml" | \
    sed s/{{\ smtp_app_username\ }}/${email}/g | \
    sed s/{{\ smtp_app_password\ }}/${token}/g | \
    kubectl apply -f -

echo "Secrets updated"
