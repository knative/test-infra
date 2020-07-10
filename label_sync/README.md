# label_sync

Update or migrate github labels on repos in a github org based on a YAML file

## Configuration

A typical labels.yaml file looks like:

```yaml
---
labels:
  - color: 00ff00
    name: lgtm
  - color: ff0000
    name: priority/P0
    previously:
      - color: 0000ff
        name: P0
  - name: dead-label
    color: cccccc
    deleteAfter: 2017-01-01T13:00:00Z
```

This will ensure that:

- there is a green `lgtm` label
- there is a red `priority/P0` label, and previous labels should be migrated to
  it:
  - if a `P0` label exists:
    - if `priority/P0` does not, modify the existing `P0` label
    - if `priority/P0` exists, `P0` labels will be deleted, `priority/P0` labels
      will be added
- if there is a `dead-label` label, it will be deleted after
  2017-01-01T13:00:00Z

## Usage

To learn more about how to run this tool locally, please check
[kubernetes/test-infra](https://github.com/kubernetes/test-infra/blob/master/label_sync/README.md#usage).

## Prow job

We currently run a Prow job every hour that synchronizes these labels if the
`labels.yaml` is changed and submitted.
