# Boskos config

This directory contains the config for our
[Boskos](https://github.com/kubernetes/test-infra/tree/master/boskos) instance.

- `config_start.yaml` Initial Boskos configuration.
- `config.yaml` Boskos configuration.
- `create_projects.sh` Script to create new Boskos projects and set permissions on them.
- `resources.yaml` Pool of projects used by Boskos.
- `custom_role.yaml` Configuration of the custom role, needs to be updated whenever we change its permissions.
- `make_custom_role.sh`: Script to generate the custom role configuration file.
- `update_custom_role.sh`: Script to update the custom role for all existing boskos projects.
