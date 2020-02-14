# Prow config

This directory contains the config for our
[Prow](https://github.com/kubernetes/test-infra/tree/master/prow) instance.

- `Makefile` Commands to interact with the Prow instance regarding configs and
  updates.
- `config_start.yaml` Initial, empty configuration for Prow.
- `cluster/*.yaml` Deployments of the Prow cluster.
- `config.yaml` Generated configuration for Prow.
- `plugins.yaml` Generated configuration of the Prow plugins.
- `jobs/config.yaml` Generated configuration of the Prow jobs.
- `testgrid.yaml` Generated Testgrid configuration.
- `config_knative.yaml` Input configuration for `config-generator` tool to generate
  `config.yaml`, `plugins.yaml` and `testgrid.yaml`.
- `run_job.sh` Convenience script to start a Prow job from command-line.
- `pj-on-kind.sh` Convenience script to start a Prow job on kind from command-line.
