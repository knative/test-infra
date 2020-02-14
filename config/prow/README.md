# Prow config

This directory contains the config for our
[Prow](https://github.com/kubernetes/test-infra/tree/master/prow) instance.

- `Makefile` Commands to interact with the Prow instance regarding configs and
  updates.
- `config_start.yaml` Initial, empty configuration for Prow.
- `cluster/*.yaml` Deployments of the Prow cluster.
- `core/*yaml` Generated core configuration for Prow.
- `jobs/config.yaml` Generated configuration of the Prow jobs.
- `testgrid.yaml` Generated Testgrid configuration.
- `config_knative.yaml` Input configuration for `config-generator` tool to generate
  `core/config.yaml`, `core/plugins.yaml`, `jobs/config.yaml` and `testgrid.yaml`.
- `run_job.sh` Convenience script to start a Prow job from command-line.
- `pj-on-kind.sh` Convenience script to start a Prow job on kind from command-line.
