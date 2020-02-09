# Staging Prow config

This directory contains the config for our staging
[Prow](https://github.com/kubernetes/test-infra/tree/master/prow) instance.

- `Makefile` Commands to interact with the staging Prow instance regarding
  configs and updates.
- `boskos_resources.yaml` Pool of projects used by Boskos.
- `config.yaml` Generated configuration for Prow.
- `jobs/config.yaml` Generated configuration for Prow jobs.
- `config_staging.yaml` Input configuration for `make_config.go` to generate
  `config.yaml`.
- `plugins.yaml` Generated configuration of the Prow plugins.
- `testgrid.yaml` Generated Testgrid configuration.
