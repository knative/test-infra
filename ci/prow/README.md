# Prow config

This directory contains the config for our
[Prow](https://github.com/kubernetes/test-infra/tree/master/prow) instance.

- `Makefile` Commands to interact with the Prow instance regarding configs and
  updates.
- `boskos_resources.yaml` Pool of projects used by Boskos.
- `create_boskos_projects.sh` Script to create new Boskos projects.
- `deployments/*.yaml` Deployments of the Prow cluster.
- `config.yaml` Generated configuration of the Prow jobs.
- `config_knative.yaml` Input configuration for `make_config.go` to generate
  `config.yaml`.
- `config_start.yaml` Initial, empty configuration for Prow.
- `make_config.go` `issue_tracker_config.go` `periodic_config.go` `testgrid_config.go`
  `templates` Tool that generates `config.yaml`, `plugins.yaml` and `testgrid.yaml`
  from `config_knative.yaml`.
- `plugins.yaml` Generated configuration of the Prow plugins.
- `run_job.sh` Convenience script to start a Prow job from command-line.
- `set_boskos_permissions.sh` Script to set up permissions for a Boskos project.
- `testgrid.yaml` Generated Testgrid configuration.
