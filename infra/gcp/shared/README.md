# Knative Shared Project

This project runs a shared gke cluster for members of the project to run workloads on.

Pay *extra* attention when modifying GKE Clusters as others are running workloads on it.

Key Infrastructure:
- shared cluster

Getting access to this project:
- We are leveraging [Google Groups for RBAC]() so make a sure a group called k8s-infra-rbac-{namespace}@knative.dev exists in knative/community repo.
- Deploy the namespace and rbac binding in the prow/
