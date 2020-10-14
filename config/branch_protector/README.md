# branch_protector

Configure github branch protection rules according to the specified policy in a
YAML file.

## Usage

To learn more about how it works and how to run this tool locally, please check
[kubernetes/test-infra](https://github.com/kubernetes/test-infra/blob/master/prow/cmd/branchprotector/README.md).

## Prow job

We currently run a postsubmit Prow job and a daily periodic Prow job that
synchronizes the branch protection rules based on the latest
`branch_protector/rules.yaml` file.
