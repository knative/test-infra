#!/bin/bash

function random_string {
  < /dev/urandom tr -dc A-Za-z0-9 | head -c14
}

kubectl create -f Kubernetes/env-variables.yaml
kubectl create secret generic grafana-pass --from-literal=grafana_admin_password="$(random_string)"
kubectl create secret generic postgres-pass --from-literal=pg_postgres_password="$(random_string)"
kubectl create secret generic gha-admin-pass --from-literal=pg_gha_admin_password="$(random_string)"
kubectl create -f Kubernetes/volumes.yaml
kubectl create -f Kubernetes/postgres-service.yaml
kubectl create -f Kubernetes/cli-home.yaml
