# Prow config

This directory contains the config for our
[Prow](https://github.com/kubernetes/test-infra/tree/master/prow) instance.

- `Makefile` Commands to interact with the Prow instance regarding configs and
  updates.
- `config_start.yaml` Initial, empty configuration for Prow.
- `deployments/*.yaml` Deployments of the Prow cluster.
- `config.yaml` Generated configuration for Prow.
- `plugins.yaml` Generated configuration of the Prow plugins.
- `jobs/config.yaml` Generated configuration of the Prow jobs.
- `testgrid.yaml` Generated Testgrid configuration.
- `config_knative.yaml` Input configuration for `config-generator` tool to
  generate `config.yaml`, `plugins.yaml` and `testgrid.yaml`.
- `run_job.sh` Convenience script to start a Prow job from command-line.
- `pj-on-kind.sh` Convenience script to start a Prow job on from command-line.
- `boskos_resources.yaml` Pool of projects used by Boskos.
- `create_boskos_projects.sh` Script to create new Boskos projects.
- `set_up_boskos_project.sh` Script to set up a Boskos project.
- `update_all_boskos_projects.sh` Script to set up all Boskos projects.
