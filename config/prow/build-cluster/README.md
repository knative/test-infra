## Create Build Cluster

`create-build-cluster.sh` is a wrapper of
[`GoogleCloudPlatform/oss-test-infra/create-build-cluster.sh`](https://github.com/GoogleCloudPlatform/oss-test-infra/blob/master/prow/create-build-cluster.sh),
with defined env vars overriding GCP project, cluster used by Knative.

Create build cluster:
```
# Creating build cluster, skip `Create a SA and secret for uploading results to GCS` as it's not needed
./create-build-cluster.sh

# Connect to knative-prow cluster
pushd ../
make get-cluster-credentials
popd

# Add the kubeconfig of newly created cluster to kubeconfig secret in Prow cluster
python3 "${GOPATH}/src/k8s.io/test-infra/gencred/merge_kubeconfig_secret.py" --src-key config --dest-key config build-cluster-kubeconfig.yaml
```

Copy secrets from old cluster to build cluster
```
# Copy over secrets needed for jobs
./add-secrets.sh
```
