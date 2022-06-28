
We use argocd to sync cluster state with the configuration specified in knative/test-infra.

ArgoCD runs in the prow cluster and manages the following clusters:
- prow
- prow-build
- ~~prow-trusted~~


ArgoCD requires onetime bootstrap and occasional upgrades:
- Deploy the contents of bootstrap `k create ns argocd && k apply -k prow/cluster/gitops/bootstrap -n argocd `
- Download argocd cli
- `argocd login -core`
- `argocd cluster add gke_knative-tests_us-central1_prow-build`
- `argocd cluster add gke_knative-tests_us-central1_prow`
- `argocd cluster add gke_knative-tests_us-central1-a_prow-trusted`
- `k apply -f prow/cluster/gitops/bootstrap/apps.yaml`


Couple of notes:
- Argo UI is accessible at https://argo.knative.dev, but it is restricted to the Google Groups that have IAP-secured Web App User on the knative-tests project.
